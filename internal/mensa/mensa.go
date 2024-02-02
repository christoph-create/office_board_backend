package mensa

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// datum;tag;warengruppe;name;kennz;preis;stud;bed;gast
type food struct {
	Date    string `json:"date"`
	Day     string `json:"day"`
	Group   string `json:"group"`
	Id      string `json:"id"`
	Name    string `json:"name"`
	Student string `json:"student"`
}

func Init() {
	fmt.Println("start listening mensa")
	http.HandleFunc("/mensa/today/main", onContentRequestMain)
	http.HandleFunc("/mensa/today/side", onContentRequestSide)
}

func onContentRequestMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enableCors(&w)
	_, week := time.Now().UTC().ISOWeek()
	resp, err := http.Get(fmt.Sprintf("http://www.stwno.de/infomax/daten-extern/csv/HS-R-tag/%v.csv", week))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	fullDataSet := []food{}

	dec := transform.NewReader(resp.Body, charmap.ISO8859_1.NewDecoder())

	scanner := bufio.NewScanner(dec)
	for scanner.Scan() {
		tmp := strings.Split(scanner.Text(), ";")
		fullDataSet = append(fullDataSet, food{
			Date:    tmp[0],
			Group:   tmp[2],
			Name:    strings.Split(tmp[3], "(")[0],
			Id:      tmp[4],
			Student: tmp[6],
		})
	}

	fullDataSet = fullDataSet[1:]

	// filter data
	dataTodayMain := []food{}

	dateToday := time.Now().Format("02.01.2006")

	// TODO: remove
	// dateToday = "30.01.2024"

	for _, f := range fullDataSet {
		if f.Date == dateToday && strings.Contains(f.Group, "H") {
			dataTodayMain = append(dataTodayMain, f)
		}
	}

	jsonFoods, err := json.Marshal(dataTodayMain)
	if err != nil {
		fmt.Println(err)
	}

	_, err = w.Write(jsonFoods)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("data send")
}

func onContentRequestSide(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enableCors(&w)
	_, week := time.Now().UTC().ISOWeek()
	resp, err := http.Get(fmt.Sprintf("http://www.stwno.de/infomax/daten-extern/csv/HS-R-tag/%v.csv", week))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	fullDataSet := []food{}

	dec := transform.NewReader(resp.Body, charmap.ISO8859_1.NewDecoder())

	scanner := bufio.NewScanner(dec)
	for scanner.Scan() {
		tmp := strings.Split(scanner.Text(), ";")
		fullDataSet = append(fullDataSet, food{
			Date:    tmp[0],
			Group:   tmp[2],
			Name:    strings.Split(tmp[3], "(")[0],
			Id:      tmp[4],
			Student: tmp[6],
		})
	}

	fullDataSet = fullDataSet[1:]

	// filter data
	dataTodaySide := []food{}

	dateToday := time.Now().Format("02.01.2006")

	// TODO: remove
	// dateToday = "26.01.2024"

	for _, f := range fullDataSet {
		if f.Date == dateToday && strings.Contains(f.Group, "B") {
			dataTodaySide = append(dataTodaySide, f)
		}
	}

	jsonFoods, err := json.Marshal(dataTodaySide)
	if err != nil {
		fmt.Println(err)
	}

	_, err = w.Write(jsonFoods)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("data send")
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
