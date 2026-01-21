package handlers

import (
	"API-demon-slayyyyy-/services"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// FavoritesHandler affiche la page des favoris de l'utilisateur
func FavoritesHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'utilisateur connecté
	user := getUserFromCookie(r)
	if user == nil {
		http.Redirect(w, r, "/?error=login_required", http.StatusSeeOther)
		return
	}

	// Récupérer les détails des favoris
	apiService := services.GetAPIService()
	authService := services.GetAuthService()

	// Cast de l'utilisateur
	userEmail := ""
	switch v := user.(type) {
	case map[string]interface{}:
		if email, ok := v["email"].(string); ok {
			userEmail = email
		}
	}

	userObj, err := authService.GetUserByEmail(userEmail)
	if err != nil {
		log.Printf("Erreur récupération user: %v", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

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
		User                 interface{}
		FavoriteCharacters   []interface{}
		FavoriteQuotes       []interface{}
		FavoriteCombatStyles []string
		FavoriteRaces        []string
	}{
		PageTitle:            "Mes Favoris - Demon Slayer",
		User:                 user,
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
	user := getUserFromCookie(r)
	if user == nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Non connecté",
		})
		return
	}

	// Récupérer les paramètres
	vars := mux.Vars(r)
	favoriteType := vars["type"] // "character", "quote", "combat_style", "race"
	valueStr := vars["value"]

	// Récupérer l'utilisateur
	authService := services.GetAuthService()
	userEmail := ""
	switch v := user.(type) {
	case map[string]interface{}:
		if email, ok := v["email"].(string); ok {
			userEmail = email
		}
	}

	userObj, err := authService.GetUserByEmail(userEmail)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Erreur utilisateur",
		})
		return
	}

	// Ajouter aux favoris selon le type
	added := false
	switch favoriteType {
	case "character", "quote":
		id, err := strconv.Atoi(valueStr)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{
				"error": "ID invalide",
			})
			return
		}

		if favoriteType == "character" {
			// Vérifier si déjà en favoris
			alreadyExists := false
			for _, fav := range userObj.Favorites.Characters {
				if fav == id {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				userObj.Favorites.Characters = append(userObj.Favorites.Characters, id)
				added = true
			}
		} else {
			// Quote
			alreadyExists := false
			for _, fav := range userObj.Favorites.Quotes {
				if fav == id {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				userObj.Favorites.Quotes = append(userObj.Favorites.Quotes, id)
				added = true
			}
		}

	case "combat_style", "race":
		if favoriteType == "combat_style" {
			// Vérifier si déjà en favoris
			alreadyExists := false
			for _, fav := range userObj.Favorites.CombatStyles {
				if fav == valueStr {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				userObj.Favorites.CombatStyles = append(userObj.Favorites.CombatStyles, valueStr)
				added = true
			}
		} else {
			// Race
			alreadyExists := false
			for _, fav := range userObj.Favorites.Races {
				if fav == valueStr {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				userObj.Favorites.Races = append(userObj.Favorites.Races, valueStr)
				added = true
			}
		}
	}

	// Sauvegarder les modifications
	if added {
		if err := authService.UpdateUser(userObj); err != nil {
			log.Printf("Erreur sauvegarde favoris: %v", err)
			respondJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Erreur de sauvegarde",
			})
			return
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"added":   added,
		"message": "Ajouté aux favoris",
	})
}

// RemoveFavoriteHandler retire un élément des favoris
func RemoveFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'utilisateur connecté
	user := getUserFromCookie(r)
	if user == nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Non connecté",
		})
		return
	}

	// Récupérer les paramètres
	vars := mux.Vars(r)
	favoriteType := vars["type"]
	valueStr := vars["value"]

	// Récupérer l'utilisateur
	authService := services.GetAuthService()
	userEmail := ""
	switch v := user.(type) {
	case map[string]interface{}:
		if email, ok := v["email"].(string); ok {
			userEmail = email
		}
	}

	userObj, err := authService.GetUserByEmail(userEmail)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Erreur utilisateur",
		})
		return
	}

	// Retirer des favoris selon le type
	removed := false
	switch favoriteType {
	case "character", "quote":
		id, err := strconv.Atoi(valueStr)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{
				"error": "ID invalide",
			})
			return
		}

		if favoriteType == "character" {
			newFavorites := []int{}
			for _, fav := range userObj.Favorites.Characters {
				if fav != id {
					newFavorites = append(newFavorites, fav)
				} else {
					removed = true
				}
			}
			userObj.Favorites.Characters = newFavorites
		} else {
			newFavorites := []int{}
			for _, fav := range userObj.Favorites.Quotes {
				if fav != id {
					newFavorites = append(newFavorites, fav)
				} else {
					removed = true
				}
			}
			userObj.Favorites.Quotes = newFavorites
		}

	case "combat_style", "race":
		if favoriteType == "combat_style" {
			newFavorites := []string{}
			for _, fav := range userObj.Favorites.CombatStyles {
				if fav != valueStr {
					newFavorites = append(newFavorites, fav)
				} else {
					removed = true
				}
			}
			userObj.Favorites.CombatStyles = newFavorites
		} else {
			newFavorites := []string{}
			for _, fav := range userObj.Favorites.Races {
				if fav != valueStr {
					newFavorites = append(newFavorites, fav)
				} else {
					removed = true
				}
			}
			userObj.Favorites.Races = newFavorites
		}
	}

	// Sauvegarder les modifications
	if removed {
		if err := authService.UpdateUser(userObj); err != nil {
			log.Printf("Erreur sauvegarde favoris: %v", err)
			respondJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Erreur de sauvegarde",
			})
			return
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"removed": removed,
		"message": "Retiré des favoris",
	})
}
