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

	"github.com/KozlovNikolai/pfp/internal/chat/repository/pgrepo"
	"github.com/KozlovNikolai/pfp/internal/chat/services"
	"github.com/KozlovNikolai/pfp/internal/chat/transport/httpserver/middlewares"
	"github.com/KozlovNikolai/pfp/internal/pkg/config"
	"github.com/KozlovNikolai/pfp/internal/pkg/pg"
	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Server is ...
type Server struct {
	router *gin.Engine
	logger *zap.Logger
}

// NewServer is ...
func NewServer() *Server {
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
	tokenService := services.NewTokenService(config.Cfg.TokenTimeDuration)
	// создаем сервер
	httpServer := NewHttpServer(
		userService,
		tokenService,
	)

	// Создание сервера
	server := &Server{
		router: gin.Default(),
		logger: logger,
	}

	// Middleware
	server.router.Use(middlewares.LoggerMiddleware(logger))
	server.router.Use(middlewares.RequestIDMiddleware())

	// CORS
	server.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost:8443", "https://127.0.0.1:8443"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// add swagger
	server.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	open := server.router.Group("/")
	open.POST("signup", httpServer.SignUp)
	open.POST("signin", httpServer.SignIn)

	// Закрытые маршруты
	// доступ только для администратора
	admin := server.router.Group("/")
	admin.Use(httpServer.CheckAdmin())

	admin.GET("users", httpServer.GetUsers)

	// доступ для любых зарегистрированных пользователей
	authorized := server.router.Group("/")
	authorized.Use(httpServer.CheckAuthorizedUser())

	authorized.GET("user", httpServer.GetUser)

	return server
}

// Run is ...
func (s *Server) Run() {
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
	if err := server.ListenAndServeTLS(config.CertFile, config.KeyFile); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal(fmt.Sprintf("Could not listen on %s", config.Cfg.Address), zap.Error(err))
	}
	<-stopped

	log.Printf("Bye!")
}
