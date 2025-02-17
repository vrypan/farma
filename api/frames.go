package api

import (
	"encoding/json"
	"fmt"

	"github.com/vrypan/farma/utils"
)

func GetFrames() string {
	frames := utils.AllFrames()
	response := Response{}
	response.Data = frames
	response.Status = "success"
	response.Message = "Frames retrieved successfully"
	output, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Error converting output to json: %v\n", err)
		return ""
	}
	return string(output)
}

func GetFrame(id uint64) string {
	frame := utils.NewFrame().FromId(id)
	output, err := json.Marshal(frame)
	if err != nil {
		fmt.Printf("Error converting output to json: %v\n", err)
		return ""
	}
	return string(output)
}
