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

func determineListenAddress() (string, error) {
	//port := os.Getenv("PORT")
	port := "5000"
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	log.Printf("NasaPortal API Live at: 5000 or PORT.env")
	return ":" + port, nil
}

func main() {

	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return origin == "http://www.ourcosmos.us" || origin == "http://localhost:8080"

		},
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "X-Custom-Header", "Origin"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})
	router, err := rest.MakeRouter(
		rest.Get("/api/apod", GetAstronomyToday),
		rest.Get("/isMars", GetMarsRoverData),
		rest.Get("/isMars/:camera/:id", GetMarsRoverDataID),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(addr, api.MakeHandler()))
}

var store = map[string]*AstronomyPicOfDay{}

var lock = sync.RWMutex{}

// HTTP Req To Nasa for Mars Rover Data
func GetMarsRoverDataID(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	camera := r.PathParam("camera")

	url := " "
	if camera == " " {
		url = "https://api.nasa.gov/mars-photos/api/v1/rovers/curiosity/photos?sol=" + id + "&api_key=iz6rQYs0Ws9LWTf2SlBgSPpyHKerfx6JUBVYCnoC"
	} else {
		url = "https://api.nasa.gov/mars-photos/api/v1/rovers/curiosity/photos?sol=" + id + "&camera=" + camera + "&api_key=iz6rQYs0Ws9LWTf2SlBgSPpyHKerfx6JUBVYCnoC"
	}

	lock.RLock()
	var isMars *MarsRovers
	isMars = getMarsRoverFromNasa(url)
	lock.RUnlock()

	if isMars == nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(isMars)

}

// HTTP Req To Nasa for Mars Rover Data
func GetMarsRoverData(w rest.ResponseWriter, r *rest.Request) {
	url := "https://api.nasa.gov/mars-photos/api/v1/rovers/curiosity/photos?sol=1&camera=fhaz&api_key=iz6rQYs0Ws9LWTf2SlBgSPpyHKerfx6JUBVYCnoC"

	lock.RLock()
	var isMars *MarsRovers
	isMars = getMarsRoverFromNasa(url)
	lock.RUnlock()

	if isMars == nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(isMars)

}

// HTTP Req To Nasa for Astronomy of the Day
func GetAstronomyToday(w rest.ResponseWriter, r *rest.Request) {
	code := "apod"
	url := "https://api.nasa.gov/planetary/" + code + "?api_key=iz6rQYs0Ws9LWTf2SlBgSPpyHKerfx6JUBVYCnoC"

	lock.RLock()
	var today *AstronomyPicOfDay
	today = getNasaData(url)
	lock.RUnlock()

	if today == nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(today)
}

func getNasaData(source string) (record *AstronomyPicOfDay) {
	url := fmt.Sprintf(source)

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

func getMarsRoverFromNasa(source string) (recordMars *MarsRovers) {
	url := fmt.Sprintf(source)

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
	if eff := json.NewDecoder(resp.Body).Decode(&recordMars); eff != nil {
		log.Println(err)
	}
	fmt.Println("Mars Rover Content: ", recordMars)
	return
}
