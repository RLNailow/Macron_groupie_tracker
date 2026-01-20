package handlers

import (
	"API-demon-slayyyyy-/services"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const ItemsPerPage = 10 // 10 personnages par page

// CharactersHandler affiche la liste des personnages (paginée)
func CharactersHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le numéro de page depuis l'URL (?page=1)
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Page par défaut
	}

	// Récupérer tous les personnages depuis l'API
	apiService := services.GetAPIService()
	allCharacters, err := apiService.GetAllCharacters()
	if err != nil {
		log.Printf("Erreur récupération personnages: %v", err)
		http.Error(w, "Erreur chargement des personnages", http.StatusInternalServerError)
		return
	}

	// Calculer la pagination
	totalCharacters := len(allCharacters)
	totalPages := (totalCharacters + ItemsPerPage - 1) / ItemsPerPage // Arrondi au supérieur

	// Vérifier que la page existe
	if page > totalPages {
		page = totalPages
	}

	// Calculer les indices de début et fin
	startIndex := (page - 1) * ItemsPerPage
	endIndex := startIndex + ItemsPerPage
	if endIndex > totalCharacters {
		endIndex = totalCharacters
	}

	// Extraire les personnages de la page actuelle
	charactersOnPage := allCharacters[startIndex:endIndex]

	// Charger le template
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/characters.html")
	if err != nil {
		log.Printf("Erreur chargement template characters: %v", err)
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
		PageTitle:    "Personnages - Demon Slayer",
		User:         getUserFromCookie(r),
		Characters:   charactersOnPage,
		CurrentPage:  page,
		TotalPages:   totalPages,
		HasPrevious:  page > 1,
		HasNext:      page < totalPages,
		PreviousPage: page - 1,
		NextPage:     page + 1,
	}

	// Rendre le template
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Erreur rendu template characters: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}
