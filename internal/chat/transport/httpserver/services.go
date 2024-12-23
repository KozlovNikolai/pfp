package httpserver

import "github.com/KozlovNikolai/pfp/internal/chat/transport/ws"

// HTTPServer is a HTTP server for ports
type HTTPServer struct {
	accountService IAccountService
	userService    IUserService
	tokenService   ITokenService
	stateService   IStateService
	chatService    IChatService
	msgService     IMessageService
	wsHandler      *ws.Handler
}

// NewHTTPServer creates a new HTTP server for ports
func NewHTTPServer(
	accountService IAccountService,
	userService IUserService,
	chatService IChatService,
	tokenService ITokenService,
	stateService IStateService,
	msgService IMessageService,
	wsHandler *ws.Handler,
) HTTPServer {
	return HTTPServer{
		accountService: accountService,
		userService:    userService,
		tokenService:   tokenService,
		stateService:   stateService,
		chatService:    chatService,
		msgService:     msgService,
		wsHandler:      wsHandler,
	}
}
