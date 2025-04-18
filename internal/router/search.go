package router

// func addSearchRoutes(
// 	mux *http.ServeMux,
// 	searchService service.SearchService,
// ) {
// 	authHandler := handler.NewAuthHandler(searchService)

// 	mux.HandleFunc("GET /auth/password", authHandler.PasswordLogInPage)
// 	mux.HandleFunc("POST /auth/password", authHandler.PasswordLogIn)

// 	// QRcode login page
// 	mux.HandleFunc("GET /auth/password/qrcode", authHandler.QRcodeLogInPage)
// 	mux.HandleFunc("POST /auth/password/qrcode", authHandler.QRcodeLogIn)

// 	mux.HandleFunc("/auth/logout", authHandler.Logout)
// }
