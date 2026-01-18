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
	url := fmt.Sprintf("%s/characters/%d", BaseURL, id)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur requête API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("personnage non trouvé")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture réponse: %w", err)
	}

	var character models.Character
	if err := json.Unmarshal(body, &character); err != nil {
		return nil, fmt.Errorf("erreur parsing JSON: %w", err)
	}

	return &character, nil
}

// GetBreathingTechniques récupère tous les styles de combat
func (s *APIService) GetBreathingTechniques() ([]models.BreathingTechnique, error) {
	url := fmt.Sprintf("%s/breathing-techniques", BaseURL)

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

	var apiResp models.BreathingTechniquesResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("erreur parsing JSON: %w", err)
	}

	return apiResp.Content, nil
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
