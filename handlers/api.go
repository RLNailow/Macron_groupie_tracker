package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"API-demon-slayyyyy-/services"
)

// APICharactersHandler retourne tous les personnages en JSON
func APICharactersHandler(w http.ResponseWriter, r *http.Request) {
	apiService := services.GetAPIService()

	// Récupérer tous les personnages
	characters, err := apiService.GetAllCharacters()
	if err != nil {
		log.Printf("Erreur API characters: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erreur serveur"})
		return
	}

	// Retourner en JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(characters); err != nil {
		log.Printf("Erreur encodage JSON: %v", err)
	}
}

// APICombatStylesHandler retourne tous les styles de combat en JSON
func APICombatStylesHandler(w http.ResponseWriter, r *http.Request) {
	apiService := services.GetAPIService()

	// Récupérer les styles de combat
	styles, err := apiService.GetBreathingTechniques()
	if err != nil {
		log.Printf("Erreur API styles: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erreur serveur"})
		return
	}

	// Retourner en JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(styles); err != nil {
		log.Printf("Erreur encodage JSON: %v", err)
	}
}
