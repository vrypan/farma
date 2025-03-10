package apiv2

import (
	"bytes"
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
				return isValid
			}
		}
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Verification failed"})
	return
}

func VerifySignature(acl ACL) gin.HandlerFunc {
	adminKey := &models.PubKey{}
	encodedAdminKey := "0:" + config.GetString("key.public")
	if err := adminKey.Decode(encodedAdminKey); err != nil {
		panic("Unable to decode admin key")
	}
	return func(c *gin.Context) {
		c.Set("ACL", acl)
		pubKey := &models.PubKey{}
		if acl == ACL_ADMIN {
			pubKey = adminKey
		}
		if acl == ACL_FRAME_OR_ADMIN {
			pubKeyHeader := c.GetHeader("X-Public-Key")
			if err := pubKey.Decode(pubKeyHeader); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error":        "Unable to decode public key",
					"Decode()":     err.Error(),
					"X-Public-Key": pubKeyHeader,
				})
				return
			}
			if bytes.Equal(pubKey.Key, adminKey.Key) {
				c.Set("ACL", ACL_ADMIN)
			} else {
				if !pubKey.InDb() {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"error": "Invalid public key",
					})
					return
				}
				c.Set("ACCESS_FRAME_ID", pubKey.FrameId)
			}
		}
		if isValid := verify(c, pubKey); isValid {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":          "Verification failed",
				"RequestHeaders": c.Request.Header,
				"PublicKey":      pubKey,
			})
		}
	}
}
