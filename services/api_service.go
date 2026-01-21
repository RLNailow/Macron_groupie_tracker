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

// GetBreathingTechniques récupère tous les styles de combat depuis l'API
func (s *APIService) GetBreathingTechniques() ([]models.BreathingTechnique, error) {
	url := fmt.Sprintf("%s/combat-styles?limit=100", BaseURL)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur requête API combat-styles: %w", err)
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

	log.Printf("✅ %d styles de combat récupérés depuis l'API", len(apiResp.Content))

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

// StyleWithCharacter associe un style à un personnage
type StyleWithCharacter struct {
	ID          int
	Name        string
	CharacterID int
}

// GetStylesWithCharacters récupère tous les styles avec le personnage qui l'utilise
func (s *APIService) GetStylesWithCharacters() ([]StyleWithCharacter, error) {
	// Récupérer tous les styles
	styles, err := s.GetBreathingTechniques()
	if err != nil {
		return nil, err
	}

	// Récupérer tous les personnages
	characters, err := s.GetAllCharacters()
	if err != nil {
		return nil, err
	}

	var result []StyleWithCharacter

	// Pour chaque style, trouver le premier personnage qui l'utilise
	for _, style := range styles {
		found := false

		// Chercher dans les personnages
		for _, char := range characters {
			// Vérifier avec plusieurs stratégies
			if matchesStyle(char.Description, char.Name, style.Name) {
				result = append(result, StyleWithCharacter{
					ID:          style.ID,
					Name:        style.Name,
					CharacterID: char.ID,
				})
				log.Printf("✅ Style '%s' → %s (ID: %d)", style.Name, char.Name, char.ID)
				found = true
				break
			}
		}

		// Si aucun personnage trouvé, utiliser le premier personnage (Tanjiro)
		if !found && len(characters) > 0 {
			result = append(result, StyleWithCharacter{
				ID:          style.ID,
				Name:        style.Name,
				CharacterID: characters[0].ID,
			})
			log.Printf("⚠️  Style '%s' → Tanjiro par défaut", style.Name)
		}
	}

	log.Printf("📊 %d styles associés à des personnages", len(result))
	return result, nil
}

// matchesStyle vérifie si un personnage correspond à un style avec plusieurs stratégies
func matchesStyle(description, characterName, styleName string) bool {
	desc := toLower(description)
	name := toLower(characterName)
	style := toLower(styleName)

	// 1. Recherche exacte du nom du style dans la description
	if contains(desc, style) {
		return true
	}

	// 2. Extraire le mot-clé principal du style
	var mainKeyword string
	if contains(style, "breathing") {
		// "Water Breathing" → "water"
		parts := split(style, " breathing")
		if len(parts) > 0 {
			mainKeyword = trim(parts[0])
		}
	} else if contains(style, "manipulation") {
		// "Blood Manipulation" → "blood"
		parts := split(style, " manipulation")
		if len(parts) > 0 {
			mainKeyword = trim(parts[0])
		}
	} else if contains(style, "demon art") {
		mainKeyword = "demon"
	}

	// 3. Chercher le mot-clé dans description ou nom
	if mainKeyword != "" && len(mainKeyword) > 3 {
		if contains(desc, mainKeyword) || contains(name, mainKeyword) {
			return true
		}
	}

	// 4. Mappings spécifiques personnage → style
	mappings := map[string][]string{
		"tanjiro":     {"sun breathing", "water breathing", "hinokami"},
		"giyu":        {"water breathing"},
		"zenitsu":     {"thunder breathing"},
		"inosuke":     {"beast breathing"},
		"shinobu":     {"insect breathing"},
		"kyojuro":     {"flame breathing"},
		"rengoku":     {"flame breathing"},
		"tengen":      {"sound breathing"},
		"uzui":        {"sound breathing"},
		"mitsuri":     {"love breathing"},
		"kanroji":     {"love breathing"},
		"muichiro":    {"mist breathing"},
		"tokito":      {"mist breathing"},
		"gyomei":      {"stone breathing"},
		"himejima":    {"stone breathing"},
		"sanemi":      {"wind breathing"},
		"shinazugawa": {"wind breathing"},
		"obanai":      {"serpent breathing"},
		"iguro":       {"serpent breathing"},
		"kanae":       {"flower breathing"},
		"kanao":       {"flower breathing"},
		"yoriichi":    {"sun breathing"},
		"kokushibo":   {"moon breathing"},
		"muzan":       {"blood", "demon"},
		"akaza":       {"destructive death", "demon"},
		"douma":       {"cryokinesis", "ice", "demon"},
		"gyutaro":     {"blood", "sickle", "demon"},
		"daki":        {"obi", "sash", "demon"},
		"enmu":        {"sleep", "dream", "demon"},
		"kaigaku":     {"thunder breathing"},
	}

	for charKeyword, styleKeywords := range mappings {
		if contains(name, charKeyword) {
			for _, styleKeyword := range styleKeywords {
				if contains(style, styleKeyword) {
					return true
				}
			}
		}
	}

	return false
}

// Fonctions utilitaires sans dépendances externes
func toLower(s string) string {
	result := ""
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			result += string(c + 32)
		} else {
			result += string(c)
		}
	}
	return result
}

func contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func split(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		match := true
		for j := 0; j < len(sep); j++ {
			if s[i+j] != sep[j] {
				match = false
				break
			}
		}
		if match {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trim(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}
	return s[start:end]
}
