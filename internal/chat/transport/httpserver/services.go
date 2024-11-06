package httpserver

// HTTPServer is a HTTP server for ports
type HTTPServer struct {
	userService  IUserService
	tokenService ITokenService
}

// NewHTTPServer creates a new HTTP server for ports
func NewHTTPServer(
	userService IUserService,
	tokenService ITokenService,
) HTTPServer {
	return HTTPServer{
		userService:  userService,
		tokenService: tokenService,
	}
}
