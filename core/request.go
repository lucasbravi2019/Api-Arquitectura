package core

import (
	"encoding/json"
	"log"
	"net/http"
)

func DecodeBody(r *http.Request, storeVar any) bool {
	err := json.NewDecoder(r.Body).Decode(&storeVar)
	if err != nil {
		log.Println(err.Error())
		return true
	}

	return Validate(storeVar)
}
