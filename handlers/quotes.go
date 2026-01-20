package handlers

import (
	"API-demon-slayyyyy-/services"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// QuotesHandler affiche la liste des citations (paginée)
func QuotesHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le numéro de page
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Récupérer tous les personnages ayant une citation
	apiService := services.GetAPIService()
	charactersWithQuotes, err := apiService.GetCharactersWithQuotes()
	if err != nil {
		log.Printf("Erreur récupération citations: %v", err)
		http.Error(w, "Erreur chargement des citations", http.StatusInternalServerError)
		return
	}

	// Calculer la pagination
	totalQuotes := len(charactersWithQuotes)
	totalPages := (totalQuotes + ItemsPerPage - 1) / ItemsPerPage

	if page > totalPages {
		page = totalPages
	}

	// Extraire les citations de la page actuelle
	startIndex := (page - 1) * ItemsPerPage
	endIndex := startIndex + ItemsPerPage
	if endIndex > totalQuotes {
		endIndex = totalQuotes
	}

	quotesOnPage := charactersWithQuotes[startIndex:endIndex]

	// Charger le template
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/quotes.html")
	if err != nil {
		log.Printf("Erreur chargement template quotes: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Données à passer au template
	data := struct {
		PageTitle    string
		User         interface{}
		Characters   interface{}
		CurrentPage  int
		TotalPages   int
		HasPrevious  bool
		HasNext      bool
		PreviousPage int
		NextPage     int
	}{
		PageTitle:    "Citations - Demon Slayer",
		User:         getUserFromCookie(r),
		Characters:   quotesOnPage,
		CurrentPage:  page,
		TotalPages:   totalPages,
		HasPrevious:  page > 1,
		HasNext:      page < totalPages,
		PreviousPage: page - 1,
		NextPage:     page + 1,
	}

	// Rendre le template
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Erreur rendu template quotes: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}
