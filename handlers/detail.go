package handlers

import (
	"API-demon-slayyyyy-/services"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

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

	// Extraire le style de combat depuis la description
	combatStyle := extractCombatStyle(character.Description)

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

// extractCombatStyle extrait le style de combat depuis la description
func extractCombatStyle(description string) string {
	if description == "" {
		return "Aucun style de combat connu"
	}

	descLower := strings.ToLower(description)

	// Liste des styles de combat connus (élargie)
	knownStyles := []string{
		"Water Breathing",
		"Thunder Breathing",
		"Flame Breathing",
		"Wind Breathing",
		"Stone Breathing",
		"Mist Breathing",
		"Serpent Breathing",
		"Insect Breathing",
		"Sound Breathing",
		"Moon Breathing",
		"Sun Breathing",
		"Beast Breathing",
		"Flower Breathing",
		"Love Breathing",
		"Blood Demon Art",
		"Demon Blood Art",
	}

	// 1. Chercher si un style est mentionné exactement
	for _, style := range knownStyles {
		if strings.Contains(descLower, strings.ToLower(style)) {
			return style
		}
	}

	// 2. Chercher des variations (mot-clé + "breath")
	breathingKeywords := []string{
		"water", "thunder", "flame", "wind", "stone",
		"mist", "serpent", "insect", "sound", "moon",
		"sun", "beast", "flower", "love",
	}

	for _, keyword := range breathingKeywords {
		if strings.Contains(descLower, keyword+" breath") {
			// Capitaliser la première lettre
			return strings.Title(keyword) + " Breathing"
		}
	}

	// 3. Chercher "Demon" pour les démons
	if strings.Contains(descLower, "demon") {
		if strings.Contains(descLower, "blood") {
			return "Blood Demon Art"
		}
		if strings.Contains(descLower, "art") {
			return "Demon Blood Art"
		}
	}

	// 4. Chercher n'importe quel mot suivi de "Breathing"
	words := strings.Fields(description)
	for i := 0; i < len(words)-1; i++ {
		if strings.ToLower(words[i+1]) == "breathing" {
			return words[i] + " Breathing"
		}
	}

	return "Aucun style de combat connu"
}
