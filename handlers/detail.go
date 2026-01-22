package handlers

import (
	"API-demon-slayyyyy-/services"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// DetailHandler affiche la fiche détaillée d'un personnage
func DetailHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID depuis l'URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ID invalide: %s", idStr)
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	// Récupérer le personnage depuis l'API
	apiService := services.GetAPIService()
	character, err := apiService.GetCharacterByID(id)
	if err != nil {
		log.Printf("Personnage %d non trouvé: %v", id, err)
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	// Récupérer le style de combat depuis l'API
	combatStyle, _ := apiService.GetCombatStyleForCharacter(id)

	// Charger le template
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/detail.html")
	if err != nil {
		log.Printf("Erreur chargement template detail: %v", err)
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
		log.Printf("Erreur rendu template detail: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}
