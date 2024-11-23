package httpserver

// HTTPServer is a HTTP server for ports
type HTTPServer struct {
	userService  IUserService
	tokenService ITokenService
	stateService IStateService
	chatService  IChatService
}

// NewHTTPServer creates a new HTTP server for ports
func NewHTTPServer(
	userService IUserService,
	chatService IChatService,
	tokenService ITokenService,
	stateService IStateService,
) HTTPServer {
	return HTTPServer{
		userService:  userService,
		tokenService: tokenService,
		stateService: stateService,
		chatService:  chatService,
	}
}
