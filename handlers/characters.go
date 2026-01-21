package handlers

import (
	"API-demon-slayyyyy-/services"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const ItemsPerPage = 10 // 10 personnages par page

// CharactersHandler affiche la liste des personnages (paginée et filtrée)
func CharactersHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le numéro de page depuis l'URL (?page=1)
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Page par défaut
	}

	// Récupérer le filtre de race (?race=Human ou ?race=Demon)
	raceFilter := r.URL.Query().Get("race")

	// Récupérer tous les personnages depuis l'API
	apiService := services.GetAPIService()
	allCharacters, err := apiService.GetAllCharacters()
	if err != nil {
		log.Printf("Erreur récupération personnages: %v", err)
		http.Error(w, "Erreur chargement des personnages", http.StatusInternalServerError)
		return
	}

	// Filtrer par race si un filtre est spécifié
	var filteredCharacters []interface{}
	if raceFilter != "" {
		for _, char := range allCharacters {
			if char.Race == raceFilter {
				filteredCharacters = append(filteredCharacters, char)
			}
		}
	} else {
		// Pas de filtre, afficher tous les personnages
		for _, char := range allCharacters {
			filteredCharacters = append(filteredCharacters, char)
		}
	}

	// Calculer la pagination
	totalCharacters := len(filteredCharacters)
	totalPages := (totalCharacters + ItemsPerPage - 1) / ItemsPerPage // Arrondi au supérieur

	// Vérifier que la page existe
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	// Calculer les indices de début et fin
	var charactersOnPage []interface{}
	if totalCharacters > 0 {
		startIndex := (page - 1) * ItemsPerPage
		endIndex := startIndex + ItemsPerPage
		if endIndex > totalCharacters {
			endIndex = totalCharacters
		}

		// Extraire les personnages de la page actuelle
		charactersOnPage = filteredCharacters[startIndex:endIndex]
	}

	// Charger le template
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/characters.html")
	if err != nil {
		log.Printf("Erreur chargement template characters: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Construire le titre selon le filtre
	pageTitle := "Personnages - Demon Slayer"
	if raceFilter != "" {
		pageTitle = raceFilter + "s - Demon Slayer"
	}

	// Construire les URLs de pagination
	var prevURL, nextURL string
	if raceFilter != "" {
		prevURL = fmt.Sprintf("/characters?race=%s&page=%d", raceFilter, page-1)
		nextURL = fmt.Sprintf("/characters?race=%s&page=%d", raceFilter, page+1)
	} else {
		prevURL = fmt.Sprintf("/characters?page=%d", page-1)
		nextURL = fmt.Sprintf("/characters?page=%d", page+1)
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
		PrevURL      string
		NextURL      string
	}{
		PageTitle:    pageTitle,
		User:         getUserFromCookie(r),
		Characters:   charactersOnPage,
		CurrentPage:  page,
		TotalPages:   totalPages,
		HasPrevious:  page > 1,
		HasNext:      page < totalPages,
		PreviousPage: page - 1,
		NextPage:     page + 1,
		PrevURL:      prevURL,
		NextURL:      nextURL,
	}

	// Rendre le template
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Erreur rendu template characters: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}
