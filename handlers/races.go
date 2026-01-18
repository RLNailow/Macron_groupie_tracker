package handlers

import (
	"API-demon-slayyyyy-/services"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// RacesHandler affiche la liste des races (paginée)
func RacesHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le numéro de page
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Récupérer toutes les races uniques
	apiService := services.GetAPIService()
	allRaces, err := apiService.GetUniqueRaces()
	if err != nil {
		log.Printf("❌ Erreur récupération races: %v", err)
		http.Error(w, "Erreur chargement des races", http.StatusInternalServerError)
		return
	}

	// Calculer la pagination
	totalRaces := len(allRaces)
	totalPages := (totalRaces + ItemsPerPage - 1) / ItemsPerPage

	if page > totalPages {
		page = totalPages
	}

	// Extraire les races de la page actuelle
	startIndex := (page - 1) * ItemsPerPage
	endIndex := startIndex + ItemsPerPage
	if endIndex > totalRaces {
		endIndex = totalRaces
	}

	racesOnPage := allRaces[startIndex:endIndex]

	// Charger le template
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/races.html")
	if err != nil {
		log.Printf("❌ Erreur chargement template races: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Données à passer au template
	data := struct {
		PageTitle    string
		User         interface{}
		Races        []string
		CurrentPage  int
		TotalPages   int
		HasPrevious  bool
		HasNext      bool
		PreviousPage int
		NextPage     int
	}{
		PageTitle:    "Races - Demon Slayer",
		User:         getUserFromCookie(r),
		Races:        racesOnPage,
		CurrentPage:  page,
		TotalPages:   totalPages,
		HasPrevious:  page > 1,
		HasNext:      page < totalPages,
		PreviousPage: page - 1,
		NextPage:     page + 1,
	}

	// Rendre le template
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("❌ Erreur rendu template races: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}
