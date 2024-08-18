package example

import (
	"encoding/json"
	"math/rand"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func GetRandomId() int {
	min := 10000
	max := 99999
	return rand.Intn(max-min) + min
}

func TodoToJSON(t Todo) (string, error) {
	json, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	return string(json), nil
}
