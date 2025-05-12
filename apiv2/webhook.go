package apiv2

import (
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/vrypan/farma/fctools"
	"github.com/vrypan/farma/models"
)

var validPathRegex = regexp.MustCompile(`^[\w/-_]*$`)

func isValidPath(path string) bool {
	return validPathRegex.MatchString(path)
}

func WebhookHandler(hub *fctools.FarcasterHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// These are public endpoints that can and will be abused.
		// Let's make sure that HTTP requests are within some reasonable limits.
		if c.Request.ContentLength > 1024 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Content Length > 1024"})
			return
		}
		if len(c.Request.URL.Path) > 128 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Path Length > 128"})
			return
		}
		if !isValidPath(c.Request.URL.Path) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Path contains invalid_characters"})
			return
		}

		frame := models.NewFrame()
		if err := frame.FromEndpoint(c.Request.URL.Path); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Unknown endpoint"})
			return
		}
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"error": "Error reading request body"})
			return
		}

		subscription, eventType := models.NewSubscription().FromHttpEvent(body)
		subscription.VerifyAppId(hub)
		subscription.FrameId = frame.Id
		if err = subscription.Save(); err != nil {
			log.Printf("Error updating db: %v\nSubscription details: %v", err, subscription.NiceString())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error updating db"})
			return
		}

		ulog := models.UserLog{
			FrameId: subscription.FrameId,
			UserId:  subscription.UserId,
			AppId:   subscription.AppId,
			EvtType: eventType,
		}
		switch eventType {
		case models.EventType_FRAME_ADDED, models.EventType_NOTIFICATIONS_ENABLED:
			subscriptionContext := models.EventContextSubscription{
				Url:   subscription.Url,
				Token: subscription.Token,
			}
			ulog.EvtContext = &models.UserLog_EventContextSubscription{EventContextSubscription: &subscriptionContext}
		case models.EventType_FRAME_REMOVED, models.EventType_NOTIFICATIONS_DISABLED:
			subscriptionContext := models.EventContextSubscription{
				Url:   "",
				Token: "",
			}
			ulog.EvtContext = &models.UserLog_EventContextSubscription{EventContextSubscription: &subscriptionContext}
		default:
			ulog.EvtContext = &models.UserLog_EventContextNone{}
		}
		err = ulog.Save()
		if err != nil {
			log.Printf("Error saving user log: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error saving user log"})
			return
		}
		c.Status(http.StatusOK)
	}
}
