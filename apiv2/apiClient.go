package apiv2

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vrypan/farma/config"
)

type ApiResult struct {
	Error  string            `json:"error"`
	Result []json.RawMessage `json:"result"`
	Next   string            `json:"next"`
}

type ApiClient struct {
	HttpPath   string
	HttpMethod string
	HttpBody   []byte
	//QueryString string
	Payload    []byte
	PrivateKey []byte
	publicKey  string
	//FrameId    string
}

func (api ApiClient) Init(
	httpMethod string,
	httpPath string,
	payload []byte,
	privateKey string,
	frameId string,
) ApiClient {
	if strings.HasPrefix(httpPath, "http") {
		api.HttpPath = httpPath
	} else {
		api.HttpPath = fmt.Sprintf("http://%s/api/v2/%s", config.GetString("host.addr"), httpPath)
	}
	api.HttpMethod = httpMethod
	api.Payload = payload

	if privateKey != "" && privateKey != "config" {
		privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKey)
		if err != nil {
			panic(err)
		}
		pubKeyBytes := ed25519.PrivateKey(privateKeyBytes).Public().(ed25519.PublicKey)
		api.PrivateKey = privateKeyBytes
		api.publicKey = frameId + ":" + base64.StdEncoding.EncodeToString(pubKeyBytes)
	}
	if string(privateKey) == "config" {
		pkStr := config.GetString("key.private")
		privateKeyBytes, err := base64.StdEncoding.DecodeString(pkStr)
		if err != nil {
			panic(err)
		}
		api.PrivateKey = privateKeyBytes
		pubKeyBytes := ed25519.PrivateKey(privateKeyBytes).Public().(ed25519.PublicKey)
		api.publicKey = frameId + ":" + base64.StdEncoding.EncodeToString(pubKeyBytes)
	}
	return api
}

func (api *ApiClient) Request(start string, limit string) ([]byte, error) {
	requestUrl, _ := url.Parse(api.HttpPath)
	qParams := requestUrl.Query()
	if limit != "" {
		qParams.Add("limit", limit)
	}
	if start != "" {
		qParams.Add("start", start)
	}
	requestUrl.RawQuery = qParams.Encode()

	client := &http.Client{}
	req, err := http.NewRequest(api.HttpMethod, requestUrl.String(), bytes.NewBuffer(api.Payload))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	date := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("X-Date", date)

	if api.PrivateKey != nil {
		signature := ed25519.Sign(api.PrivateKey,
			[]byte(api.HttpMethod+"\n"+requestUrl.Path+"\n"+date),
		)
		sig := base64.StdEncoding.EncodeToString(signature)
		req.Header.Set("X-Signature", sig)
		req.Header.Set("X-Public-Key", api.publicKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if (resp.StatusCode != http.StatusOK) && (resp.StatusCode != http.StatusCreated) {
		return body, fmt.Errorf("Server returned status code %d", resp.StatusCode)
	}
	return body, nil

}
