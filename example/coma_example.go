package main

import (
	comasdkgo "coma-sdk-go"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Application struct {
		Graceful struct {
			Duration      string `json:"DURATION"`
			SleepDuration string `json:"SLEEP_DURATION"`
		} `json:"GRACEFUL"`
		PrintHello struct {
			Enable bool `json:"ENABLE"`
		} `json:"PRINT_HELLO"`
		Name string `json:"NAME"`
		Port int    `json:"PORT"`
	} `json:"APPLICATION"`
}

var (
	cfg *Config
)

func handler(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	if cfg.Application.PrintHello.Enable {
		// w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("hello world"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func main() {
	coma, err := comasdkgo.New(
		"http://localhost:3001/swagger/index.html#/Config/get_v1_configuration",
		"localhost",
		"3001",
		"EoCKgsUO2rMZdz1pqlJ0rvXTSCLDhjuomEyY",
		comasdkgo.SetRetry(10),
		comasdkgo.SetRetryWaitTime(5*time.Second))
	if err != nil {
		log.Fatal("err: connect ", err)
		return
	}

	coma.Observe(&cfg)

	http.HandleFunc("/", handler)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
