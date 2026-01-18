package models

import "time"

// User représente un utilisateur du site
type User struct {
	ID        int       `json:"id"`         // Identifiant unique
	Email     string    `json:"email"`      // Email de l'utilisateur
	Password  string    `json:"password"`   // Mot de passe hashé avec bcrypt
	CreatedAt time.Time `json:"created_at"` // Date de création du compte
	Favorites Favorites `json:"favorites"`  // Favoris de l'utilisateur (pour plus tard)
}

// Favorites représente les éléments mis en favoris par un utilisateur
type Favorites struct {
	Characters   []int    `json:"characters"`    // IDs des personnages favoris
	Quotes       []int    `json:"quotes"`        // IDs des citations favorites
	CombatStyles []string `json:"combat_styles"` // Noms des styles de combat favoris
	Races        []string `json:"races"`         // Noms des races favorites
}

// UserCredentials représente les données de connexion
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UsersData représente la structure du fichier users.json
type UsersData struct {
	Users  []User `json:"users"`   // Liste de tous les utilisateurs
	NextID int    `json:"next_id"` // Prochain ID disponible
}
