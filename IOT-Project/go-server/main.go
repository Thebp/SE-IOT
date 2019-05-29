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
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	LedConfig LedConfig `json:"led_config"`
	Boards    []Board   `json:"boards"`
}

type LedConfig struct {
	Intensity          int     `json:"intensity"`
	Color              RGBData `json:"color"`
	DaylightHarvesting bool    `json:"daylight_harvesting"`
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
	GetRoom(id string) (room Room, err error)
	GetRooms() ([]Room, error)
	GetBoard(id string) (board Board, err error)
	InsertRoom(room Room) (Room, error)
	UpdateRoom(room Room) error
	GetBoardsByRoom(roomId string) ([]Board, error)
	GetUnassignedBoards() ([]Board, error)
	//GetBoards() ([]Board, error)
	InsertBoard(board Board) error
	UpdateBoard(board Board) error
	InsertLightData(data LightData) error
	GetLatestLightData(roomID string) ([]LightData, error)
}

type server struct {
	database Database
	mqtt     mqtt.Client
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

func getAverageLightlevel(lightdata []LightData) int {
	if len(lightdata) == 0 {
		return 0
	}
	total := 0
	for _, lightdatum := range lightdata {
		total += lightdatum.LightLevel
	}
	return total / len(lightdata)
}

func getIntensityFromLightlevel(lightlevel int) int {
	if lightlevel > 100 {
		lightlevel = 100
	}
	return 100 - lightlevel
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

type inMemoryDatabase struct {
	boards      map[string]Board
	rooms       map[string]Room
	lightData   []LightData
	roomCounter int
}

func (d *inMemoryDatabase) InsertRoom(room Room) (Room, error) {
	room.ID = strconv.Itoa(d.roomCounter)
	d.rooms[room.ID] = room
	d.roomCounter++
	return room, nil
}

func (d *inMemoryDatabase) UpdateRoom(room Room) error {
	if _, ok := d.rooms[room.ID]; ok {
		d.rooms[room.ID] = room
		return nil
	}
	return fmt.Errorf("No room with ID %v", room.ID)
}

func (d *inMemoryDatabase) GetRoom(id string) (Room, error) {
	if room, ok := d.rooms[id]; ok {
		return room, nil
	}
	return Room{}, nil
}

func (d *inMemoryDatabase) GetRooms() ([]Room, error) {
	rooms := make([]Room, 0, len(d.rooms))
	for _, room := range d.rooms {
		boards, _ := d.GetBoardsByRoom(room.ID)
		room.Boards = boards
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (d *inMemoryDatabase) GetBoardsByRoom(roomId string) ([]Board, error) {
	boards := make([]Board, 0)
	for _, board := range d.boards {
		if board.RoomID == roomId {
			boards = append(boards, board)
		}
	}
	return boards, nil
}

func (d *inMemoryDatabase) GetUnassignedBoards() ([]Board, error) {
	boards := make([]Board, 0)
	for _, board := range d.boards {
		if board.RoomID == "" {
			boards = append(boards, board)
		}
	}
	return boards, nil
}

func (d *inMemoryDatabase) GetBoard(id string) (Board, error) {
	if board, ok := d.boards[id]; ok {
		return board, nil
	}
	return Board{}, nil
}

func (d *inMemoryDatabase) UpdateBoard(board Board) error {
	if _, ok := d.boards[board.ID]; ok {
		d.boards[board.ID] = board
		return nil
	}
	return fmt.Errorf("No board with ID %v", board.ID)
}

func (d *inMemoryDatabase) InsertBoard(board Board) error {
	d.boards[board.ID] = board
	return nil
}

func (d *inMemoryDatabase) InsertLightData(lightData LightData) error {
	d.lightData = append(d.lightData, lightData)
	return nil
}

func (d *inMemoryDatabase) GetLatestLightData(roomID string) ([]LightData, error) {
	lightdata := make([]LightData, 0)
	boards, _ := d.GetBoardsByRoom(roomID)
	for _, board := range boards {
		lightdata = append(lightdata, d.getLatestLightData(board.ID))
	}
	return lightdata, nil
}

func (d *inMemoryDatabase) getLatestLightData(boardID string) LightData {
	latest := LightData{}
	for _, lightdata := range d.lightData {
		if lightdata.Time.Unix() > latest.Time.Unix() {
			latest = lightdata
		}
	}
	return latest
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

	database := inMemoryDatabase{make(map[string]Board), make(map[string]Room), make([]LightData, 0), 0}

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
	r.HandleFunc("/rooms/{roomId}", server.putRoom).Methods("PUT")
	r.HandleFunc("/rooms/{roomId}/boards", server.getBoardsByRoom).Methods("GET")
	r.HandleFunc("/unassigned_boards", server.getUnassignedBoards).Methods("GET")
	r.HandleFunc("/boards/{boardId}", server.putBoard).Methods("PUT")

	fmt.Println("Listening...")
	http.ListenAndServe(":50002", r)
}
