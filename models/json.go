package models

type Info struct {
	CurrentSha string `json:"current_sha"`
}

type Dcdr struct {
	Info     Info  `json:"info"`
	Features Flags `json:"features"`
}
type Flags map[string]interface{}

type DcdrMap struct {
	Dcdr Dcdr `json:"dcdr"`
}
