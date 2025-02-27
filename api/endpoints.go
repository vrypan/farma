package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vrypan/farma/config"
	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/models"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func H_FrameAdd(c *gin.Context) {
	var requestBody struct {
		Name    string `json:"name"`
		Domain  string `json:"domain"`
		Webhook string `json:"webhook"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	frame := models.NewFrame()
	err := frame.FromName(requestBody.Name)
	if err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "frame exists",
		})
		return
	}
	if err != db.ERR_NOT_FOUND {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	frame.Name = requestBody.Name
	frame.Domain = requestBody.Domain
	if requestBody.Webhook == "" {
		frame.Webhook = "/f/" + uuid.New().String()
	} else {
		frame.Webhook = requestBody.Webhook
	}

	if err := frame.Save(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusCreated, frame)
}

func H_SubscriptionsGet(c *gin.Context) {
	frameId := c.Param("frameId")[1:]
	var prefix []byte
	if frameId == "" {
		prefix = []byte("s:id:")
	} else {
		prefix = []byte("s:id:" + frameId + ":")
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
	data, next, err := db.GetPrefixP(prefix, start, limit)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	list := make([]json.RawMessage, len(data))
	for i, item := range data {
		var pb models.Subscription
		proto.Unmarshal(item, &pb)
		j, err := protojson.Marshal(&pb)
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
func H_FramesGet(c *gin.Context) {
	frameId := c.Param("id")[1:]
	var prefix []byte
	if frameId == "" {
		prefix = []byte("f:id:")
	} else {
		prefix = []byte("f:id:" + frameId + ":")
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

	data, next, err := db.GetPrefixP(prefix, start, limit)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	list := make([]json.RawMessage, len(data))

	for i, item := range data {
		var pb models.Frame
		err := proto.Unmarshal(item, &pb)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		j, err := protojson.Marshal(&pb)
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

func H_LogsGet(c *gin.Context) {
	userId := c.Param("userId")[1:]
	var prefix []byte
	if userId == "" {
		prefix = []byte("l:user:")
	} else {
		prefix = []byte("l:user:" + userId + ":")
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

	data, next, err := db.GetPrefixP(prefix, start, limit)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	list := make([]json.RawMessage, len(data))

	for i, item := range data {
		var pb models.UserLog
		err := proto.Unmarshal(item, &pb)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		j, err := protojson.Marshal(&pb)
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

func H_DbKeysGet(c *gin.Context) {
	prefix := []byte(c.Param("prefix")[1:])

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

	data, next, err := db.GetKeysWithPrefix(prefix, start, limit)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	list := make([]string, len(data))

	for i, key := range data {
		list[i] = string(key)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": list,
		"next":   next,
	})
}

func H_Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": config.FARMA_VERSION,
	})
}

func H_Notify(c *gin.Context) {
	var ver int
	var err error
	var requestBody struct {
		FrameId uint64   `json:"frameId"`
		Title   string   `json:"title"`
		Body    string   `json:"body"`
		Url     string   `json:"url"`
		UserIds []uint64 `json:"userIds"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	frame := models.NewFrame().FromId(requestBody.FrameId)
	if frame.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "FRAME NOT FOUND"})
		return
	}

	if requestBody.Url == "" {
		requestBody.Url = "https://" + frame.Domain
	}

	keys := make(map[string][][]byte)

	if len(requestBody.UserIds) == 0 {
		requestBody.UserIds = append(requestBody.UserIds, 0)
	}
	for _, userId := range requestBody.UserIds {
		prefix := []byte("s:url:" + strconv.Itoa(int(requestBody.FrameId)) + ":")
		if userId != 0 {
			prefix = append(prefix, strconv.Itoa(int(userId))+":"...)
		}

		startKey := prefix
		for {
			urlKeys, nextKey, err := db.GetKeysWithPrefix(prefix, startKey, 1000)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
			for _, urlKeyBytes := range urlKeys {
				urlKey := models.UrlKey{}.DecodeBytes(urlKeyBytes)
				if urlKey.Status == models.SubscriptionStatus_SUBSCRIBED || urlKey.Status == models.SubscriptionStatus_RATE_LIMITED {
					keys[urlKey.Endpoint] = append(keys[urlKey.Endpoint], urlKeyBytes)
				}
			}
			startKey = nextKey
			if len(urlKeys) < 1000 {
				break
			}
		}
	}

	notificationId := ""
	notificationCount := 0
	for url, urlKeys := range keys {
		notification := models.NewNotification(
			notificationId,
			requestBody.Title,
			requestBody.Body,
			requestBody.Url,
			url,
			urlKeys,
		)
		notificationId = notification.Id
		notificationCount += len(urlKeys)
		if err := notification.Send(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		ver, err = notification.Update()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"NotificationId":      notificationId,
		"NotificationVersion": ver,
		"Count":               notificationCount,
	})
}

func H_NotificationsGet(c *gin.Context) {
	notificationId := c.Param("id")[1:]
	var prefix []byte
	if notificationId == "" {
		prefix = []byte("n:id:")
	} else {
		prefix = []byte("n:id:" + notificationId + ":")
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
	data, next, err := db.GetPrefixP(prefix, start, limit)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	list := make([]json.RawMessage, len(data))
	for i, item := range data {
		var pb models.Notification
		proto.Unmarshal(item, &pb)
		j, err := protojson.Marshal(&pb)
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
