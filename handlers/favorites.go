package handlers

import (
	"API-demon-slayyyyy-/models"
	"API-demon-slayyyyy-/services"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Helper: vérifie si un ID existe dans une slice
func containsInt(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// Helper: vérifie si une string existe dans une slice
func containsString(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// Helper: retire un ID d'une slice
func removeInt(slice []int, val int) []int {
	result := []int{}
	for _, item := range slice {
		if item != val {
			result = append(result, item)
		}
	}
	return result
}

// Helper: retire une string d'une slice
func removeString(slice []string, val string) []string {
	result := []string{}
	for _, item := range slice {
		if item != val {
			result = append(result, item)
		}
	}
	return result
}

// FavoritesHandler affiche la page des favoris de l'utilisateur
func FavoritesHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'utilisateur connecté
	userObj := getUserFromCookie(r)
	if userObj == nil {
		http.Redirect(w, r, "/?error=login_required", http.StatusSeeOther)
		return
	}

	// Récupérer les détails des favoris
	apiService := services.GetAPIService()

	// Récupérer les personnages favoris
	var favoriteCharacters []interface{}
	for _, charID := range userObj.Favorites.Characters {
		char, err := apiService.GetCharacterByID(charID)
		if err == nil {
			favoriteCharacters = append(favoriteCharacters, char)
		}
	}

	// Récupérer tous les personnages pour les citations favorites
	allCharacters, _ := apiService.GetAllCharacters()
	var favoriteQuotes []interface{}
	for _, quoteID := range userObj.Favorites.Quotes {
		for _, char := range allCharacters {
			if char.ID == quoteID {
				favoriteQuotes = append(favoriteQuotes, char)
				break
			}
		}
	}

	// Charger le template
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/favorites.html")
	if err != nil {
		log.Printf("Erreur chargement template favorites: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Données à passer au template
	data := struct {
		PageTitle            string
		User                 *models.User
		FavoriteCharacters   []interface{}
		FavoriteQuotes       []interface{}
		FavoriteCombatStyles []string
		FavoriteRaces        []string
	}{
		PageTitle:            "Mes Favoris - Demon Slayer",
		User:                 userObj,
		FavoriteCharacters:   favoriteCharacters,
		FavoriteQuotes:       favoriteQuotes,
		FavoriteCombatStyles: userObj.Favorites.CombatStyles,
		FavoriteRaces:        userObj.Favorites.Races,
	}

	// Rendre le template
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Erreur rendu template favorites: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
}

// AddFavoriteHandler ajoute un élément aux favoris
func AddFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'utilisateur connecté
	userObj := getUserFromCookie(r)
	if userObj == nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Non connecté"})
		return
	}

	// Récupérer les paramètres
	vars := mux.Vars(r)
	favoriteType := vars["type"] // "character", "quote", "combat_style", "race"
	valueStr := vars["value"]

	// Ajouter aux favoris selon le type
	added := false
	switch favoriteType {
	case "character", "quote":
		id, err := strconv.Atoi(valueStr)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "ID invalide"})
			return
		}
		if favoriteType == "character" && !containsInt(userObj.Favorites.Characters, id) {
			userObj.Favorites.Characters = append(userObj.Favorites.Characters, id)
			added = true
		} else if favoriteType == "quote" && !containsInt(userObj.Favorites.Quotes, id) {
			userObj.Favorites.Quotes = append(userObj.Favorites.Quotes, id)
			added = true
		}
	case "combat_style":
		if !containsString(userObj.Favorites.CombatStyles, valueStr) {
			userObj.Favorites.CombatStyles = append(userObj.Favorites.CombatStyles, valueStr)
			added = true
		}
	case "race":
		if !containsString(userObj.Favorites.Races, valueStr) {
			userObj.Favorites.Races = append(userObj.Favorites.Races, valueStr)
			added = true
		}
	}

	// Sauvegarder les modifications
	if added {
		if err := services.GetAuthService().UpdateUser(userObj); err != nil {
			log.Printf("Erreur sauvegarde favoris: %v", err)
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erreur de sauvegarde"})
			return
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "added": added, "message": "Ajouté aux favoris"})
}

// RemoveFavoriteHandler retire un élément des favoris
func RemoveFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'utilisateur connecté
	userObj := getUserFromCookie(r)
	if userObj == nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Non connecté"})
		return
	}

	// Récupérer les paramètres
	vars := mux.Vars(r)
	favoriteType := vars["type"]
	valueStr := vars["value"]

	// Retirer des favoris selon le type
	removed := false
	switch favoriteType {
	case "character", "quote":
		id, err := strconv.Atoi(valueStr)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "ID invalide"})
			return
		}
		if favoriteType == "character" && containsInt(userObj.Favorites.Characters, id) {
			userObj.Favorites.Characters = removeInt(userObj.Favorites.Characters, id)
			removed = true
		} else if favoriteType == "quote" && containsInt(userObj.Favorites.Quotes, id) {
			userObj.Favorites.Quotes = removeInt(userObj.Favorites.Quotes, id)
			removed = true
		}
	case "combat_style":
		if containsString(userObj.Favorites.CombatStyles, valueStr) {
			userObj.Favorites.CombatStyles = removeString(userObj.Favorites.CombatStyles, valueStr)
			removed = true
		}
	case "race":
		if containsString(userObj.Favorites.Races, valueStr) {
			userObj.Favorites.Races = removeString(userObj.Favorites.Races, valueStr)
			removed = true
		}
	}

	// Sauvegarder les modifications
	if removed {
		if err := services.GetAuthService().UpdateUser(userObj); err != nil {
			log.Printf("Erreur sauvegarde favoris: %v", err)
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erreur de sauvegarde"})
			return
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "removed": removed, "message": "Retiré des favoris"})
}
