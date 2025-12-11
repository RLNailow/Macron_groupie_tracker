package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Character struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

func charactersHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://www.demonslayer-api.com/api/characters")
	if err != nil {
		http.Error(w, "Impossible de récupérer les personnages", 500)
		return
	}
	defer resp.Body.Close()

	var characters []Character
	json.NewDecoder(resp.Body).Decode(&characters)

	tmpl := template.Must(template.ParseFiles("templates/characters.html"))
	tmpl.Execute(w, characters)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/characters", charactersHandler)

	log.Println("Serveur lancé sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}
