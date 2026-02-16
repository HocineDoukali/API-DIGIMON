package routes

import (
	"guide/controllers"
	"net/http"
)

// digimonsRoutes configure toutes les routes liées aux Digimons
func digimonsRoutes(router *http.ServeMux) {
	// ============================================================
	// LISTE ET PAGINATION
	// ============================================================
	
	// Liste complète des Digimons (première page)
	router.HandleFunc("/digimons", controllers.DisplayListDigimons)
	
	// Liste paginée des Digimons avec navigation
	router.HandleFunc("/digimons/paginated", controllers.DisplayListDigimonsWithPagination)

	// ============================================================
	// RECHERCHE
	// ============================================================
	
	// Recherche simple par nom
	router.HandleFunc("/digimons/search", controllers.DisplaySearch)
	
	// Recherche avancée (avec option exacte)
	router.HandleFunc("/digimons/search/advanced", controllers.DisplaySearchAdvanced)

	// ============================================================
	// FILTRAGE
	// ============================================================
	
	// Formulaire de filtrage
	router.HandleFunc("/digimons/filter/form", controllers.DisplayFilterForm)
	
	// Filtrage standard (niveau, attribut, X-Antibody)
	router.HandleFunc("/digimons/filter", controllers.DisplayFilter)
	
	// Filtrage avancé (avec filtres multiples en mémoire)
	router.HandleFunc("/digimons/filter/advanced", controllers.DisplayFilterAdvanced)

	// ============================================================
	// DÉTAILS
	// ============================================================
	
	// Détails d'un Digimon par ID
	router.HandleFunc("/digimon/details", controllers.DisplayDigimonDetails)
	
	// Détails d'un Digimon par nom
	router.HandleFunc("/digimon/details/name", controllers.DisplayDigimonDetailsByName)

	// ============================================================
	// PAR RESSOURCES
	// ============================================================
	
	// Liste des Digimons par attribut (Vaccine, Virus, Data, etc.)
	router.HandleFunc("/digimons/by-attribute", controllers.DisplayDigimonsByAttribute)
	
	// Liste des Digimons par niveau (Rookie, Champion, Ultimate, etc.)
	router.HandleFunc("/digimons/by-level", controllers.DisplayDigimonsByLevel)
}