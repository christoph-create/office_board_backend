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

var dataBase *gorm.DB

type ScoreEntry struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	OnePoint  uint
	HalfPoint uint
	Total     uint
}

type GitAccident struct {
	ID   uint `gorm:"primaryKey"`
	Time time.Time
}

type TimeSince struct {
	Seconds uint
}

func Init() error {
	err := initDB()
	if err != nil {
		return err
	}

	http.HandleFunc("/score/update", updateScore)
	http.HandleFunc("/score/load", loadScores)

	http.HandleFunc("/accident/load", loadLastAccident)
	http.HandleFunc("/accident/reset", resetAccident)

	return nil
}

func initDB() error {
	var err error
	dataBase, err = gorm.Open(sqlite.Open("/mnt/office_board.db"), &gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Silent)})
	if err != nil {
		return err
	}

	err = dataBase.AutoMigrate(&ScoreEntry{})
	if err != nil {
		return err
	}

	err = dataBase.AutoMigrate(&GitAccident{})
	if err != nil {
		return err
	}

	return err
}

func updateScore(w http.ResponseWriter, r *http.Request) {
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

	dataBase.Save(data)
}

func loadScores(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	var data []ScoreEntry

	dataBase.Find(&data)
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

func loadLastAccident(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	var accident GitAccident

	dataBase.Last(&accident)
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

func resetAccident(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	dataBase.Save(&GitAccident{Time: time.Now()})
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
