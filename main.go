package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

func main() {

	bands := GetArtist()
	for _, band := range bands {
		http.HandleFunc("/artist/"+band.Name, artistHandler(band.Id))
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", firstPageHandler)
	http.HandleFunc("/artist", mainPageHandler)

	fmt.Println("starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

func firstPageHandler(w http.ResponseWriter, r *http.Request) {
	if !checkURL(r.URL.Path) {
		// fmt.Println("error 404, page not found")
		fmt.Fprintf(w, "404 Not Found")
		return
	}

	tmpl, err := template.ParseFiles("templates/first.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	if !checkURL(r.URL.Path) {
		// fmt.Println("error 404, page not found")
		fmt.Fprintf(w, "404 Not FOund")
		return
	}

	artists := GetArtist()
	tmpl, err := template.ParseFiles("templates/artisthandler.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	// fmt.Println("Templated is parsed lol")

	err = tmpl.Execute(w, artists)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
	// fmt.Println("checkcheck")
}

func artistHandler(id int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Found band")
		Band := Bandinfo{
			Relations: GetRelation(id),
			Artist:    GetArtist()[id-1],
		}
		tmpl, _ := template.ParseFiles("templates/artistinfo.html")
		tmpl.Execute(w, Band)
	}
}

var client *http.Client

func GetJson(url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

// struct for the bands
type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationdate"`
	FirstAlbum   string   `json:"firstalbum"`
}

type Relations struct {
	DatesLocations interface{} `json:"datesLocations"`
}

type Bandinfo struct {
	Relations Relations
	Artist    Artist
}

func GetArtist() []Artist {
	client = &http.Client{Timeout: 10 * time.Second}
	url := "https://groupietrackers.herokuapp.com/api/artists"
	var artists []Artist
	err := GetJson(url, &artists)
	if err != nil {
		fmt.Printf("error getting Artist: %s\n", err.Error())
	}
	return artists
}

func GetRelation(id int) Relations {
	str := strconv.Itoa(id)
	url := "https://groupietrackers.herokuapp.com/api/relation/" + str
	var relations Relations
	err := GetJson(url, &relations)
	if err != nil {
		fmt.Printf("error getting Artist: %s\n", err.Error())
	}
	return relations
}

func checkURL(url string) bool {
	artInfo := GetArtist()
	for _, i := range artInfo {
		if url == "/"+i.Name || url == "/" || url == "/artist" {
			return true
		}
	}
	return false
}
