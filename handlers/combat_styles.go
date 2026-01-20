package handlers

import (
	"API-demon-slayyyyy-/services"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// CombatStylesHandler affiche la liste des styles de combat (paginée)
func CombatStylesHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le numéro de page
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Récupérer tous les styles de combat depuis l'API
	apiService := services.GetAPIService()
	allStyles, err := apiService.GetBreathingTechniques()
	if err != nil {
		log.Printf("Erreur récupération styles de combat: %v", err)
		http.Error(w, "Erreur chargement des styles de combat", http.StatusInternalServerError)
		return
	}

	// Calculer la pagination
	totalStyles := len(allStyles)
	totalPages := (totalStyles + ItemsPerPage - 1) / ItemsPerPage

	if page > totalPages {
		page = totalPages
	}

	// Extraire les styles de la page actuelle
	startIndex := (page - 1) * ItemsPerPage
	endIndex := startIndex + ItemsPerPage
	if endIndex > totalStyles {
		endIndex = totalStyles
	}

	stylesOnPage := allStyles[startIndex:endIndex]

	// Charger le template
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/combat_styles.html")
	if err != nil {
		log.Printf("Erreur chargement template combat_styles: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Données à passer au template
	data := struct {
		PageTitle    string
		User         interface{}
		Styles       interface{}
		CurrentPage  int
		TotalPages   int
		HasPrevious  bool
		HasNext      bool
		PreviousPage int
		NextPage     int
	}{
		PageTitle:    "Styles de combat - Demon Slayer",
		User:         getUserFromCookie(r),
		Styles:       stylesOnPage,
		CurrentPage:  page,
		TotalPages:   totalPages,
		HasPrevious:  page > 1,
		HasNext:      page < totalPages,
		PreviousPage: page - 1,
		NextPage:     page + 1,
	}

	// Rendre le template
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Erreur rendu template combat_styles: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}
