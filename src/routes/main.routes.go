package routes

import (
	"net/http"
)

// MainRouter initialise et retourne le routeur principal de l'application
func MainRouter() *http.ServeMux {

	// Création du routeur principal
	mainRouter := http.NewServeMux()

	// Enregistrement des routes Digimon
	digimonsRoutes(mainRouter)
	
	// Routes de test (si vous en avez besoin)
	testRoutes(mainRouter)

	// Configuration du serveur de fichiers statiques (CSS, images, etc.)
	fileServerHandler := http.FileServer(http.Dir("./../assets"))

	// Route permettant de servir les fichiers statiques via /static/
	mainRouter.Handle("/static/", http.StripPrefix("/static/", fileServerHandler))

	return mainRouter
}