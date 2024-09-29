package httpserver

// HttpServer is a HTTP server for ports
type HttpServer struct {
	userService  IUserService
	tokenService ITokenService
}

// NewHttpServer creates a new HTTP server for ports
func NewHttpServer(
	userService IUserService,
	tokenService ITokenService,
) HttpServer {
	return HttpServer{
		userService:  userService,
		tokenService: tokenService,
	}
}
