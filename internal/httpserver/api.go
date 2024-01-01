package httpserver

func (hs *HttpServer) apiRoutes() {
	hs.router.GET("/api/realms", hs.GetRealms)
}
