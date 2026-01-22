package handlers

import (
	"API-demon-slayyyyy-/services"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CombatStylesHandler affiche la liste des styles de combat (paginée)
func CombatStylesHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le numéro de page
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Récupérer tous les styles depuis l'API
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

	if page > totalPages && totalStyles > 0 {
		page = totalPages
	}

	// Extraire les styles de la page actuelle
	var stylesOnPage []interface{}
	if totalStyles > 0 {
		startIndex := (page - 1) * ItemsPerPage
		endIndex := startIndex + ItemsPerPage
		if endIndex > totalStyles {
			endIndex = totalStyles
		}
		for i := startIndex; i < endIndex; i++ {
			stylesOnPage = append(stylesOnPage, allStyles[i])
		}
	}

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

// CombatStyleDetailHandler affiche la page détail d'un style de combat
func CombatStyleDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID depuis l'URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ID invalide: %s", idStr)
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	// Récupérer le style de combat depuis l'API
	apiService := services.GetAPIService()
	style, err := apiService.GetBreathingTechniqueByID(id)
	if err != nil {
		log.Printf("Style de combat %d non trouvé: %v", id, err)
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	// Debug: afficher le nombre de personnages
	log.Printf("Style '%s' a %d personnages", style.Name, len(style.Characters))

	// Charger le template
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/combat_style_detail.html")
	if err != nil {
		log.Printf("Erreur chargement template combat_style_detail: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Données à passer au template
	data := struct {
		PageTitle string
		User      interface{}
		Style     interface{}
	}{
		PageTitle: style.Name + " - Demon Slayer",
		User:      getUserFromCookie(r),
		Style:     style,
	}

	// Rendre le template
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Erreur rendu template combat_style_detail: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}
