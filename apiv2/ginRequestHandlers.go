package apiv2

import (
	"encoding/base64"
	"encoding/csv"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vrypan/farma/config"
	"github.com/vrypan/farma/models"
)

// Check that the frameId passed in the parameters match ACCESS_FRAME_ID
// set during the signature verification. When ACL==ACL_ADMIN, ACCESS_FRAME_ID
// is not set, so this returns true for any frameId.
func validateAdminAccess(c *gin.Context) bool {
	if acl, ok := c.Get("ACL"); !ok || acl != ACL_ADMIN {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not admin key"})
		return false
	}
	return true
}

// Check that the frameId passed in the parameters match ACCESS_FRAME_ID
// set during the signature verification. When ACL==ACL_ADMIN, ACCESS_FRAME_ID
// is not set, so this returns true for any frameId.
func validateFrameAccess(c *gin.Context) bool {
	frameId := c.Param("frameId")
	if only_frame, is_set := c.Get("ACCESS_FRAME_ID"); is_set && frameId != only_frame {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return false
	}
	return true
}

// DONE
func H_FrameAdd(c *gin.Context) {
	// No validation, this handler is only accessible to ACL_ADMIN
	// which was already validated during the signature verification.
	var requestBody struct {
		Name   string `json:"name"`
		Domain string `json:"domain"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	frame := models.NewFrame()
	frame.Name = requestBody.Name
	frame.Domain = requestBody.Domain
	privKey, err := frame.PublicKey.GenerateKey(frame.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	encodedPrivKey := base64.StdEncoding.EncodeToString([]byte(privKey))

	if err := frame.Save(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"frame":       frame,
		"private_key": encodedPrivKey,
		"public_key":  frame.PublicKey.Encode(),
	})
}

func H_FrameUpdate(c *gin.Context) {
	if !validateFrameAccess(c) {
		return
	}
	frameId := c.Param("frameId")

	var requestBody struct {
		Name      string `json:"name"`
		Domain    string `json:"domain"`
		PublicKey string `json:"public_key"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	frame := models.NewFrame()
	frame = frame.FromId(frameId)
	if frame == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "frame not found"})
		return
	}

	if requestBody.Name != "" {
		frame.Name = requestBody.Name
	}
	if requestBody.Domain != "" {
		frame.Domain = requestBody.Domain
	}

	if requestBody.PublicKey != "" {
		publicKeyBytes, err := base64.StdEncoding.DecodeString(requestBody.PublicKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid base64 public_key"})
			return
		}
		frame.PublicKey = &models.PubKey{
			FrameId: frame.Id,
			Key:     publicKeyBytes,
		}
	}
	if err := frame.Save(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"frame":      frame,
		"public_key": frame.PublicKey.Encode(),
	})
}

func H_SubscriptionsGet(c *gin.Context) {
	if !validateFrameAccess(c) {
		return
	}
	frameId := c.Param("frameId")[1:]
	var prefix string

	if acl, ok := c.Get("ACL"); ok && acl == ACL_ADMIN && frameId == "" {
		// allow admins to query all subscriptions, regardless of frame
		prefix = "s:id:"
	} else {
		prefix = "s:id:" + frameId + ":"
	}

	getData(c, prefix, &models.Subscription{})
}

func H_FramesGet(c *gin.Context) {
	if !validateFrameAccess(c) {
		return
	}
	frameId := c.Param("frameId")
	prefix := "f:id:" + frameId
	getData(c, prefix, &models.Frame{})
	//debug(c)
}

func H_FramesGetAll(c *gin.Context) {
	if !validateFrameAccess(c) {
		return
	}
	prefix := "f:id:"
	getData(c, prefix, &models.Frame{})
}

func H_LogsGet(c *gin.Context) {
	if !validateFrameAccess(c) {
		return
	}
	frameId := c.Param("frameId")
	userId := c.Param("userId")[1:]
	prefix := "l:user:" + frameId + ":"
	if userId != "" {
		prefix += userId + ":"
	}
	getData(c, prefix, &models.UserLog{})
}

func H_Notify(c *gin.Context) {
	if !validateFrameAccess(c) {
		return
	}
	var ver int
	var err error
	var requestBody struct {
		FrameId string   `json:"frameId"`
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
	if frame.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "FRAME NOT FOUND"})
		return
	}

	if requestBody.Url == "" {
		requestBody.Url = "https://" + frame.Domain
	}

	// [url][token][fid]
	keys := make(map[string]map[string]uint64)
	// Map notification URLs to Application Ids
	appUrls := make(map[string]uint64)
	if len(requestBody.UserIds) > 0 {
		// Notify only specific subscribers
		for _, userId := range requestBody.UserIds {
			subscriptions, err := models.SubscriptionsByFrameUser(requestBody.FrameId, userId)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
			for _, s := range subscriptions {
				if s.Status == models.SubscriptionStatus_SUBSCRIBED || s.Status == models.SubscriptionStatus_RATE_LIMITED {
					if url, exists := keys[s.Url]; exists {
						url[s.Token] = s.UserId
					} else {
						keys[s.Url] = make(map[string]uint64)
						keys[s.Url][s.Token] = s.UserId
						appUrls[s.Url] = s.AppId
					}
				}
			}
		}
	} else {
		// Notify all frame subscribers
		var start []byte
		for {
			subscriptions, next, err := models.SubscriptionsByFrame(requestBody.FrameId, start, 1000)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
			for _, s := range subscriptions {
				if s.Status == models.SubscriptionStatus_SUBSCRIBED || s.Status == models.SubscriptionStatus_RATE_LIMITED {
					if url, exists := keys[s.Url]; exists {
						url[s.Token] = s.UserId
					} else {
						keys[s.Url] = make(map[string]uint64)
						keys[s.Url][s.Token] = s.UserId
						appUrls[s.Url] = s.AppId
					}
				}
			}
			if len(subscriptions) < 1000 {
				break
			}
			start = next
		}
	}

	notificationId := ""
	notificationCount := 0
	for url, tokens := range keys {
		// Calling notification.Send() for each client endpoint
		// One client may have multiple endpoints, and different
		// clients will have different endpoints.
		notification := models.NewNotification(
			requestBody.FrameId,
			appUrls[url],
			notificationId,
			requestBody.Title,
			requestBody.Body,
			requestBody.Url,
			url,
			keys[url],
		)
		notificationId = notification.Id
		notificationCount += len(tokens)
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

// DONE
func H_NotificationsGet(c *gin.Context) {
	if !validateFrameAccess(c) {
		return
	}
	frameId := c.Param("frameId")
	notificationId := c.Param("notificationId")
	if strings.HasPrefix(notificationId, "/") {
		notificationId = notificationId[1:]
	}
	if strings.HasSuffix(notificationId, "/") {
		notificationId = notificationId[:len(notificationId)-1]
	}
	prefix := "n:id:" + frameId + ":"
	if notificationId != "" {
		prefix += notificationId + ":"
	}
	getData(c, prefix, &models.Notification{})
}

// DONE
func H_DbKeysGet(c *gin.Context) {
	// No validation, this handler is only accessible to ACL_ADMIN
	// which was already validated during the signature verification.
	prefix := c.Param("prefix")[1:]
	getKeys(c, prefix)
}

// DONE
func H_Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": config.FARMA_VERSION,
	})
}

func H_NewKeypair(c *gin.Context) {
	frameId := c.Param("frameId")
	k := models.PubKey{}
	privKey, err := k.GenerateKey(frameId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	encodedPrivKey := base64.StdEncoding.EncodeToString([]byte(privKey))
	c.JSON(http.StatusCreated, gin.H{
		"private_key": encodedPrivKey,
		"public_key":  k.Encode(),
	})
}

func debug(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"debug": gin.H{
			"keys":    c.Keys,
			"headers": c.Request.Header,
			"params":  c.Params,
		},
	})
}

func H_SubscriptionsImportCSV(c *gin.Context) {
	if !validateFrameAccess(c) {
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	frameId := c.Param("frameId")
	appUrl := c.PostForm("appUrl")
	if appUrl == "" {
		appUrl = "https://api.warpcast.com/v1/frame-notifications"
	}
	appIdStr := c.PostForm("appId")
	if appIdStr == "" {
		appIdStr = "9152"
	}
	appId, err := strconv.ParseUint(appIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appId"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not open file"})
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	var entries []csvEntry

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not read CSV line"})
			return
		}
		if line[0] == "fid" {
			// Skip the header line
			continue
		}
		id, err := strconv.ParseUint(line[0], 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not convert CSV line[0] to uint64"})
			return
		}
		entries = append(entries, csvEntry{
			fid:               id,
			notificationToken: line[1],
		})
	}

	importData := ImportData{
		frameId: frameId,
		appId:   appId,
		appUrl:  appUrl,
		data:    entries,
	}

	imported := 0
	if imported, err = CreateSubscriptionsFromCSV(importData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": gin.H{
			"message": "Subscriptions imported successfully",
			"entries": imported,
		},
	})
}
