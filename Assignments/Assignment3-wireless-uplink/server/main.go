package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type LightData struct {
	LightLevel  string `json:"light"`
	Temperature string `json:"temp"`
	Board       string `json:"board"`
}

type DataLogger struct {
	Filename string
}

func (l *DataLogger) Init() error {
	l.Filename = "data-" + strconv.FormatInt(time.Now().Unix(), 10) + ".csv"
	f, err := os.Create(l.Filename)
	if err != nil {
		return err
	}
	f.Write([]byte("Lightlevel, Temperature, Board, Time\n"))
	f.Close()
	return nil
}

func (l *DataLogger) Log(message string) error {
	f, err := os.OpenFile(l.Filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	f.Write([]byte(message))
	f.Close()
	return nil
}

func main() {
	logger := DataLogger{}
	err := logger.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/data", logger.postLight)

	http.ListenAndServe(":50001", r)
}

func (l *DataLogger) postLight(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("IGmJuizNPufozspY4Xry")
	if auth != "NuB19BA6TV8bmvTpgOgo" {
		fmt.Println("401: Unauthorized")
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
		return
	}
	var message LightData
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		w.Write([]byte("Message malformed"))
		return
	}
	err = l.Log(message.LightLevel + ", " + message.Temperature + ", " + message.Board + ", " + strconv.FormatInt(time.Now().Unix(), 10) + "\n")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		w.Write([]byte("Internal error"))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("Received"))
}
