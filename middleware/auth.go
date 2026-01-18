package middleware

import (
	"API-demon-slayyyyy-/services"
	"context"
	"net/http"
)

type contextKey string

const UserContextKey contextKey = "user"

// AuthRequired est un middleware qui vérifie si l'utilisateur est connecté
// (Pour l'instant non utilisé, mais prêt pour les favoris plus tard)
func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Récupérer le cookie de session
		cookie, err := r.Cookie("session_token")
		if err != nil {
			// Pas de cookie = pas connecté
			http.Redirect(w, r, "/?error=login_required", http.StatusSeeOther)
			return
		}

		// Vérifier si le token est valide
		// (Pour l'instant, on utilise l'email comme token simple)
		// Dans une vraie app, tu utiliserais JWT ou des sessions Redis
		authService := services.GetAuthService()
		user, err := authService.GetUserByEmail(cookie.Value)
		if err != nil {
			// Token invalide
			http.Redirect(w, r, "/?error=invalid_session", http.StatusSeeOther)
			return
		}

		// Ajouter l'utilisateur au contexte de la requête
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext récupère l'utilisateur depuis le contexte
func GetUserFromContext(r *http.Request) interface{} {
	return r.Context().Value(UserContextKey)
}
