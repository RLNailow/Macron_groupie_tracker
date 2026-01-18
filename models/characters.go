package models

// Character représente un personnage de Demon Slayer
type Character struct {
	ID            int    `json:"id"`             // ID unique du personnage
	Name          string `json:"name"`           // Nom du personnage
	Age           int    `json:"age"`            // Âge du personnage
	Gender        string `json:"gender"`         // Sexe (Male/Female)
	Race          string `json:"race"`           // Race (Human/Demon/Hybrid)
	Description   string `json:"description"`    // Description du personnage
	Img           string `json:"img"`            // URL de l'image du personnage
	AffiliationID int    `json:"affiliation_id"` // ID de l'affiliation
	ArcID         int    `json:"arc_id"`         // ID de l'arc narratif
	Quote         string `json:"quote"`          // Citation du personnage
}

// BreathingTechnique représente un style de combat
type BreathingTechnique struct {
	ID          int    `json:"id"`          // ID unique
	Name        string `json:"name"`        // Nom du style (ex: "Water Breathing")
	Description string `json:"description"` // Description
}
