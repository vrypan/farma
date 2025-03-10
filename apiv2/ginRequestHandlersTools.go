package apiv2

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	db "github.com/vrypan/farma/localdb"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func getData(c *gin.Context, prefix string, model any) {
	if strings.HasPrefix(prefix, "/") {
		prefix = strings.TrimPrefix(prefix, "/")
	}
	prefixBytes := []byte(prefix)

	var startBytes []byte
	var err error
	if s := c.DefaultQuery("start", ""); s != "" {
		startBytes, err = base64.StdEncoding.DecodeString(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to decode start value"})
			return
		}
	} else {
		startBytes = prefixBytes
	}

	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	//fmt.Printf("GetPrefixP(%s, %s, %d)\n", prefixBytes, startBytes, limit)
	data, next, err := db.GetPrefixP(prefixBytes, startBytes, limit)

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

func getKeys(c *gin.Context, prefix string) {
	if strings.HasPrefix(prefix, "/") {
		prefix = strings.TrimPrefix(prefix, "/")
	}
	prefixBytes := []byte(prefix)

	var startBytes []byte
	var err error
	if s := c.DefaultQuery("start", ""); s != "" {
		startBytes, err = base64.StdEncoding.DecodeString(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to decode start value"})
			return
		}
	} else {
		startBytes = prefixBytes
	}

	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	data, next, err := db.GetKeysWithPrefix(prefixBytes, startBytes, limit)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	list := make([]string, len(data))

	for i, item := range data {
		list[i] = string(item)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": list,
		"next":   next,
	})
}
