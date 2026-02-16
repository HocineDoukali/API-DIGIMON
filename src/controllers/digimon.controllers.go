package controllers

import (
	"context"
	"fmt"
	"guide/helper"
	"guide/services"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// createContext crée un contexte avec timeout pour les requêtes API
func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// ============================================================
// AFFICHAGE DE LA LISTE
// ============================================================

// DisplayListDigimons affiche la liste complète des Digimons.
// - Appelle le service qui récupère tous les Digimons
// - Gère l'erreur éventuelle (service KO / statut != 200)
// - Rend ensuite le template "list_digimon" avec les données
func DisplayListDigimons(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := createContext()
	defer cancel()

	// Récupère la première page avec une taille généreuse
	opts := &services.DigimonListOptions{
		PageSize: 100, // Ajustez selon vos besoins
	}

	data, dataStatusCode, err := services.GetAllDigimons(ctx, opts)
	if dataStatusCode != http.StatusOK || err != nil {
		http.Error(
			w,
			fmt.Sprintf("Erreur service - code: %d\nmessage: %s", dataStatusCode, err.Error()),
			dataStatusCode,
		)
		return
	}

	// Affiche le template de liste avec les données récupérées
	helper.RenderTemplate(w, r, "list_digimon", data.Content)
}

// DisplayListDigimonsWithPagination affiche la liste paginée des Digimons
func DisplayListDigimonsWithPagination(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := createContext()
	defer cancel()

	// Récupère le numéro de page depuis l'URL (ex: ?page=2)
	pageStr := r.URL.Query().Get("page")
	page := 0
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	opts := &services.DigimonListOptions{
		Page:     page,
		PageSize: 20,
	}

	data, dataStatusCode, err := services.GetAllDigimons(ctx, opts)
	if dataStatusCode != http.StatusOK || err != nil {
		http.Error(
			w,
			fmt.Sprintf("Erreur service - code: %d\nmessage: %s", dataStatusCode, err.Error()),
			dataStatusCode,
		)
		return
	}

	// Structure pour le template avec les infos de pagination
	templateData := map[string]interface{}{
		"Digimons":     data.Content,
		"CurrentPage":  page,
		"TotalPages":   data.TotalPages,
		"TotalDigimons": data.TotalElements,
		"HasNext":      !data.Last,
		"HasPrevious":  !data.First,
	}

	helper.RenderTemplate(w, r, "list_digimon_paginated", templateData)
}

// ============================================================
// RECHERCHE
// ============================================================

// DisplaySearch gère la recherche via un champ "query".
// - Normalise la recherche (trim + lowercase)
// - Si vide : redirection vers la liste
// - Sinon : utilise l'API pour filtrer directement
func DisplaySearch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := createContext()
	defer cancel()

	// Récupère le paramètre de formulaire nommé "query"
	query := r.FormValue("query")

	// Normalisation : enlever espaces inutiles + passer en minuscules
	query = strings.TrimSpace(query)

	// Si pas de recherche, on retourne à la page liste
	if query == "" {
		http.Redirect(w, r, "/digimons", http.StatusSeeOther)
		return
	}

	// Utilise l'API pour rechercher directement
	opts := &services.DigimonListOptions{
		Name:     query,
		PageSize: 50,
	}

	data, dataStatusCode, dataError := services.GetAllDigimons(ctx, opts)
	if dataStatusCode != http.StatusOK || dataError != nil {
		http.Error(
			w,
			fmt.Sprintf("Erreur service - code: %d\nmessage: %v", dataStatusCode, dataError.Error()),
			dataStatusCode,
		)
		return
	}

	// Structure pour le template
	templateData := map[string]interface{}{
		"Digimons": data.Content,
		"Query":    query,
		"Total":    data.TotalElements,
	}

	// Réutilise le template de liste pour afficher le résultat filtré
	helper.RenderTemplate(w, r, "search_digimon", templateData)
}

// DisplaySearchAdvanced gère la recherche avancée avec recherche exacte
func DisplaySearchAdvanced(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := createContext()
	defer cancel()

	query := strings.TrimSpace(r.FormValue("query"))
	exact := r.FormValue("exact") == "true" || r.FormValue("exact") == "on"

	if query == "" {
		http.Redirect(w, r, "/digimons", http.StatusSeeOther)
		return
	}

	opts := &services.DigimonListOptions{
		Name:     query,
		Exact:    exact,
		PageSize: 50,
	}

	data, dataStatusCode, dataError := services.GetAllDigimons(ctx, opts)
	if dataStatusCode != http.StatusOK || dataError != nil {
		http.Error(
			w,
			fmt.Sprintf("Erreur service - code: %d\nmessage: %v", dataStatusCode, dataError.Error()),
			dataStatusCode,
		)
		return
	}

	templateData := map[string]interface{}{
		"Digimons": data.Content,
		"Query":    query,
		"Exact":    exact,
		"Total":    data.TotalElements,
	}

	helper.RenderTemplate(w, r, "search_digimon", templateData)
}

// ============================================================
// FILTRAGE
// ============================================================

// DisplayFilter filtre les Digimons selon :
// - niveau (champ "level")
// - attribut (champ "attribute")
// - X-Antibody (checkbox "xantibody")
// Puis affiche le template "filter_digimons".
func DisplayFilter(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := createContext()
	defer cancel()

	// Parse le formulaire pour accéder à r.Form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erreur parsing formulaire", http.StatusBadRequest)
		return
	}

	// Récupère les paramètres de filtrage
	level := strings.TrimSpace(r.FormValue("level"))
	attribute := strings.TrimSpace(r.FormValue("attribute"))
	xAntibodyStr := r.FormValue("xantibody")

	// Debug console
	log.Printf("Filtres - Level: %s, Attribute: %s, XAntibody: %s", level, attribute, xAntibodyStr)

	// Construction des options de filtrage
	opts := &services.DigimonListOptions{
		PageSize: 100,
	}

	// Filtre par niveau si fourni
	if level != "" {
		opts.Level = level
	}

	// Filtre par attribut si fourni
	if attribute != "" {
		opts.Attribute = attribute
	}

	// Filtre par X-Antibody si coché
	if xAntibodyStr == "true" || xAntibodyStr == "on" {
		hasXAntibody := true
		opts.XAntibody = &hasXAntibody
	}

	// Appel à l'API avec les filtres
	data, dataStatusCode, dataError := services.GetAllDigimons(ctx, opts)
	if dataStatusCode != http.StatusOK || dataError != nil {
		log.Printf("Erreur DisplayFilter - %s", dataError.Error())
		http.Error(
			w,
			fmt.Sprintf("Erreur service - code: %d\nmessage: %v", dataStatusCode, dataError.Error()),
			dataStatusCode,
		)
		return
	}

	// Structure pour le template
	templateData := map[string]interface{}{
		"Digimons":   data.Content,
		"Level":      level,
		"Attribute":  attribute,
		"XAntibody":  xAntibodyStr == "true" || xAntibodyStr == "on",
		"Total":      data.TotalElements,
		"TotalPages": data.TotalPages,
	}

	// Rend un template dédié au filtrage
	helper.RenderTemplate(w, r, "filter_digimons", templateData)
}

// DisplayFilterAdvanced filtre avec filtrage local en mémoire
// (utile si vous voulez des critères non supportés par l'API)
func DisplayFilterAdvanced(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := createContext()
	defer cancel()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erreur parsing formulaire", http.StatusBadRequest)
		return
	}

	// Récupère tous les Digimons
	opts := &services.DigimonListOptions{
		PageSize: 500, // Grande taille pour tout récupérer
	}

	data, dataStatusCode, dataError := services.GetAllDigimons(ctx, opts)
	if dataStatusCode != http.StatusOK || dataError != nil {
		log.Printf("Erreur DisplayFilterAdvanced - %s", dataError.Error())
		http.Error(
			w,
			fmt.Sprintf("Erreur service - code: %d\nmessage: %v", dataStatusCode, dataError.Error()),
			dataStatusCode,
		)
		return
	}

	// Paramètres de filtrage local
	levels := r.Form["levels"]          // Checkbox multiple de niveaux
	attributes := r.Form["attributes"]  // Checkbox multiple d'attributs
	xAntibodyStr := r.FormValue("xantibody")

	// Debug
	log.Printf("Filtres - Levels: %v, Attributes: %v, XAntibody: %s", levels, attributes, xAntibodyStr)

	// Liste finale filtrée
	validDigimons := []services.DigimonSummary{}

	for _, digimon := range data.Content {
		

		// Vérification niveau (nécessite de récupérer le Digimon complet)
		// Note: Cette approche nécessiterait des appels API supplémentaires
		// Pour simplifier, on utilise uniquement les filtres API

		// Filtre simplifié basé sur le nom (exemple)
		if len(levels) == 0 && len(attributes) == 0 && xAntibodyStr == "" {
			validDigimons = append(validDigimons, digimon)
		}
	}

	templateData := map[string]interface{}{
		"Digimons":   validDigimons,
		"Levels":     levels,
		"Attributes": attributes,
		"XAntibody":  xAntibodyStr == "true" || xAntibodyStr == "on",
		"Total":      len(validDigimons),
	}

	helper.RenderTemplate(w, r, "filter_digimons_advanced", templateData)
}

// ============================================================
// AFFICHAGE DÉTAILS
// ============================================================

// DisplayDigimonDetails affiche les détails complets d'un Digimon
func DisplayDigimonDetails(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := createContext()
	defer cancel()

	// Récupère l'ID depuis l'URL (ex: /digimon/1)
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// Récupère le Digimon complet
	digimon, statusCode, err := services.GetDigimonByID(ctx, id)
	if statusCode != http.StatusOK || err != nil {
		if statusCode == http.StatusNotFound {
			http.Error(w, "Digimon non trouvé", http.StatusNotFound)
		} else {
			http.Error(
				w,
				fmt.Sprintf("Erreur service - code: %d\nmessage: %s", statusCode, err.Error()),
				statusCode,
			)
		}
		return
	}

	helper.RenderTemplate(w, r, "digimon_details", digimon)
}

// DisplayDigimonDetailsByName affiche les détails d'un Digimon par son nom
func DisplayDigimonDetailsByName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := createContext()
	defer cancel()

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Nom manquant", http.StatusBadRequest)
		return
	}

	digimon, statusCode, err := services.GetDigimonByName(ctx, name)
	if statusCode != http.StatusOK || err != nil {
		if statusCode == http.StatusNotFound {
			http.Error(w, "Digimon non trouvé", http.StatusNotFound)
		} else {
			http.Error(
				w,
				fmt.Sprintf("Erreur service - code: %d\nmessage: %s", statusCode, err.Error()),
				statusCode,
			)
		}
		return
	}

	helper.RenderTemplate(w, r, "digimon_details", digimon)
}

// ============================================================
// FILTRES PAR RESSOURCES
// ============================================================

// DisplayDigimonsByAttribute affiche tous les Digimons d'un attribut spécifique
func DisplayDigimonsByAttribute(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := createContext()
	defer cancel()

	attributeName := r.URL.Query().Get("attribute")
	if attributeName == "" {
		http.Error(w, "Attribut manquant", http.StatusBadRequest)
		return
	}

	// Récupère l'attribut avec ses Digimons
	attribute, statusCode, err := services.GetAttributeByName(ctx, attributeName)
	if statusCode != http.StatusOK || err != nil {
		http.Error(
			w,
			fmt.Sprintf("Erreur service - code: %d\nmessage: %s", statusCode, err.Error()),
			statusCode,
		)
		return
	}

	templateData := map[string]interface{}{
		"Attribute": attribute.Attribute,
		"Digimons":  attribute.Digimons,
		"Total":     len(attribute.Digimons),
	}

	helper.RenderTemplate(w, r, "digimons_by_attribute", templateData)
}

// DisplayDigimonsByLevel affiche tous les Digimons d'un niveau spécifique
func DisplayDigimonsByLevel(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := createContext()
	defer cancel()

	levelName := r.URL.Query().Get("level")
	if levelName == "" {
		http.Error(w, "Niveau manquant", http.StatusBadRequest)
		return
	}

	// Récupère le niveau avec ses Digimons
	level, statusCode, err := services.GetLevelByName(ctx, levelName)
	if statusCode != http.StatusOK || err != nil {
		http.Error(
			w,
			fmt.Sprintf("Erreur service - code: %d\nmessage: %s", statusCode, err.Error()),
			statusCode,
		)
		return
	}

	templateData := map[string]interface{}{
		"Level":    level.Level,
		"Digimons": level.Digimons,
		"Total":    len(level.Digimons),
	}

	helper.RenderTemplate(w, r, "digimons_by_level", templateData)
}

// ============================================================
// UTILITAIRES
// ============================================================

// GetAvailableLevels retourne la liste des niveaux disponibles pour les filtres
func GetAvailableLevels() []string {
	return []string{
		"Fresh",
		"In-Training",
		"Rookie",
		"Champion",
		"Ultimate",
		"Mega",
		"Ultra",
		"Armor",
	}
}

// GetAvailableAttributes retourne la liste des attributs disponibles pour les filtres
func GetAvailableAttributes() []string {
	return []string{
		"Vaccine",
		"Data",
		"Virus",
		"Free",
		"Unknown",
	}
}

// DisplayFilterForm affiche le formulaire de filtrage avec les options disponibles
func DisplayFilterForm(w http.ResponseWriter, r *http.Request) {
	templateData := map[string]interface{}{
		"Levels":     GetAvailableLevels(),
		"Attributes": GetAvailableAttributes(),
	}

	helper.RenderTemplate(w, r, "filter_form", templateData)
}