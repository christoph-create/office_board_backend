package datamanager

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var singleton *dataManager

type dataManager struct {
	db *gorm.DB
}

type ScoreEntry struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	OnePoint  uint
	HalfPoint uint
	Total     uint
}

type RatingEntry struct {
	ID      string `gorm:"primaryKey"`
	Name    string
	Rating  int
	RatedBy string
}

type GitAccident struct {
	ID   uint `gorm:"primaryKey"`
	Time time.Time
}

type TimeSince struct {
	Seconds uint
}

func New() (DataManager, error) {
	if singleton != nil {
		return singleton, nil
	}
	newDB, err := singleton.initDB()
	if err != nil {
		return nil, err
	}
	singleton = &dataManager{db: newDB}
	http.HandleFunc("/score/update", singleton.updateScore)
	http.HandleFunc("/score/load", singleton.loadScores)

	http.HandleFunc("/accident/load", singleton.loadLastAccident)
	http.HandleFunc("/accident/reset", singleton.resetAccident)

	return singleton, nil
}

func (dataManager *dataManager) initDB() (*gorm.DB, error) {
	var err error
	dataManager.db, err = gorm.Open(sqlite.Open("/mnt/office_board.db"), &gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Silent)})
	if err != nil {
		return nil, err
	}

	err = dataManager.db.AutoMigrate(&ScoreEntry{})
	if err != nil {
		return nil, err
	}

	err = dataManager.db.AutoMigrate(&GitAccident{})
	if err != nil {
		return nil, err
	}

	return dataManager.db, err
}

func (dataManager *dataManager) updateScore(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	request, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("failed to read request body: ", err)
		return
	}

	var data ScoreEntry
	err = json.Unmarshal(request, &data)
	if err != nil {
		fmt.Println("failed to read request body: ", err)
		return
	}

	dataManager.db.Save(data)
}

func (dataManager *dataManager) updateRating(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	request, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("failed to read request body: ", err)
		return
	}

	var data RatingEntry
	err = json.Unmarshal(request, &data)
	if err != nil {
		fmt.Println("failed to read request body: ", err)
		return
	}

	dataManager.db.Save(data)
}

func (dataManager *dataManager) loadScores(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	var data []ScoreEntry

	dataManager.db.Find(&data)
	dataJson, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	_, err = w.Write(dataJson)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (dataManager *dataManager) loadLastAccident(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	var accident GitAccident

	dataManager.db.Last(&accident)
	seconds := uint(time.Since(accident.Time).Seconds())

	secondsJson, err := json.Marshal(&TimeSince{Seconds: seconds})
	if err != nil {
		fmt.Println(err)
	}

	_, err = w.Write(secondsJson)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (dataManager *dataManager) resetAccident(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	dataManager.db.Save(&GitAccident{Time: time.Now()})
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (dataManager *dataManager) GetAverageRating(name string) int {
	return 5
}
