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

func isValidPath(path string) bool {
	matched, _ := regexp.MatchString(`^[\w/-_]*$`, path)

	return matched
}

func WebhookHandler(hub *fctools.FarcasterHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		//func NotificationsH(c *gin.Context, hub *fctools.FarcasterHub) {
		// These are public endpoints that can and will be abused.
		// Let's make sure that HTTP requests are within some reasonable limits.
		if c.Request.ContentLength > 1024 {
			c.AbortWithStatus(http.StatusBadRequest)
			c.String(http.StatusBadRequest, "Content Length > 1024")
			return
		}
		if len(c.Request.URL.Path) > 128 {
			c.AbortWithStatus(http.StatusBadRequest)
			c.String(http.StatusBadRequest, "Path Length > 128")
			return
		}
		if !isValidPath(c.Request.URL.Path) {
			c.AbortWithStatus(http.StatusBadRequest)
			c.String(http.StatusBadRequest, "Path contains invalid_characters")
			return
		}

		frame := models.NewFrame()
		if err := frame.FromEndpoint(c.Request.URL.Path); err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			c.String(http.StatusNotFound, "Unknown endpoint")
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusNoContent)
			c.String(http.StatusNoContent, "Error reading request body")
			return
		}

		subscription, eventType := models.NewSubscription().FromHttpEvent(body)
		subscription.VerifyAppId(hub)
		subscription.FrameId = frame.Id
		if err = subscription.Save(); err != nil {
			log.Println("Error updating db.", err)
			log.Println("Subscription details:", subscription.NiceString())
			c.AbortWithStatus(http.StatusInternalServerError)
			c.String(http.StatusInternalServerError, "Error updating db")
			return
		}
		ulog := models.UserLog{
			FrameId:    subscription.FrameId,
			UserId:     subscription.UserId,
			AppId:      subscription.AppId,
			EvtType:    eventType,
			EvtContext: &models.UserLog_EventContextNone{},
		}
		err = ulog.Save()
		c.Status(http.StatusOK)
	}
}
