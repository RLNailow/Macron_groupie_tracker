package models

type Character struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Race  string // ajouté localement
}

func GetRace(name string) string {
	demons := map[string]bool{
		"Akaza":     true,
		"Doma":      true,
		"Kokushibo": true,
		"Gyutaro":   true,
		"Daki":      true,
		"Kaigaku":   true,
		"Rui":       true,
		"Enmu":      true,
		"Gyokko":    true,
		"Nakime":    true,
		"Hantengu":  true,
		"Kyogai":    true,
		"Muzan":     true,
		"Tamayo":    true,
		"Yushiro":   true,
	}

	if demons[name] {
		return "demon"
	}
	return "human"
}
