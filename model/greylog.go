package model

import (
	"encoding/json"
)

type Greylog struct {
	Message      string `json:"message"`
	Level        int    `json:"level"`
	ShortMessage string `json:"short_message"`
	FullMessage  string `json:"full_message"`
	Host         string `json:"host"`
	CustomFields map[string]interface{}
}

func NewGreylog() Greylog {
	gl := Greylog{}
	cf := make(map[string]interface{})
	gl.CustomFields = cf
	return gl
}

func (g Greylog) String() string {

	entries := map[string]interface{}{}

	entries["message"] = g.Message
	entries["level"] = g.Level
	entries["short_message"] = g.ShortMessage
	entries["full_message"] = g.FullMessage
	entries["host"] = g.Host

	for k, v := range g.CustomFields {

		entries[k] = v
	}

	message, err := json.Marshal(entries)
	if err != nil {
		panic(err)
	}

	return string(message)
}
