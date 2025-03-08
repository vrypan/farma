package apiv2

import (
	"encoding/base64"
	"net/http"

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
		"frame":      frame,
		"PrivateKey": encodedPrivKey,
		"PublicKey":  frame.PublicKey.Encode(),
	})
}

// DONE
func H_FrameUpdate(c *gin.Context) {
	if !validateFrameAccess(c) {
		return
	}
	frameId := c.Param("id")

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

	if frame.PublicKey != nil {
		if err := frame.PublicKey.Delete(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete f:pk: " + err.Error(),
			})
			return
		}
	}

	// Now update the frame data and save.

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
		if err := frame.PublicKey.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save " + frame.PublicKey.DbKey() + " " + err.Error(),
			})
			return
		}
	}
	if err := frame.Save(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"frame":  frame,
		"pubKey": frame.PublicKey.Encode(),
	})
}

// DONE
func H_SubscriptionsGet(c *gin.Context) {
	if !validateFrameAccess(c) {
		return
	}
	frameId := c.Param("frameId")[1:]
	prefix := "s:id:" + frameId + ":"
	getData(c, prefix, &models.Subscription{})
}

func H_FramesGet(c *gin.Context) {
	if !validateFrameAccess(c) {
		return
	}
}

// DONE
func H_LogsGet(c *gin.Context) {
	if !validateFrameAccess(c) {
		return
	}
	frameId := c.Param("frameId")[1:]
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
}

// DONE
func H_NotificationsGet(c *gin.Context) {
	if !validateFrameAccess(c) {
		return
	}
	if !validateFrameAccess(c) {
		return
	}
	frameId := c.Param("frameId")[1:]
	notificationId := c.Param("notificationId")[1:]
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
	getData(c, prefix, []byte{})
}

// DONE
func H_Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": config.FARMA_VERSION,
	})
}
