package main

import "time"

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
