package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// ErrorHandler affiche la page 404
func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	// Charger le template 404
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/404.html")
	if err != nil {
		log.Printf("❌ Erreur chargement template 404: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Définir le status code à 404
	w.WriteHeader(http.StatusNotFound)

	// Données à passer au template
	data := struct {
		PageTitle string
		User      interface{}
	}{
		PageTitle: "404 - Page non trouvée",
		User:      getUserFromCookie(r),
	}

	// Rendre le template
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("❌ Erreur rendu template 404: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}
