package api

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
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

func ApiCall(method string, endpoint string, methodPath string, id string, body []byte, rawQuery string) ([]byte, error) {
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
		Query:  rawQuery,
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
