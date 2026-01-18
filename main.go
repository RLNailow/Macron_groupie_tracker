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
	router.HandleFunc("/characters/{id:[0-9]+}", handlers.CharacterDetailHandler).Methods("GET")
	router.HandleFunc("/races", handlers.RacesHandler).Methods("GET")
	router.HandleFunc("/combat-styles", handlers.CombatStylesHandler).Methods("GET")
	router.HandleFunc("/quotes", handlers.QuotesHandler).Methods("GET")
	router.HandleFunc("/random", handlers.RandomHandler).Methods("GET")

	// ==================== ROUTES PROTÉGÉES ====================
	// (Si tu veux ajouter des favoris plus tard, on les protège ici)

	// Exemple : router.Handle("/favorites", middleware.AuthRequired(http.HandlerFunc(handlers.FavoritesHandler))).Methods("GET")

	// ==================== FICHIERS STATIQUES ====================

	// Servir les fichiers CSS, JS, images
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// ==================== GESTION 404 ====================

	// Toute route non trouvée affiche la page 404
	router.NotFoundHandler = http.HandlerFunc(handlers.ErrorHandler)

	// ==================== DÉMARRAGE DU SERVEUR ====================

	log.Println("🔥 Serveur démarré sur http://localhost:8080")
	log.Println("📁 Fichiers statiques servis depuis /static/")
	log.Println("🎨 Templates chargés depuis /templates/")

	// Lancer le serveur sur le port 8080
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Erreur démarrage serveur:", err)
	}
}
