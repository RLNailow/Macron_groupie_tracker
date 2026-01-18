package handlers

import (
	"API-demon-slayyyyy-/services"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// RandomHandler affiche la page "aléatoire" avec le bouton "CHOOSE FOR ME"
func RandomHandler(w http.ResponseWriter, r *http.Request) {
	// Si c'est une requête POST, rediriger vers un personnage aléatoire
	if r.Method == "POST" {
		redirectToRandomCharacter(w, r)
		return
	}

	// Sinon, afficher la page avec le bouton
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/random.html")
	if err != nil {
		log.Printf("❌ Erreur chargement template random: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	data := struct {
		PageTitle string
		User      interface{}
	}{
		PageTitle: "Aléatoire - Demon Slayer",
		User:      getUserFromCookie(r),
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("❌ Erreur rendu template random: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}

// redirectToRandomCharacter redirige vers un personnage aléatoire
func redirectToRandomCharacter(w http.ResponseWriter, r *http.Request) {
	// Récupérer tous les personnages
	apiService := services.GetAPIService()
	characters, err := apiService.GetAllCharacters()
	if err != nil {
		log.Printf("❌ Erreur récupération personnages: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Sélectionner un personnage aléatoire
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(characters))
	randomCharacter := characters[randomIndex]

	// Rediriger vers la page de détail du personnage
	redirectURL := fmt.Sprintf("/characters/%d", randomCharacter.ID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
