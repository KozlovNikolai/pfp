package httpserver

import "github.com/KozlovNikolai/pfp/internal/chat/transport/ws"

// HTTPServer is a HTTP server for ports
type HTTPServer struct {
	userService  IUserService
	tokenService ITokenService
	stateService IStateService
	chatService  IChatService
	msgService   IMessageService
	wsHandler    *ws.Handler
}

// NewHTTPServer creates a new HTTP server for ports
func NewHTTPServer(
	userService IUserService,
	chatService IChatService,
	tokenService ITokenService,
	stateService IStateService,
	msgService IMessageService,
	wsHandler *ws.Handler,
) HTTPServer {
	return HTTPServer{
		userService:  userService,
		tokenService: tokenService,
		stateService: stateService,
		chatService:  chatService,
		msgService:   msgService,
		wsHandler:    wsHandler,
	}
}
