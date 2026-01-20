package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// HomeHandler affiche la page d'accueil
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Charger le template
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/home.html")
	if err != nil {
		log.Printf("Erreur chargement template home: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Données à passer au template
	data := struct {
		PageTitle string
		User      interface{} // Pour plus tard, gérer l'utilisateur connecté
	}{
		PageTitle: "Demon Slayer - Kimetsu no Yaiba",
		User:      getUserFromCookie(r), // Vérifier si utilisateur connecté
	}

	// Rendre le template
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Erreur rendu template home: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}
