package models

// APIResponse représente la structure de réponse de l'API Demon Slayer
type APIResponse struct {
	Pagination Pagination  `json:"pagination"` // Informations de pagination
	Content    []Character `json:"content"`    // Liste des personnages
}

// Pagination contient les infos de pagination de l'API
type Pagination struct {
	TotalElements  int    `json:"totalElements"`  // Nombre total d'éléments
	ElementsOnPage int    `json:"elementsOnPage"` // Éléments sur cette page
	CurrentPage    int    `json:"currentPage"`    // Page actuelle
	TotalPages     int    `json:"totalPages"`     // Nombre total de pages
	PreviousPage   string `json:"previousPage"`   // URL page précédente
	NextPage       string `json:"nextPage"`       // URL page suivante
}

// BreathingTechniquesResponse représente la réponse pour les styles de combat
type BreathingTechniquesResponse struct {
	Pagination Pagination           `json:"pagination"`
	Content    []BreathingTechnique `json:"content"`
}
