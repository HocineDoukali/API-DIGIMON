package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Configuration de l'API
const (
	digimonAPIBaseURL = "https://digi-api.com/api/v1"
	defaultTimeout    = 10 * time.Second
)

// Client HTTP réutilisable
var httpClient = &http.Client{
	Timeout: defaultTimeout,
}

// ============================================================
// STRUCTURES DE DONNÉES
// ============================================================

// Image représente les différentes images d'un Digimon
type Image struct {
	Href        string `json:"href"`
	Transparent string `json:"transparent,omitempty"`
}

// DigimonType représente un type de Digimon
type DigimonType struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Image string `json:"image,omitempty"`
}

// DigimonAttribute représente un attribut (Vaccine, Virus, Data, etc.)
type DigimonAttribute struct {
	ID        int    `json:"id"`
	Attribute string `json:"attribute"`
	Image     string `json:"image,omitempty"`
}

// DigimonLevel représente le niveau d'évolution
type DigimonLevel struct {
	ID    int    `json:"id"`
	Level string `json:"level"`
}

// DigimonField représente un champ/famille
type DigimonField struct {
	ID    int    `json:"id"`
	Field string `json:"field"`
	Image string `json:"image,omitempty"`
}

// DigimonSkill représente une compétence
type DigimonSkill struct {
	ID          int    `json:"id"`
	Skill       string `json:"skill"`
	Description string `json:"description,omitempty"`
}

// Digimon représente un Digimon complet
type Digimon struct {
	ID          int                `json:"id"`
	Name        string             `json:"name"`
	XAntibody   bool               `json:"xAntibody"`
	Images      []Image            `json:"images"`
	Levels      []DigimonLevel     `json:"levels"`
	Types       []DigimonType      `json:"types"`
	Attributes  []DigimonAttribute `json:"attributes"`
	Fields      []DigimonField     `json:"fields"`
	Skills      []DigimonSkill     `json:"skills"`
	Descriptions []Description    `json:"descriptions,omitempty"`
}

// Description représente une description du Digimon
type Description struct {
	Origin      string `json:"origin"`
	Language    string `json:"language"`
	Description string `json:"description"`
}

// DigimonListResponse représente la réponse paginée de la liste
type DigimonListResponse struct {
	Content          []DigimonSummary `json:"content"`
	Pageable         Pageable         `json:"pageable"`
	TotalElements    int              `json:"totalElements"`
	TotalPages       int              `json:"totalPages"`
	Last             bool             `json:"last"`
	First            bool             `json:"first"`
	Size             int              `json:"size"`
	Number           int              `json:"number"`
	NumberOfElements int              `json:"numberOfElements"`
	Empty            bool             `json:"empty"`
}

// DigimonSummary représente un Digimon dans la liste (version simplifiée)
type DigimonSummary struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Href  string `json:"href"`
	Image string `json:"image"`
}

// Pageable contient les informations de pagination
type Pageable struct {
	Sort       Sort `json:"sort"`
	PageNumber int  `json:"pageNumber"`
	PageSize   int  `json:"pageSize"`
	Offset     int  `json:"offset"`
	Paged      bool `json:"paged"`
	Unpaged    bool `json:"unpaged"`
}

// Sort contient les informations de tri
type Sort struct {
	Sorted   bool `json:"sorted"`
	Unsorted bool `json:"unsorted"`
	Empty    bool `json:"empty"`
}

// ============================================================
// FONCTIONS PRINCIPALES
// ============================================================

// GetDigimonByID récupère un Digimon spécifique par son ID
func GetDigimonByID(ctx context.Context, id int) (*Digimon, int, error) {
	url := fmt.Sprintf("%s/digimon/%d", digimonAPIBaseURL, id)
	return fetchDigimon(ctx, url)
}

// GetDigimonByName récupère un Digimon spécifique par son nom
func GetDigimonByName(ctx context.Context, name string) (*Digimon, int, error) {
	url := fmt.Sprintf("%s/digimon/%s", digimonAPIBaseURL, name)
	return fetchDigimon(ctx, url)
}

// fetchDigimon est une fonction helper pour récupérer un Digimon
func fetchDigimon(ctx context.Context, url string) (*Digimon, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, http.StatusInternalServerError,
			fmt.Errorf("erreur création requête: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError,
			fmt.Errorf("erreur requête HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode,
			fmt.Errorf("code HTTP inattendu: %d", resp.StatusCode)
	}

	var digimon Digimon
	if err := json.NewDecoder(resp.Body).Decode(&digimon); err != nil {
		return nil, http.StatusInternalServerError,
			fmt.Errorf("erreur décodage JSON: %w", err)
	}

	return &digimon, resp.StatusCode, nil
}

// DigimonListOptions contient les options de filtrage pour la liste
type DigimonListOptions struct {
	Name       string // Recherche par nom similaire
	Exact      bool   // Recherche exacte du nom
	Attribute  string // Filtrer par attribut (Vaccine, Virus, Data, etc.)
	XAntibody  *bool  // Filtrer par X-Antibody (nil = pas de filtre)
	Level      string // Filtrer par niveau (Fresh, In-Training, Rookie, etc.)
	Page       int    // Numéro de page (commence à 0)
	PageSize   int    // Taille de la page (par défaut: 20)
}

// GetAllDigimons récupère la liste paginée des Digimons avec options de filtrage
func GetAllDigimons(ctx context.Context, opts *DigimonListOptions) (*DigimonListResponse, int, error) {
	url := fmt.Sprintf("%s/digimon", digimonAPIBaseURL)
	
	// Construction de l'URL avec les paramètres de requête
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, http.StatusInternalServerError,
			fmt.Errorf("erreur création requête: %w", err)
	}

	// Ajout des paramètres de requête si options fournies
	if opts != nil {
		q := req.URL.Query()
		
		if opts.Name != "" {
			q.Add("name", opts.Name)
		}
		if opts.Exact {
			q.Add("exact", "true")
		}
		if opts.Attribute != "" {
			q.Add("attribute", opts.Attribute)
		}
		if opts.XAntibody != nil {
			if *opts.XAntibody {
				q.Add("xAntibody", "true")
			} else {
				q.Add("xAntibody", "false")
			}
		}
		if opts.Level != "" {
			q.Add("level", opts.Level)
		}
		if opts.Page > 0 {
			q.Add("page", fmt.Sprintf("%d", opts.Page))
		}
		if opts.PageSize > 0 {
			q.Add("pageSize", fmt.Sprintf("%d", opts.PageSize))
		}
		
		req.URL.RawQuery = q.Encode()
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError,
			fmt.Errorf("erreur requête HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode,
			fmt.Errorf("code HTTP inattendu: %d", resp.StatusCode)
	}

	var listResponse DigimonListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
		return nil, http.StatusInternalServerError,
			fmt.Errorf("erreur décodage JSON: %w", err)
	}

	return &listResponse, resp.StatusCode, nil
}

// ============================================================
// FONCTIONS POUR LES AUTRES RESSOURCES
// ============================================================

// Attribute représente un attribut complet
type Attribute struct {
	ID        int              `json:"id"`
	Attribute string           `json:"attribute"`
	Digimons  []DigimonSummary `json:"digimons"`
}

// GetAttributeByID récupère un attribut par son ID
func GetAttributeByID(ctx context.Context, id int) (*Attribute, int, error) {
	url := fmt.Sprintf("%s/attribute/%d", digimonAPIBaseURL, id)
	return fetchAttribute(ctx, url)
}

// GetAttributeByName récupère un attribut par son nom
func GetAttributeByName(ctx context.Context, name string) (*Attribute, int, error) {
	url := fmt.Sprintf("%s/attribute/%s", digimonAPIBaseURL, name)
	return fetchAttribute(ctx, url)
}

func fetchAttribute(ctx context.Context, url string) (*Attribute, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, http.StatusInternalServerError,
			fmt.Errorf("erreur création requête: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError,
			fmt.Errorf("erreur requête HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode,
			fmt.Errorf("code HTTP inattendu: %d", resp.StatusCode)
	}

	var attribute Attribute
	if err := json.NewDecoder(resp.Body).Decode(&attribute); err != nil {
		return nil, http.StatusInternalServerError,
			fmt.Errorf("erreur décodage JSON: %w", err)
	}

	return &attribute, resp.StatusCode, nil
}

// Level représente un niveau complet
type Level struct {
	ID       int              `json:"id"`
	Level    string           `json:"level"`
	Digimons []DigimonSummary `json:"digimons"`
}

// GetLevelByID récupère un niveau par son ID
func GetLevelByID(ctx context.Context, id int) (*Level, int, error) {
	url := fmt.Sprintf("%s/level/%d", digimonAPIBaseURL, id)
	return fetchLevel(ctx, url)
}

// GetLevelByName récupère un niveau par son nom
func GetLevelByName(ctx context.Context, name string) (*Level, int, error) {
	url := fmt.Sprintf("%s/level/%s", digimonAPIBaseURL, name)
	return fetchLevel(ctx, url)
}

func fetchLevel(ctx context.Context, url string) (*Level, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, http.StatusInternalServerError,
			fmt.Errorf("erreur création requête: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError,
			fmt.Errorf("erreur requête HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode,
			fmt.Errorf("code HTTP inattendu: %d", resp.StatusCode)
	}

	var level Level
	if err := json.NewDecoder(resp.Body).Decode(&level); err != nil {
		return nil, http.StatusInternalServerError,
			fmt.Errorf("erreur décodage JSON: %w", err)
	}

	return &level, resp.StatusCode, nil
}

// ============================================================
// VERSIONS SIMPLIFIÉES SANS CONTEXTE (pour compatibilité)
// ============================================================

// GetDigimonByIDSimple version simplifiée sans contexte
func GetDigimonByIDSimple(id int) (*Digimon, int, error) {
	return GetDigimonByID(context.Background(), id)
}

// GetDigimonByNameSimple version simplifiée sans contexte
func GetDigimonByNameSimple(name string) (*Digimon, int, error) {
	return GetDigimonByName(context.Background(), name)
}

// GetAllDigimonsSimple version simplifiée sans contexte
func GetAllDigimonsSimple(opts *DigimonListOptions) (*DigimonListResponse, int, error) {
	return GetAllDigimons(context.Background(), opts)
}