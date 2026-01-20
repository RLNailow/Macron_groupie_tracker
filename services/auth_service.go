package services

import (
	"API-demon-slayyyyy-/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	UsersFilePath = "data/users.json" // Chemin vers le fichier JSON des utilisateurs
)

var (
	ErrUserNotFound      = errors.New("utilisateur non trouvé")
	ErrInvalidPassword   = errors.New("mot de passe invalide")
	ErrUserAlreadyExists = errors.New("utilisateur déjà existant")
	ErrInvalidEmail      = errors.New("email invalide")
)

// AuthService gère l'authentification et les utilisateurs
type AuthService struct {
	mu        sync.RWMutex     // Mutex pour la sécurité concurrentielle
	usersData models.UsersData // Données des utilisateurs
}

var authServiceInstance *AuthService

// InitAuthService initialise le service d'authentification
func InitAuthService() error {
	authServiceInstance = &AuthService{}

	// Créer le dossier data/ s'il n'existe pas
	if err := os.MkdirAll("data", 0755); err != nil {
		return fmt.Errorf("erreur création dossier data: %w", err)
	}

	// Charger ou créer le fichier users.json
	if err := authServiceInstance.loadUsers(); err != nil {
		// Si le fichier n'existe pas, créer un fichier vide
		authServiceInstance.usersData = models.UsersData{
			Users:  []models.User{},
			NextID: 1,
		}
		if err := authServiceInstance.saveUsers(); err != nil {
			return fmt.Errorf("erreur création fichier users: %w", err)
		}
	}

	log.Println("Auth Service initialisé")
	return nil
}

// GetAuthService retourne l'instance du service d'authentification
func GetAuthService() *AuthService {
	return authServiceInstance
}

// Register crée un nouveau compte utilisateur
func (s *AuthService) Register(email, password string) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Vérifier si l'email existe déjà
	for _, user := range s.usersData.Users {
		if user.Email == email {
			return nil, ErrUserAlreadyExists
		}
	}

	// Créer le nouvel utilisateur (mot de passe en clair car local)
	newUser := models.User{
		ID:        s.usersData.NextID,
		Email:     email,
		Password:  password, // ← Stocké en clair (OK pour du local)
		CreatedAt: time.Now(),
		Favorites: models.Favorites{
			Characters:   []int{},
			Quotes:       []int{},
			CombatStyles: []string{},
			Races:        []string{},
		},
	}

	// Ajouter l'utilisateur
	s.usersData.Users = append(s.usersData.Users, newUser)
	s.usersData.NextID++

	// Sauvegarder dans le fichier
	if err := s.saveUsers(); err != nil {
		return nil, fmt.Errorf("erreur sauvegarde user: %w", err)
	}

	log.Printf("✅ Nouvel utilisateur créé: %s (ID: %d)", email, newUser.ID)
	return &newUser, nil
}

// Login vérifie les identifiants et retourne l'utilisateur
func (s *AuthService) Login(email, password string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Chercher l'utilisateur par email
	for _, user := range s.usersData.Users {
		if user.Email == email {
			// Vérifier le mot de passe (comparaison simple en clair)
			if user.Password != password {
				return nil, ErrInvalidPassword
			}
			log.Printf("Connexion réussie: %s", email)
			return &user, nil
		}
	}

	return nil, ErrUserNotFound
}

// GetUserByID récupère un utilisateur par son ID
func (s *AuthService) GetUserByID(id int) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.usersData.Users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, ErrUserNotFound
}

// GetUserByEmail récupère un utilisateur par son email
func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.usersData.Users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, ErrUserNotFound
}

// loadUsers charge les utilisateurs depuis le fichier JSON
func (s *AuthService) loadUsers() error {
	data, err := os.ReadFile(UsersFilePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.usersData)
}

// saveUsers sauvegarde les utilisateurs dans le fichier JSON
func (s *AuthService) saveUsers() error {
	data, err := json.MarshalIndent(s.usersData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(UsersFilePath, data, 0644)
}
