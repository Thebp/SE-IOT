package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
)

type server struct {
	database Database
	mqtt     mqtt.Client
}

type Database interface {
	GetRoom(id string) (room Room, err error)
	GetRooms() ([]Room, error)
	GetBoard(id string) (board Board, err error)
	InsertRoom(room Room) (Room, error)
	UpdateRoom(room Room) error
	GetBoardsByRoom(roomId string) ([]Board, error)
	GetUnassignedBoards() ([]Board, error)
	InsertBoard(board Board) error
	UpdateBoard(board Board) error
	InsertLightData(data LightData) error
	GetLatestLightData(roomID string) ([]LightData, error)
}

func NewServer(database Database, mqtt mqtt.Client) *server {
	if token := mqtt.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return &server{database, mqtt}

}

func (s *server) run() {
	s.mqtt.Subscribe("board_discovery", 0, s.boardDiscovery)
	s.mqtt.Subscribe("lightdata", 0, s.processLightData)

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	})

	r.HandleFunc("/boards/{boardId}/ping", s.pingBoard).Methods("POST")
	r.HandleFunc("/rooms", s.getRooms).Methods("GET")
	r.HandleFunc("/rooms", s.postRoom).Methods("POST")
	r.HandleFunc("/rooms/{roomId}", s.putRoom).Methods("PUT")
	r.HandleFunc("/rooms/{roomId}/boards", s.getBoardsByRoom).Methods("GET")
	r.HandleFunc("/unassigned_boards", s.getUnassignedBoards).Methods("GET")
	r.HandleFunc("/boards/{boardId}", s.putBoard).Methods("PUT")

	fmt.Println("Listening...")
	http.ListenAndServe(":50002", r)
}

func (s *server) updateRoomLight(roomID string) {
	room, err := s.database.GetRoom(roomID)
	if err != nil {
		fmt.Println(err)
	}

	intensity := room.LedConfig.Intensity

	if room.LedConfig.DaylightHarvesting {
		lightdata, err := s.database.GetLatestLightData(roomID)
		if err != nil {
			fmt.Println(err)
			return
		}

		lightlevel := getAverageLightlevel(lightdata)
		intensity = getIntensityFromLightlevel(lightlevel)
	}

	boards, err := s.database.GetBoardsByRoom(roomID)
	if err != nil {
		fmt.Println(err)
		return
	}

	color := room.LedConfig.Color
	color.Red = int(float32(color.Red) / 100.0 * float32(intensity))
	color.Green = int(float32(color.Green) / 100.0 * float32(intensity))
	color.Blue = int(float32(color.Blue) / 100.0 * float32(intensity))

	color_json, err := json.Marshal(color)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, board := range boards {
		s.sendMessage(fmt.Sprintf("%v/led/rgb", board.ID), color_json)
	}

}

func (s *server) boardDiscovery(client mqtt.Client, msg mqtt.Message) {
	var board Board
	if err := json.Unmarshal(msg.Payload(), &board); err != nil {
		fmt.Println(err)
		return
	}
	b, err := s.database.GetBoard(board.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if b.ID == "" {
		err = s.database.InsertBoard(board)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("New board discovered: %v\n", board.ID)
	}
}

func (s *server) processLightData(client mqtt.Client, msg mqtt.Message) {
	var lightData LightData
	if err := json.Unmarshal(msg.Payload(), &lightData); err != nil {
		panic(err)
	}
	if lightData.BoardID == "" {
		fmt.Println("No BoardID on lightdata")
		return
	}
	lightData.Time = time.Now()
	err := s.database.InsertLightData(lightData)
	if err != nil {
		fmt.Println(err)
	}
	board, err := s.database.GetBoard(lightData.BoardID)
	if err != nil {
		fmt.Println(err)
	}
	go s.updateRoomLight(board.RoomID)
}

func (s *server) sendMessage(topic string, payload interface{}) {
	token := s.mqtt.Publish(topic, 0, false, payload)
	token.Wait()
}

func (s *server) pingBoard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	boardId := vars["boardId"]

	board, err := s.database.GetBoard(boardId)
	if err != nil {
		fmt.Println(err)
	}

	if board.ID != "" {
		s.sendMessage(fmt.Sprintf("%v/led/ping", board.ID), `{"ping":"ping"}`)
	}
	w.WriteHeader(200)
}

func (s *server) getRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rooms, err := s.database.GetRooms()
	if err != nil {
		fmt.Println(err)
	}

	rooms_json, err := json.Marshal(rooms)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Write(rooms_json)
}

func (s *server) getBoardsByRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	roomId := vars["roomId"]

	boards, err := s.database.GetBoardsByRoom(roomId)
	if err != nil {
		fmt.Println(err)
	}

	boards_json, err := json.Marshal(boards)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Write(boards_json)
}

func (s *server) getUnassignedBoards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	boards, err := s.database.GetUnassignedBoards()
	if err != nil {
		fmt.Println(err)
	}

	boards_json, err := json.Marshal(boards)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Write(boards_json)
}

func (s *server) putBoard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	boardId := vars["boardId"]

	var board Board
	err := json.NewDecoder(r.Body).Decode(&board)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Request body is invalid format", http.StatusBadRequest)
		return
	}

	if board.ID != boardId {
		http.Error(w, "Board ID in request body is different from ID in URL", http.StatusBadRequest)
		return
	}

	err = s.database.UpdateBoard(board)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	go s.updateRoomLight(board.RoomID)

	w.WriteHeader(200)
}

func (s *server) postRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var room Room
	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Request body is invalid format", http.StatusBadRequest)
		return
	}

	room, err = s.database.InsertRoom(room)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func (s *server) putRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	roomId := vars["roomId"]

	var room Room
	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Request body is invalid format", http.StatusBadRequest)
		return
	}

	if room.ID != roomId {
		http.Error(w, "Room ID in request body is different from ID in URL", http.StatusBadRequest)
		return
	}

	err = s.database.UpdateRoom(room)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	go s.updateRoomLight(room.ID)

	w.WriteHeader(200)
}
