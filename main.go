package main

import (
	"API-demon-slayyyyy-/handlers"
	"API-demon-slayyyyy-/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialiser le service API (pour appeler l'API Demon Slayer)
	services.InitAPIService()

	// Initialiser le service d'authentification (gestion des utilisateurs)
	if err := services.InitAuthService(); err != nil {
		log.Fatal("Erreur initialisation auth service:", err)
	}

	// Créer le routeur Gorilla Mux
	router := mux.NewRouter()

	// ==================== ROUTES PUBLIQUES ====================

	// Page d'accueil
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	// Authentification
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/logout", handlers.LogoutHandler).Methods("GET")

	// Pages principales (accessibles sans login)
	router.HandleFunc("/characters", handlers.CharactersHandler).Methods("GET")
	router.HandleFunc("/characters/{id:[0-9]+}", handlers.DetailHandler).Methods("GET")
	router.HandleFunc("/races", handlers.RacesHandler).Methods("GET")
	router.HandleFunc("/combat-styles", handlers.CombatStylesHandler).Methods("GET")
	router.HandleFunc("/combat-styles/{id:[0-9]+}", handlers.CombatStyleDetailHandler).Methods("GET")
	router.HandleFunc("/quotes", handlers.QuotesHandler).Methods("GET")
	router.HandleFunc("/random", handlers.RandomHandler).Methods("GET")

	// API interne pour obtenir un ID aléatoire
	router.HandleFunc("/api/random-character", handlers.GetRandomCharacterID).Methods("GET")

	// API pour la recherche (évite les problèmes CORS)
	router.HandleFunc("/api/characters", handlers.APICharactersHandler).Methods("GET")
	router.HandleFunc("/api/combat-styles", handlers.APICombatStylesHandler).Methods("GET")

	// ==================== ROUTES PROTÉGÉES ====================

	// Favoris (nécessite d'être connecté)
	router.HandleFunc("/favorites", handlers.FavoritesHandler).Methods("GET")
	router.HandleFunc("/favorites/add/{type}/{value}", handlers.AddFavoriteHandler).Methods("POST")
	router.HandleFunc("/favorites/remove/{type}/{value}", handlers.RemoveFavoriteHandler).Methods("POST")

	// ==================== FICHIERS STATIQUES ====================

	// Servir les fichiers CSS, JS, images
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// ==================== GESTION 404 ====================

	// Toute route non trouvée affiche la page 404
	router.NotFoundHandler = http.HandlerFunc(handlers.ErrorHandler)

	// ==================== DÉMARRAGE DU SERVEUR ====================

	log.Println("Serveur démarré sur http://localhost:8080")

	// Lancer le serveur sur le port 8080
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Erreur démarrage serveur:", err)
	}
}
