package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/utils"
	"google.golang.org/protobuf/proto"
)

func H_FramesGet(c *gin.Context) {
	idStr := c.Param("id")[1:]
	if idStr == "" {
		frames := utils.AllFrames()
		c.JSON(http.StatusOK, frames)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_ID: " + idStr})
		return
	}

	frame := utils.NewFrame().FromId(id)
	frames := append([]*utils.Frame{}, frame)
	if frame == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FRAME_NOT_FOUND"})
		return
	}
	c.JSON(http.StatusOK, frames)
}

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
	frame := utils.NewFrame()
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
	var prefix string
	if frameId == "" {
		prefix = "s:id:"
	} else {
		prefix = "s:id:" + frameId + ":"
	}

	limitStr := c.DefaultQuery("limit", "1000")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	data, _, err := db.GetPrefixP([]byte(prefix), []byte(prefix), limit)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	list := make([]*utils.Subscription, len(data))
	for i, item := range data {
		var pb utils.Subscription
		proto.Unmarshal(item, &pb)
		list[i] = &pb
	}
	c.JSON(http.StatusOK, list)
}

func H_LogsGet(c *gin.Context) {
	userId := c.Param("userId")[1:]
	var prefix string
	if userId == "" {
		prefix = "l:user:"
	} else {
		prefix = "l:user:" + userId + ":"
	}

	limitStr := c.DefaultQuery("limit", "1000")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	data, _, err := db.GetPrefixP([]byte(prefix), []byte(prefix), limit)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	list := make([]*utils.UserLog, len(data))
	for i, item := range data {
		var pb utils.UserLog
		proto.Unmarshal(item, &pb)
		list[i] = &pb
	}
	c.JSON(http.StatusOK, list)
}

func H_Notify(c *gin.Context) {
	var requestBody struct {
		FrameId uint64 `json:"frameId"`
		Title   string `json:"title"`
		Body    string `json:"body"`
		Url     string `json:"url"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	frame := utils.NewFrame().FromId(requestBody.FrameId)
	if frame.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "FRAME NOT FOUND"})
		return
	}

	// Warpcast will crash when an notificationUrl is clicked.
	if requestBody.Url == "" {
		requestBody.Url = "https://" + frame.Domain
	}

	keys := make(map[string][][]byte)

	prefix := []byte("s:url:" + strconv.Itoa(int(requestBody.FrameId)) + ":")

	startKey := prefix
	for {
		urlKeys, nextKey, err := db.GetKeysWithPrefix(prefix, startKey, 1000)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		for _, urlKeyBytes := range urlKeys {
			urlKey := utils.UrlKey{}.DecodeBytes(urlKeyBytes)
			status := urlKey.Status
			url := urlKey.Endpoint
			if status == utils.SubscriptionStatus_SUBSCRIBED || status == utils.SubscriptionStatus_RATE_LIMITED {
				keys[url] = append(keys[url], urlKeyBytes)
			}
		}
		startKey = nextKey
		if len(urlKeys) < 1000 {
			break
		}
	}

	notificationId := ""
	notificationCount := 0
	for url, urlKeys := range keys {
		notification := utils.NewNotification(
			notificationId,
			requestBody.Title,
			requestBody.Body,
			requestBody.Url,
			url,
			urlKeys,
		)
		notificationId = notification.Id
		notificationCount += len(urlKeys)
		err := notification.Send()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}
	responseJson := struct {
		NotificationId string
		Count          int
	}{
		NotificationId: notificationId,
		Count:          notificationCount,
	}
	c.JSON(http.StatusOK, responseJson)
}
