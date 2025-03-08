package apiv2

import (
	"crypto/ed25519"
	"encoding/base64"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vrypan/farma/config"
	"github.com/vrypan/farma/models"
)

type ACL int

const (
	ACL_ADMIN ACL = iota
	ACL_FRAME_OR_ADMIN
)

func verify(c *gin.Context, pubKey *models.PubKey) (isValid bool) {
	reqMethod, reqPath := c.Request.Method, c.Request.URL.Path
	date, signature := c.GetHeader("Date"), c.GetHeader("X-Signature")

	if edPubKey := ed25519.PublicKey(pubKey.Key); edPubKey != nil {
		procTime, _ := time.Parse(time.RFC1123, date)
		if int(math.Abs(float64(time.Now().UTC().Sub(procTime).Seconds()))) <= 10 {
			if sig, _ := base64.StdEncoding.DecodeString(signature); len(sig) > 0 {
				isValid = ed25519.Verify(edPubKey, []byte(reqMethod+"\n"+reqPath+"\n"+date), sig)
				return
			}
		}
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Verification failed"})
	return
}

func VerifySignature(acl ACL) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ACL", acl)
		pubKey := &models.PubKey{}
		if acl == ACL_ADMIN {
			encodedAdminKey := config.GetString("key.public")
			if pubKey.Decode(encodedAdminKey) == nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unable to decode admin key"})
				return
			}
		} else {
			c.Set("ACCESS_FRAME_ID", pubKey.FrameId)
			if pubKey.Decode(c.GetHeader("X-Public-Key")) == nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unable to decode public key"})
				return
			}
			if !pubKey.Exists() {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid public key",
				})
				return
			}
		}
		if isValid := verify(c, pubKey); isValid {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Verification failed"})
		}
	}
}
