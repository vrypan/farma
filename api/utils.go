package api

import (
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/vrypan/farma/config"
)

func ApiCall(method string, endpoint string, methodPath string, id string, body []byte) ([]byte, error) {
	apiEndpointPath := methodPath

	if endpoint == "" {
		endpoint = fmt.Sprintf("http://%s/api/v1/%s", config.GetString("host.addr"), apiEndpointPath)
	} else {
		endpoint += apiEndpointPath
	}

	endpoint = fmt.Sprintf("%s%s", endpoint, id)

	keyPrivate := config.GetString("key.private")
	if keyPrivate == "" {
		return nil,
			fmt.Errorf("Missing private key. Use command '$ farma config set key.private <private_key>' or set the environment variable FARMA_KEY_PRIVATE.")

	}

	keyPrivateBytes, err := hex.DecodeString(keyPrivate[2:])
	if err != nil {
		return nil, fmt.Errorf("Error decoding private key: %v\n", err)

	}

	hostUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("Error parsing endpoint: %v\n", err)

	}

	request := Request{
		Method: method,
		Path:   hostUrl.Path,
		Body:   body,
	}
	request.Sign(keyPrivateBytes)

	res, err := request.Send(fmt.Sprintf("%s://%s", hostUrl.Scheme, hostUrl.Host))
	if err != nil {
		return res, fmt.Errorf("Error sending request: %v\n", err)
	}
	return res, nil
}
