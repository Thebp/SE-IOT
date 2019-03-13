package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type LightData struct {
	Data string `json:"data"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/light", postLight)

	http.ListenAndServe(":5467", r)
}

func postLight(w http.ResponseWriter, r *http.Request) {
	var message LightData
	_ = json.NewDecoder(r.Body).Decode(&message)
	fmt.Println(message.Data)

	message_json, _ := json.Marshal(message)
	w.WriteHeader(200)
	w.Write([]byte(message_json))
}
