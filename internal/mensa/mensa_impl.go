package mensa

import (
	"encoding/json"
	"fmt"
	datamanager "main/internal/dataManager"
	"net/http"
)

var singleton *mensa

type mensa struct {
	dataManager datamanager.DataManager
}

type Food struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Prices   Prices   `json:"prices"`
	Notes    []string `json:"notes"`
}

type RelevantFood struct {
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Prices   Prices   `json:"prices"`
	Rating   int      `json:"rating"`
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

func Init(db datamanager.DataManager) (Mensa, error) {
	if singleton != nil {
		return singleton, nil
	}
	singleton = &mensa{dataManager: db}
	fmt.Println("start listening mensa")
	http.HandleFunc("/mensa/main", singleton.onContentRequestMain)
	http.HandleFunc("/mensa/side", singleton.onContentRequestSide)
	return singleton, nil
}

func (mensa *mensa) onContentRequestMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, "Date parameter is missing", http.StatusBadRequest)
		return
	}

	enableCors(&w)
	data, err := mensa.getFood(mainDish, date)
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

func (mensa *mensa) onContentRequestSide(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, "Date parameter is missing", http.StatusBadRequest)
		return
	}

	enableCors(&w)
	data, err := mensa.getFood(sideDish, date)
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

func (mensa *mensa) getFood(category string, date string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("https://openmensa.org/api/v2/canteens/195/days/%s/meals", date))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	foodData := &[]Food{}
	err = json.NewDecoder(resp.Body).Decode(foodData)
	if err != nil {
		return nil, err
	}
	dishes := []RelevantFood{}
	for _, meal := range *foodData {
		if meal.Category == category {
			dishes = append(dishes, RelevantFood{Name: meal.Name, Category: meal.Category, Prices: meal.Prices, Notes: meal.Notes, Rating: mensa.dataManager.GetAverageRating(meal.Name)})
		}
	}

	data, err := json.Marshal(dishes)
	if err != nil {
		return nil, err
	}
	return data, nil
}
