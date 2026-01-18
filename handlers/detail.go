package handlers

import (
	"API-demon-slayyyyy-/services"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CharacterDetailHandler affiche la fiche détaillée d'un personnage
func CharacterDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID depuis l'URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("❌ ID invalide: %s", idStr)
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	// Récupérer le personnage depuis l'API
	apiService := services.GetAPIService()
	character, err := apiService.GetCharacterByID(id)
	if err != nil {
		log.Printf("❌ Personnage %d non trouvé: %v", id, err)
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	// Extraire le style de combat depuis la description
	combatStyle := extractCombatStyle(character.Description)

	// Charger le template
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/character_detail.html")
	if err != nil {
		log.Printf("❌ Erreur chargement template character_detail: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Données à passer au template
	data := struct {
		PageTitle   string
		User        interface{}
		Character   interface{}
		CombatStyle string
	}{
		PageTitle:   character.Name + " - Demon Slayer",
		User:        getUserFromCookie(r),
		Character:   character,
		CombatStyle: combatStyle,
	}

	// Rendre le template
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("❌ Erreur rendu template character_detail: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}

// extractCombatStyle extrait le style de combat depuis la description
func extractCombatStyle(description string) string {
	// Regex simple pour trouver "X Breathing" dans la description
	// Pour l'instant, on retourne juste un message par défaut
	// Tu peux améliorer ça avec une vraie regex
	if description == "" {
		return "Aucun style de combat connu"
	}

	// Chercher "Breathing" dans la description
	// (À améliorer avec une vraie regex)
	return "Style de combat trouvé dans la description"
}
