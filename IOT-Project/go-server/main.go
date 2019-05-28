package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
)

type Board struct {
	ID     string `json:"id"`
	RoomID string `json:"room_id"`
}

type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type LightData struct {
	LightLevel int       `json:"lightlevel"`
	BoardID    string    `json:"board_id"`
	Time       time.Time `json:"time"`
}

type RGBData struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
}

type Database interface {
	//GetRoom(id string) (room Room, err error)
	GetRooms() ([]Room, error)
	GetBoard(id string) (board Board, err error)
	InsertRoom(room Room) (Room, error)
	GetBoardsByRoom(roomId string) ([]Board, error)
	GetUnassignedBoards() ([]Board, error)
	//GetBoards() ([]Board, error)
	InsertBoard(board Board) error
	UpdateBoard(board Board) error
	InsertLightData(data LightData) error
}

type server struct {
	database Database
	mqtt     mqtt.Client
}

func (s *server) boardDiscovery(client mqtt.Client, msg mqtt.Message) {
	var board Board
	if err := json.Unmarshal(msg.Payload(), &board); err != nil {
		panic(err)
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
}

func (s *server) sendMessage(topic string, payload interface{}) {
	// json_string, err := json.Marshal(payload)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	fmt.Printf("Sending %v to %v\n", payload, topic)
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

type inMemoryDatabase struct {
	boards      []Board
	rooms       []Room
	lightData   []LightData
	roomCounter int
}

func (d *inMemoryDatabase) InsertRoom(room Room) (Room, error) {
	room.ID = strconv.Itoa(d.roomCounter)
	d.rooms = append(d.rooms, room)
	d.roomCounter++
	return room, nil
}

func (d *inMemoryDatabase) GetRooms() ([]Room, error) {
	return d.rooms, nil
}

func (d *inMemoryDatabase) GetBoardsByRoom(roomId string) ([]Board, error) {
	board_list := make([]Board, 0)
	for _, board := range d.boards {
		if board.RoomID == roomId {
			board_list = append(board_list, board)
		}
	}
	return board_list, nil
}

func (d *inMemoryDatabase) GetUnassignedBoards() ([]Board, error) {
	board_list := make([]Board, 0)
	for _, board := range d.boards {
		if board.RoomID == "" {
			board_list = append(board_list, board)
		}
	}
	return board_list, nil
}

func (d *inMemoryDatabase) GetBoard(id string) (Board, error) {
	for _, board := range d.boards {
		if board.ID == id {
			return board, nil
		}
	}
	return Board{}, nil
}

func (d *inMemoryDatabase) UpdateBoard(board Board) error {
	for i, b := range d.boards {
		if b.ID == board.ID {
			d.boards[i] = board
			return nil
		}
	}
	return fmt.Errorf("No board with ID %v", board.ID)
}

func (d *inMemoryDatabase) InsertBoard(board Board) error {
	fmt.Printf("Inserting board %v\n", board)
	d.boards = append(d.boards, board)
	return nil
}

func (d *inMemoryDatabase) InsertLightData(lightData LightData) error {
	fmt.Printf("Inserting lightData: %v\n", lightData)
	d.lightData = append(d.lightData, lightData)
	return nil
}

func main() {

	opts := mqtt.NewClientOptions()
	opts.AddBroker("mndkk.dk:1883")
	opts.SetClientID("server")
	opts.SetUsername("iot")
	opts.SetPassword("uS831ACCL6sZHz4")

	mqtt_client := mqtt.NewClient(opts)
	if token := mqtt_client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	database := inMemoryDatabase{make([]Board, 0), make([]Room, 0), make([]LightData, 0), 0}

	server := server{&database, mqtt_client}
	server.mqtt.Subscribe("board_discovery", 0, server.boardDiscovery)
	server.mqtt.Subscribe("lightdata", 0, server.processLightData)

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	})

	r.HandleFunc("/boards/{boardId}/ping", server.pingBoard).Methods("POST")
	r.HandleFunc("/rooms", server.getRooms).Methods("GET")
	r.HandleFunc("/rooms", server.postRoom).Methods("POST")
	r.HandleFunc("/rooms/{roomId}/boards", server.getBoardsByRoom).Methods("GET")
	r.HandleFunc("/unassigned_boards", server.getUnassignedBoards).Methods("GET")
	r.HandleFunc("/boards/{boardId}", server.putBoard).Methods("PUT")

	fmt.Println("Listening...")
	http.ListenAndServe(":50002", r)
}
