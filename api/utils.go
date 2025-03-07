package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vrypan/farma/config"
	db "github.com/vrypan/farma/localdb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ApiCallData struct {
	FrameId   int
	PublicKey string
	Method    string
	Endpoint  string
	Path      string
	Id        string
	Body      string
	RawQuery  string
}

func (api *ApiCallData) Call() ([]byte, error) {
	endpoint := api.Path

	if api.Endpoint == "" {
		endpoint = fmt.Sprintf("http://%s/api/v1/%s", config.GetString("host.addr"), api.Path)
	} else {
		endpoint += api.Path
	}
	endpoint = fmt.Sprintf("%s%s", endpoint, api.Id)

	keyPrivate := config.GetString("key.private")
	if keyPrivate == "" {
		return nil,
			fmt.Errorf("Missing private key. Use command '$ farma config set key.private <private_key>' or set the environment variable FARMA_KEY_PRIVATE.")

	}

	keyPrivateBytes, err := base64.RawStdEncoding.DecodeString(keyPrivate)
	if err != nil {
		return nil, fmt.Errorf("Error decoding private key: %v\n", err)

	}

	hostUrl, err := url.Parse(api.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("Error parsing endpoint: %v\n", err)

	}

	request := Request{
		Method: api.Method,
		Path:   hostUrl.Path,
		Body:   []byte(api.Body),
		Query:  api.RawQuery,
	}
	request.SignEd25519(keyPrivateBytes)

	res, err := request.Send(fmt.Sprintf("%s://%s", hostUrl.Scheme, hostUrl.Host))
	if err != nil {
		return res, fmt.Errorf("Error sending request: %v\n", err)
	}
	return res, nil
}

func retrieveData(c *gin.Context, prefixKey string, prefixVal string, model any) {
	if strings.HasPrefix(prefixVal, "/") {
		prefixVal = strings.TrimPrefix(prefixVal, "/")
	}
	var prefix []byte
	if prefixVal == "" {
		prefix = []byte(prefixKey)
	} else {
		prefix = []byte(prefixKey + prefixVal + ":")
	}

	var start []byte
	var err error
	if s := c.DefaultQuery("start", ""); s != "" {
		start, err = base64.StdEncoding.DecodeString(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to decode start value"})
			return
		}
	} else {
		start = prefix
	}

	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	fmt.Println("prefix:", string(prefix), "start:", string(start), "limit:", limit)
	data, next, err := db.GetPrefixP(prefix, start, limit)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	list := make([]json.RawMessage, len(data))

	for i, item := range data {
		err := proto.Unmarshal(item, model.(proto.Message))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		j, err := protojson.Marshal(model.(proto.Message))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		list[i] = j
	}

	c.JSON(http.StatusOK, gin.H{
		"result": list,
		"next":   next,
	})
}

func OnlyFrame(c *gin.Context, frameId string) error {
	of, exists := c.Get("OnlyFrame")
	if exists {
		return nil
	}
	if frameId == of.(string) {
		return nil
	}
	return errors.New("Frame ID and X-Public-Key mismatch")
}
