package mensa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Food struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Prices   Prices   `json:"prices"`
	Notes    []string `json:"notes"`
}

type Prices struct {
	Students float32 `json:"students"`
}

const (
	mainDish string = "Hauptgerichte"
	sideDish string = "Beilagen"
	soup     string = "Suppe"
	dessert  string = "Nachspeisen"
)

func Init() {
	fmt.Println("start listening mensa")
	http.HandleFunc("/mensa/today/main", onContentRequestMain)
	http.HandleFunc("/mensa/today/side", onContentRequestSide)
}

func onContentRequestMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enableCors(&w)
	data, err := getFood(mainDish)
	if err != nil {
		fmt.Println(err)
	}
	_, err = w.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("data send")
}

func onContentRequestSide(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enableCors(&w)
	data, err := getFood(sideDish)
	if err != nil {
		fmt.Println(err)
	}
	_, err = w.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("data send")
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getFood(category string) ([]byte, error) {
	day := time.Now().Format("2006-01-02")
	resp, err := http.Get(fmt.Sprintf("https://openmensa.org/api/v2/canteens/195/days/%s/meals", day))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	foodData := &[]Food{}
	err = json.NewDecoder(resp.Body).Decode(foodData)
	if err != nil {
		return nil, err
	}
	dishes := []Food{}
	for _, meal := range *foodData {
		if meal.Category == category {
			dishes = append(dishes, meal)
		}
	}

	data, err := json.Marshal(dishes)
	if err != nil {
		return nil, err
	}
	return data, nil
}
