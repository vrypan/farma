package api

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (r Response) String() string {
	output, err := json.Marshal(r)
	if err != nil {
		return Error("JSON_ERROR", err)
	}
	return string(output)
}

func Error(message string, e error) string {
	errResponse := Response{
		Status:  "error",
		Message: message,
		Data:    e,
	}
	errOutput, err := json.Marshal(errResponse)
	if err != nil {
		return "{}"
	}
	return string(errOutput)
}
func (r Response) Format(status, message string, data interface{}) string {
	r.Status = status
	r.Message = message
	r.Data = data
	output, err := json.Marshal(r)
	if err != nil {
		errResponse := Response{
			Status:  "error",
			Message: err.Error(),
		}
		errOutput, err := json.Marshal(errResponse)
		if err != nil {
			return "{}"
		}
		return string(errOutput)
	}
	return string(output)
}

type Api struct {
	PubKeys     map[string]int
	jsonHeader  map[string]interface{}
	jsonPayload map[string]interface{}
}

func New() *Api {
	return &Api{
		PubKeys:     make(map[string]int),
		jsonHeader:  make(map[string]interface{}),
		jsonPayload: make(map[string]interface{}),
	}
}

func (a *Api) AddKey(k string) {
	if !strings.HasPrefix(k, "0x") {
		k = "0x" + k
	}
	a.PubKeys[k] = 1
}

func (a *Api) Prepare(payload string) error {

	var jsonBody, jsonHeader map[string]interface{}

	if err := json.Unmarshal([]byte(payload), &jsonBody); err != nil {
		return err
	}

	header, _ := jsonBody["header"].(string)
	headerBytes, err := base64.RawURLEncoding.DecodeString(header)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(headerBytes, &jsonHeader); err != nil {
		return err
	}

	requestKey, _ := jsonHeader["key"].(string)

	if !strings.HasPrefix(requestKey, "0x") {
		requestKey = "0x" + requestKey
	}

	if _, ok := a.PubKeys[requestKey]; !ok {
		return fmt.Errorf("Invalid key: %s", requestKey)
	}

	signedData, _ := jsonBody["payload"].(string)
	signedDataBytes := []byte(header + "." + signedData)

	signature, _ := jsonBody["signature"].(string)
	signatureBytes, _ := base64.RawURLEncoding.DecodeString(signature)

	if !ed25519.Verify(common.FromHex(requestKey), signedDataBytes, signatureBytes) {
		return fmt.Errorf("Invalid signature")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(signedData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(payloadBytes, &a.jsonPayload)
	if err != nil {
		return err
	}

	a.jsonHeader = jsonHeader

	// All checkc have passed. The signature has been validated.
	// jsonPayload and jsonHeader are loaded with the corresponding JSON data.
	// The API consumer can now call Execute() to run the command.
	return nil
}

func (a Api) Execute() (string, error) {
	switch a.jsonPayload["command"] {
	case "notification/send":
		params := a.jsonPayload["params"].(map[string]interface{})
		r := Notify(
			params["frame"].(string),
			params["title"].(string),
			params["body"].(string),
			params["url"].(string),
		)
		return r, nil
	case "frames/get":
		params := a.jsonPayload["params"].(map[string]interface{})
		if params["id"] != nil {
			id := uint64(params["id"].(float64))
			return FrameGet(id), nil
		}
		return FramesGet(), nil
	case "frames/add":
		params := a.jsonPayload["params"].(map[string]interface{})
		if params != nil && params["name"] != nil && params["domain"] != nil && params["webhook"] != nil {
			return FrameAdd(params["name"].(string), params["domain"].(string), params["webhook"].(string)), nil
		} else {
			return Error("FRAME_MISSING_PARAMS", nil), nil
		}
	case "db/keys":
		return DbKeys(), nil
	case "logs/get":
		userId := 0
		limit := 1000
		params := a.jsonPayload["params"].(map[string]interface{})
		if params != nil && params["userId"] != nil && params["limit"] != nil {
			userId = params["name"].(int)
			limit = params["name"].(int)
		}
		return ShowLogs(uint64(userId), limit), nil
	case "subscriptions/get":
		frameId := 0
		limit := 1000
		params := a.jsonPayload["params"].(map[string]interface{})
		if params != nil && params["frameId"] != nil && params["limit"] != nil {
			frameId = params["frameId"].(int)
			limit = params["limit"].(int)
		}
		return ShowSubscriptions(uint64(frameId), limit), nil
	default:
		return Error("INVALID_COMMAND", nil), nil
	}
}
