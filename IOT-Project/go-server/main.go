package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type LightData struct {
	LightLevel string `json:"data"`
	Board      string `json:"board"`
}

type LightServer struct {
	Writer io.Writer
}

func main() {

	f, err := os.Create("light-" + strconv.FormatInt(time.Now().Unix(), 10) + ".csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	server := LightServer{f}
	r := mux.NewRouter()
	r.HandleFunc("/light", server.postLight)

	http.ListenAndServe(":50001", r)
}

func (s LightServer) postLight(w http.ResponseWriter, r *http.Request) {
	var message LightData
	_ = json.NewDecoder(r.Body).Decode(&message)
	//fmt.Println(message.Data)
	s.Writer.Write([]byte("Light: " + message.LightLevel + ", Board: " + message.Board + "\n\n"))

	light, err := strconv.Atoi(message.LightLevel)
	if err != nil {
		fmt.Println(err)
		return
	}
	if light > 100 {
		light = 100.0
	}
	intensity := int(255.0 - (float32(light) / 100.0 * 255.0))
	hex, err := strconv.ParseInt(fmt.Sprintf("%02x%02x%02x", intensity, intensity, intensity), 16, 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	message_json, _ := json.Marshal(hex)
	w.WriteHeader(200)
	w.Write([]byte(message_json))
}
