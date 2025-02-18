package api

import (
	"fmt"

	db "github.com/vrypan/farma/localdb"
)

func DbKeys() string {
	response := Response{}
	keyStrings := []string{}

	keys, err := db.GetKeys([]byte(""), 1000)
	if err != nil {
		fmt.Println(err)
	}
	for _, k := range keys {
		keyStrings = append(keyStrings, string(k))
	}
	return response.Format("OK", "", keyStrings)
}
