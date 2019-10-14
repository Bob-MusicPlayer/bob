package model

type Config struct {
	Hostname string    `json:"hostname"`
	Port     int       `json:"port"`
	LogLevel string    `json:"logLevel"`
	Player   []*Player `json:"player"`
}

type Player struct {
	Source   string `json:"source"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	IconUrl  string `json:"iconUrl"`
}
