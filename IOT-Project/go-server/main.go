package main

import (
	"fmt"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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
	database := inMemoryDatabase{make(map[string]Board), make(map[string]Room), make([]LightData, 0), 0}

	server := NewServer(&database, mqtt_client)
	server.run()
}
