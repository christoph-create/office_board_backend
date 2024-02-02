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

func Init() error {
	err := initDB()
	if err != nil {
		return err
	}
	// TODO remove when create entry is implemeted
	fill()

	http.HandleFunc("/score/update", updateScore)
	http.HandleFunc("/score/load", loadScores)

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

func fill() {
	dataBase.Save(&ScoreEntry{"1", "Peter", 2, 1, 3})
	dataBase.Save(&ScoreEntry{"2", "Golo", 2, 1, 11})
	dataBase.Save(&ScoreEntry{"3", "Michi G", 2, 1, 9})
	dataBase.Save(&ScoreEntry{"4", "Silvio", 0, 0, 5})
	dataBase.Save(&ScoreEntry{"5", "Christoph", 1, 2, 11})
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
