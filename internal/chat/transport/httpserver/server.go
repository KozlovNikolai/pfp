package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"github.com/KozlovNikolai/pfp/internal/chat/repository/pgrepo"
	"github.com/KozlovNikolai/pfp/internal/chat/services"
	"github.com/KozlovNikolai/pfp/internal/chat/transport/httpserver/middlewares"
	"github.com/KozlovNikolai/pfp/internal/chat/transport/ws"
	"github.com/KozlovNikolai/pfp/internal/pkg/config"
	"github.com/KozlovNikolai/pfp/internal/pkg/pg"
)

// Router is ...
type Router struct {
	router *gin.Engine
	logger *zap.Logger
}

// NewRouter is ...
func NewRouter() *Router {
	// Инициализация логгера Zap
	//	logger, err := zap.NewProduction()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	var userRepo services.IUserRepository
	// Выбор репозитория
	switch config.Cfg.RepoType {
	case "postgres":
		pgDB, err := pg.Dial(config.Cfg.Storage)
		if err != nil {
			logger.Fatal("pg.Dial failed: %w", zap.Error(err))
		}
		userRepo = pgrepo.NewUserRepo(pgDB)

	default:
		logger.Fatal("Invalid repository type")
	}
	// создаем сервисы
	userService := services.NewUserService(userRepo)
	tokenService := services.NewTokenService(
		userRepo,
		config.Cfg.TokenTimeDuration,
	)

	// создаем http сервер
	httpServer := NewHTTPServer(
		userService,
		tokenService,
	)
	hub := ws.NewHub() // создаем hub
	go hub.Run()
	wsHandler := ws.NewHandler(hub) // создаем websocket handler
	// Создание роутера
	server := &Router{
		router: gin.Default(),
		logger: logger,
	}

	// Middleware
	server.router.Use(middlewares.LoggerMiddleware(logger))
	server.router.Use(middlewares.RequestIDMiddleware())

	// CORS
	server.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*", "https://localhost:8443", "https://127.0.0.1:8443"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// add swagger /docs/index.html
	server.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	open := server.router.Group("/")
	open.POST("signup", httpServer.SignUp)
	open.POST("signin", httpServer.SignIn)
	//open.GET("signout", httpServer.SignOut)
	// доступ для админов
	admin := server.router.Group("/admin/")
	admin.Use(httpServer.CheckAdmin())
	admin.GET("users", httpServer.GetUsers)

	// доступ для любых зарегистрированных пользователей
	authorized := server.router.Group("/auth/")
	//authorized.Use(httpServer.CheckAuthorizedUser())

	authorized.GET("user", httpServer.GetUser)
	authorized.GET("signout", httpServer.SignOut)
	authorized.GET("login", httpServer.LoginUserByToken)
	// создаем websocket сервер

	// websocket routes
	authorized.POST("/ws/createRoom", wsHandler.CreateRoom)
	authorized.GET("/ws/getRooms", wsHandler.GetRooms)
	authorized.GET("/ws/getClients/:roomID", wsHandler.GetClients)
	authorized.GET("/ws/joinRoom/:roomID", wsHandler.JoinRoom)

	return server
}

// Run is ...
func (s *Router) Run() {
	defer func() {
		_ = s.logger.Sync() // flushes buffer, if any
	}()
	// Настройка сервера с таймаутами
	server := &http.Server{
		Addr:         config.Cfg.Address,
		Handler:      s.router,
		ReadTimeout:  config.Cfg.Timeout,
		WriteTimeout: config.Cfg.Timeout,
		IdleTimeout:  config.Cfg.IdleTimout,
	}

	// listen to OS signals and gracefully shutdown HTTP server
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(stopped)
	}()
	if err := server.ListenAndServeTLS(config.CertFile, config.KeyFile); err != nil &&
		err != http.ErrServerClosed {
		s.logger.Fatal(fmt.Sprintf("Could not listen on %s", config.Cfg.Address), zap.Error(err))
	}
	<-stopped

	log.Printf("Bye!")
}
