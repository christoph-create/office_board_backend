package rvv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Rvv interface{}

type rvv struct {
	client *http.Client
}

var singleton *rvv

func New() (Rvv, error) {
	if singleton != nil {
		return singleton, nil
	}

	singleton = &rvv{
		client: &http.Client{},
	}

	go func() {
		fmt.Println("start listening rvv")
		http.HandleFunc("/rvv/techcampus", singleton.fetchSchedule)
		err := http.ListenAndServe(":8124", nil)
		if err != nil {
			fmt.Println(err)
		}
	}()

	return singleton, nil
}

type Departure struct {
	Platform     string    `json:"platform"`
	StopName     string    `json:"stopName"`
	DateTime     *StopTime `json:"dateTime"`
	RealDateTime *StopTime `json:"realDateTime"`
	ServingLine  Line      `json:"servingLine"`
}

type StopTime struct {
	Year    string `json:"year"`
	Month   string `json:"month"`
	Day     string `json:"day"`
	Weekday string `json:"weekday"`
	Hour    string `json:"hour"`
	Minute  string `json:"minute"`
}

type Line struct {
	Direction string `json:"direction"`
	Number    string `json:"number"`
}

type DepartureJson struct {
	Direction string `json:"direction"`
	Number    string `json:"number"`
	Time      string `json:"time"`
	TimePlus  string `json:"timePlus"`
}

func (rv *rvv) fetchSchedule(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start rvv")

	w.Header().Set("Content-Type", "application/json")
	enableCors(&w)

	url := "https://mobile.defas-fgi.de/beg/json/XML_DM_REQUEST"
	params := map[string]string{
		"outputFormat":           "JSON",
		"language":               "de",
		"stateless":              "1",
		"type_dm":                "stop",
		"mode":                   "direct",
		"useRealtime":            "1",
		"ptOptionActive":         "1",
		"mergeDep":               "1",
		"deleteAssignedStops_dm": "1",
		"limit":                  "20",
		"name_dm":                "TechCampus/OTH",
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("failed to create HTTP request: %v", err.Error())
		return
	}

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	res, err := rv.client.Do(req)
	if err != nil {
		fmt.Println("failed to perform HTTP request: %w", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Printf("HTTP request failed with status code: %d", res.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("failed to read response body: %w", err)
		return
	}

	// fmt.Println(string(body))

	var response struct {
		DepartureList []Departure `json:"departureList"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("failed to parse JSON response: %w", err)
		return
	}

	responseFilterd := []DepartureJson{}

	for _, departure := range response.DepartureList {
		// to city
		if departure.Platform == "2" {
			if len(departure.DateTime.Minute) < 2 {
				departure.DateTime.Minute = fmt.Sprintf("0%s", departure.DateTime.Minute)
			}
			timeString := fmt.Sprintf("%s:%s", departure.DateTime.Hour, departure.DateTime.Minute)
			timeVal, err := time.Parse("15:04", timeString)
			if err != nil {
				fmt.Println("failed to parse time string: %w", err)
				return
			}

			var diff string

			if departure.RealDateTime != nil {
				if len(departure.RealDateTime.Minute) < 2 {
					departure.RealDateTime.Minute = fmt.Sprintf("0%s", departure.RealDateTime.Minute)
				}
				timePlusString := fmt.Sprintf("%s:%s", departure.RealDateTime.Hour, departure.RealDateTime.Minute)
				timePlusVal, err := time.Parse("15:04", timePlusString)
				if err != nil {
					fmt.Println("failed to parse time string: %w", err)
					return
				}

				diff = fmt.Sprint(timePlusVal.Sub(timeVal).Minutes())
			}

			newDeparture := DepartureJson{
				Direction: departure.ServingLine.Direction,
				Number:    departure.ServingLine.Number,
				Time:      timeString,
				TimePlus:  diff,
			}
			responseFilterd = append(responseFilterd, newDeparture)
		}
	}

	jsonDeparture, err := json.Marshal(responseFilterd)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(jsonDeparture)
	fmt.Println("data send rvv")
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
