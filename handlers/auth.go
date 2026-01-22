package handlers

import (
	"API-demon-slayyyyy-/models"
	"API-demon-slayyyyy-/services"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// LoginHandler traite la connexion
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Décoder les données JSON du body
	var creds models.UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Données invalides"})
		return
	}

	// Vérifier les identifiants
	authService := services.GetAuthService()
	user, err := authService.Login(creds.Email, creds.Password)
	if err != nil {
		log.Printf("Échec login: %s - %v", creds.Email, err)
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Email ou mot de passe invalide"})
		return
	}

	// Créer un cookie de session et répondre
	setSessionCookie(w, user.Email)
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Connexion réussie",
		"user":    map[string]interface{}{"id": user.ID, "email": user.Email},
	})
}

// RegisterHandler traite l'inscription
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Décoder les données JSON
	var creds models.UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Données invalides"})
		return
	}

	// Validation basique
	if len(creds.Email) < 3 || len(creds.Password) < 6 {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Email ou mot de passe trop court (min 6 caractères)"})
		return
	}

	// Créer le compte
	authService := services.GetAuthService()
	user, err := authService.Register(creds.Email, creds.Password)
	if err != nil {
		log.Printf("Échec register: %s - %v", creds.Email, err)
		respondJSON(w, http.StatusConflict, map[string]string{"error": "Cet email est déjà utilisé"})
		return
	}

	// Connecter automatiquement l'utilisateur et répondre
	setSessionCookie(w, user.Email)
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Compte créé avec succès",
		"user":    map[string]interface{}{"id": user.ID, "email": user.Email},
	})
}

// LogoutHandler déconnecte l'utilisateur
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Supprimer le cookie de session
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Expire dans le passé
		HttpOnly: true,
		Path:     "/",
	})

	// Rediriger vers la page d'accueil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// setSessionCookie crée un cookie de session
func setSessionCookie(w http.ResponseWriter, email string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    email,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
}

// getUserFromCookie vérifie si un utilisateur est connecté via son cookie
func getUserFromCookie(r *http.Request) *models.User {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}
	user, err := services.GetAuthService().GetUserByEmail(cookie.Value)
	if err != nil {
		return nil
	}
	return user
}

// respondJSON envoie une réponse JSON
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
