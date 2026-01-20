package services

import (
	"API-demon-slayyyyy-/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	// URL de base de l'API Demon Slayer
	BaseURL = "https://www.demonslayer-api.com/api/v1"
)

// APIService gère toutes les requêtes vers l'API Demon Slayer
type APIService struct {
	client *http.Client
}

var apiServiceInstance *APIService

// InitAPIService initialise le service API
func InitAPIService() {
	apiServiceInstance = &APIService{
		client: &http.Client{
			Timeout: 10 * time.Second, // Timeout de 10 secondes
		},
	}
	log.Println("✅ API Service initialisé")
}

// GetAPIService retourne l'instance du service API
func GetAPIService() *APIService {
	return apiServiceInstance
}

// GetAllCharacters récupère tous les personnages (limité à 100)
func (s *APIService) GetAllCharacters() ([]models.Character, error) {
	url := fmt.Sprintf("%s/characters?limit=100", BaseURL)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur requête API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture réponse: %w", err)
	}

	var apiResp models.APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("erreur parsing JSON: %w", err)
	}

	return apiResp.Content, nil
}

// GetCharacterByID récupère un personnage spécifique par son ID
func (s *APIService) GetCharacterByID(id int) (*models.Character, error) {
	// Récupérer tous les personnages et chercher par ID
	// (Plus fiable que l'endpoint /characters/{id} qui semble bugué)
	characters, err := s.GetAllCharacters()
	if err != nil {
		return nil, err
	}

	// Chercher le personnage par ID
	for _, char := range characters {
		if char.ID == id {
			return &char, nil
		}
	}

	return nil, fmt.Errorf("personnage non trouvé")
}

// GetBreathingTechniques extrait les styles de combat depuis les descriptions des personnages
func (s *APIService) GetBreathingTechniques() ([]models.BreathingTechnique, error) {
	// Récupérer tous les personnages
	characters, err := s.GetAllCharacters()
	if err != nil {
		return nil, err
	}

	// Map pour stocker les styles uniques
	stylesMap := make(map[string]bool)

	// Extraire les styles depuis les descriptions
	for _, char := range characters {
		if char.Description != "" {
			// Chercher "X Breathing" dans la description
			// Patterns courants : "Water Breathing", "Thunder Breathing", etc.
			description := char.Description

			// Liste des styles connus
			knownStyles := []string{
				"Water Breathing", "Thunder Breathing", "Flame Breathing",
				"Wind Breathing", "Stone Breathing", "Mist Breathing",
				"Serpent Breathing", "Insect Breathing", "Sound Breathing",
				"Moon Breathing", "Sun Breathing", "Beast Breathing",
				"Flower Breathing", "Love Breathing",
			}

			for _, style := range knownStyles {
				if contains(description, style) {
					stylesMap[style] = true
				}
			}
		}
	}

	// Convertir en slice
	var techniques []models.BreathingTechnique
	id := 1
	for style := range stylesMap {
		techniques = append(techniques, models.BreathingTechnique{
			ID:          id,
			Name:        style,
			Description: "",
		})
		id++
	}

	return techniques, nil
}

// contains vérifie si une chaîne contient une sous-chaîne (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// GetUniqueRaces extrait toutes les races uniques depuis les personnages
func (s *APIService) GetUniqueRaces() ([]string, error) {
	characters, err := s.GetAllCharacters()
	if err != nil {
		return nil, err
	}

	// Map pour stocker les races uniques
	racesMap := make(map[string]bool)
	for _, char := range characters {
		if char.Race != "" {
			racesMap[char.Race] = true
		}
	}

	// Convertir la map en slice
	races := make([]string, 0, len(racesMap))
	for race := range racesMap {
		races = append(races, race)
	}

	return races, nil
}

// GetCharactersWithQuotes récupère tous les personnages ayant une citation
func (s *APIService) GetCharactersWithQuotes() ([]models.Character, error) {
	characters, err := s.GetAllCharacters()
	if err != nil {
		return nil, err
	}

	// Filtrer les personnages ayant une citation
	var withQuotes []models.Character
	for _, char := range characters {
		if char.Quote != "" {
			withQuotes = append(withQuotes, char)
		}
	}

	return withQuotes, nil
}
