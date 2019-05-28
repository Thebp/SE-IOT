package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	//GetRooms() ([]Room, error)
	GetBoard(id string) (board Board, err error)
	//GetBoardsByRoom(roomId string) ([]Board, error)
	//GetBoards() ([]Board, error)
	InsertBoard(board Board) error
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

	s.sendMessage(fmt.Sprintf("%v/led/ping", board.ID), `{"ping":"ping"}`)
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

type inMemoryDatabase struct {
	boards    []Board
	rooms     []Room
	lightData []LightData
}

func (d *inMemoryDatabase) GetBoard(id string) (Board, error) {
	for _, board := range d.boards {
		if board.ID == id {
			return board, nil
		}
	}
	return Board{}, nil
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

	// _, err := os.Create("light-" + strconv.FormatInt(time.Now().Unix(), 10) + ".csv")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	opts := mqtt.NewClientOptions()
	opts.AddBroker("mndkk.dk:1883")
	opts.SetClientID("server")
	opts.SetUsername("iot")
	opts.SetPassword("uS831ACCL6sZHz4")

	mqtt_client := mqtt.NewClient(opts)
	if token := mqtt_client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	database := inMemoryDatabase{make([]Board, 0), make([]Room, 0), make([]LightData, 0)}

	server := server{&database, mqtt_client}
	server.mqtt.Subscribe("board_discovery", 0, server.boardDiscovery)
	server.mqtt.Subscribe("lightdata", 0, server.processLightData)

	r := mux.NewRouter()
	//r.HandleFunc("/light", server.postLight)

	fmt.Println("Listening...")
	http.ListenAndServe(":50001", r)
}
