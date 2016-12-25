package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func getNasaData(source string) {
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

	// Load data from JSON
	var record AstronomyPicOfDay

	// Decode Json for reading
	if eff := json.NewDecoder(resp.Body).Decode(&record); eff != nil {
		log.Println(err)
	}

	fmt.Println("NASA Content of the Day: ", record.Explanation)
}

func main() {
	nasaSource := "apod"
	getNasaData(nasaSource)
}
