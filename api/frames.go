package api

import (
	"encoding/json"

	"github.com/google/uuid"
	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/utils"
)

func FramesGet() string {
	frames := utils.AllFrames()
	response := Response{}
	response.Data = frames
	response.Status = "SUCCESS"
	response.Message = "Frames retrieved successfully"
	output, err := json.Marshal(response)
	if err != nil {
		return response.Format("error", "ERROR_JSON_MARSHAL", err)
	}
	return string(output)
}

func FrameGet(id uint64) string {
	frame := utils.NewFrame().FromId(id)
	response := Response{}
	response.Data = frame
	if frame == nil {
		response.Status = "ERROR"
		response.Message = "FRAME_NOT_FOUND"
		output, err := json.Marshal(response)
		if err != nil {
			return response.Format("error", "ERROR_JSON_MARSHAL", err)
		}
		return string(output)
	}
	response.Status = "SUCCESS"
	response.Message = "Frame retrieved successfully"
	output, err := json.Marshal(response)
	if err != nil {
		return response.Format("ERROR", "ERROR_JSON_MARSHAL", err)
	}
	return string(output)
}

func FrameAdd(frameName, frameDomain, frameWebhook string) string {
	response := Response{}
	if frameWebhook == "" {
		frameWebhook = "/f/" + uuid.New().String()
	}
	if len(frameName) > 32 {
		return response.Format("ERROR", "FRAME_NAME_TOO_LONG",
			struct{ description string }{"Frame name must be up to 32 characters"})
	}
	frame := utils.NewFrame()
	err := frame.FromName(frameName)
	if err == nil {
		return response.Format("ERROR", "FRAME_EXISTS", nil)
	}
	if err != db.ERR_NOT_FOUND {
		return response.Format("ERROR", "FRAME_ERROR", err)
	}
	frame.Name = frameName
	frame.Domain = frameDomain
	frame.Webhook = frameWebhook
	if err := frame.Save(); err != nil {
		return response.Format("ERROR", "FRAME_ERROR", err)
	}
	response.Data = frame
	response.Status = "SUCCESS"
	response.Message = "Frame added successfully"
	output, err := json.Marshal(response)
	if err != nil {
		return response.Format("ERROR", "ERROR_JSON_MARSHAL", err)
	}
	return string(output)
}
