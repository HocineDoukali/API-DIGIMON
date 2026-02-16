package helper

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Variable globale qui contiendra tous les templates chargés
var listeTemplate *template.Template

// Load charge tous les fichiers HTML depuis le dossier ../templates
func Load() {
	// Chargement des fichiers .html dans le dossier templates
	temp, tempErr := template.ParseGlob("../templates/*.html")
	if tempErr != nil {
		// En cas d'erreur, le programme s'arrête avec un message d'erreur
		log.Fatalf("Erreur template - %s", tempErr.Error())
		return
	}
	// Affectation des templates à la variable globale
	listeTemplate = temp
	fmt.Println("Template - chargement des templates terminé")
}

// RenderTemplate exécute le template spécifié et écrit le résultat dans la réponse HTTP
func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	var buffer bytes.Buffer

	// Exécution du template avec les données fournies
	errRender := listeTemplate.ExecuteTemplate(&buffer, name, data)
	if errRender != nil {
		// Si une erreur survient, on retourne une erreur 500 au client
		fmt.Println(errRender)
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}

	// Écriture du contenu généré dans la réponse HTTP
	buffer.WriteTo(w)
}