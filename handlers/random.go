package handlers

import (
	"API-demon-slayyyyy-/services"
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// RandomHandler affiche la page "aléatoire" avec le bouton "CHOOSE FOR ME"
func RandomHandler(w http.ResponseWriter, r *http.Request) {
	// Afficher la page avec le bouton
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

// getUserFromCookie est déjà définie dans auth.go
// (tous les fichiers du package handlers partagent les fonctions)

// GetRandomCharacterID retourne l'ID d'un personnage aléatoire (API interne)
func GetRandomCharacterID(w http.ResponseWriter, r *http.Request) {
	// Récupérer tous les personnages
	apiService := services.GetAPIService()
	characters, err := apiService.GetAllCharacters()
	if err != nil {
		log.Printf("❌ Erreur récupération personnages: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	if len(characters) == 0 {
		http.Error(w, "Aucun personnage trouvé", http.StatusNotFound)
		return
	}

	// Sélectionner un personnage aléatoire
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(characters))
	randomCharacter := characters[randomIndex]

	// Retourner l'ID en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"id": randomCharacter.ID,
	})
}
