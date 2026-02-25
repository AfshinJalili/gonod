package server

func (s *Server) registerRoutes() {

	s.mux.Handle("GET /health", s.healthHandler)

	s.mux.HandleFunc("POST /register", s.authHandler.Register)
	s.mux.HandleFunc("POST /login", s.authHandler.Login)
}
