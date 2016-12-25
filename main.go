package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/ant0ine/go-json-rest/rest"
)

// Nasa "AstronomyPicOfDay"" JSON struct
type AstronomyPicOfDay struct {
	Copyright      string `json:"copyright"`
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	Hdurl          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	URL            string `json:"url"`
}

// Nasa "Mars Rovers" JSON struct
type MarsRovers struct {
	Photos []struct {
		ID     int `json:"id"`
		Sol    int `json:"sol"`
		Camera struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			RoverID  int    `json:"rover_id"`
			FullName string `json:"full_name"`
		} `json:"camera"`
		ImgSrc    string `json:"img_src"`
		EarthDate string `json:"earth_date"`
		Rover     struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			LandingDate string `json:"landing_date"`
			LaunchDate  string `json:"launch_date"`
			Status      string `json:"status"`
			MaxSol      int    `json:"max_sol"`
			MaxDate     string `json:"max_date"`
			TotalPhotos int    `json:"total_photos"`
			Cameras     []struct {
				Name     string `json:"name"`
				FullName string `json:"full_name"`
			} `json:"cameras"`
		} `json:"rover"`
	} `json:"photos"`
}

func main() {
	//nasaSource := "apod"

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/apod", GetAstronomyToday),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))

}

var store = map[string]*AstronomyPicOfDay{}

var lock = sync.RWMutex{}

func GetAstronomyToday(w rest.ResponseWriter, r *rest.Request) {
	code := "apod"
	lock.RLock()
	var today *AstronomyPicOfDay
	today = getNasaData(code)
	lock.RUnlock()

	if today == nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(today)
}

func getNasaData(source string) (record *AstronomyPicOfDay) {
	url := fmt.Sprintf("https://api.nasa.gov/planetary/" + source + "?api_key=iz6rQYs0Ws9LWTf2SlBgSPpyHKerfx6JUBVYCnoC")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// HTTP Client
	client := &http.Client{}

	// Sends HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	// Close the body at end
	defer resp.Body.Close()

	// Decode Json for reading
	if eff := json.NewDecoder(resp.Body).Decode(&record); eff != nil {
		log.Println(err)
	}

	fmt.Println("NASA Content of the Day: ", record)
	return
}
