package httpserver

// HTTPServer is a HTTP server for ports
type HTTPServer struct {
	userChatService IUserChatService
	// userService     IUserService
	tokenService ITokenService
	stateService IStateService
}

// NewHTTPServer creates a new HTTP server for ports
func NewHTTPServer(
	userChatService IUserChatService,
	// userService IUserService,
	tokenService ITokenService,
	stateService IStateService,
) HTTPServer {
	return HTTPServer{
		userChatService: userChatService,
		// userService:     userService,
		tokenService: tokenService,
		stateService: stateService,
	}
}
