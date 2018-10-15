package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	timeSleep = 10
)

var (
	err error
	ch1 = make(chan []byte, 1)
)

type (
	InJsonStruct struct {
		Title  string `json:"title"`
		People int    `json:"people"`
		Area   int    `json:"area"`
	}

	OutJsonStruct struct {
		Country string         `json:"country"`
		City    []InJsonStruct `json:"city"`
	}

	NewOutStruct struct {
		CountrySity       string  `json:"country_city"`
		PopulationDensity float64 `json:"population_density"`
	}
)

func main() {
	defer close(ch1)
	readDir()
	country := getDeserializeJSON()
	setSerilazeJSON(country)
}

func setErr(err interface{}) {
	log.Fatal("Internal error: ", err)
}

func readDir() {
	var FileJSON = "InFile.json"
	var srcJSON []byte
	fmt.Println("Читаем текущую директорию")

	for i := timeSleep; i >= 0; i-- {
		srcJSON, err = ioutil.ReadFile(FileJSON)
		if err != nil {
			time.Sleep(time.Minute * 1)
			fmt.Printf("Положите фаил в директорию. Ждём %d\n", i)
		}
	}
	fmt.Println("Прочитали: ", string(srcJSON))
	ch1 <- srcJSON
}

func getDeserializeJSON() OutJsonStruct {
	srcJSON := <-ch1
	var newJSON OutJsonStruct
	if err := json.Unmarshal(srcJSON, &newJSON); err != nil {
		setErr(err)
	}
	fmt.Println("Десериализовали в: ", newJSON)
	return newJSON
}

func setSerilazeJSON(outJSON OutJsonStruct) {
	var newJSON []NewOutStruct
	for _, c := range outJSON.City {
		var nj NewOutStruct
		nj.CountrySity = string(outJSON.Country) + "." + string(c.Title)
		nj.PopulationDensity = float64(c.People / c.Area)
		newJSON = append(newJSON, nj)
	}
	fmt.Println("Данные для сериализации: ", newJSON)

	writeFileJSON, err := json.Marshal(newJSON)
	if err != nil {
		setErr(err)
	}

	f, err := os.Create("OutFile.json")
	_, err = f.Write([]byte(writeFileJSON))
	if err != nil {
		setErr(err)
	}
	fmt.Printf("Готово, фаил: %s\n", f.Name())
}
