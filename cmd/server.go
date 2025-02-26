package cmd

import (
	"encoding/hex"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"crypto/ed25519"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/vrypan/farma/api"
	"github.com/vrypan/farma/config"
	"github.com/vrypan/farma/fctools"
	db "github.com/vrypan/farma/localdb"
)

var ginServerCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Run:   ginServer,
}

func init() {
	rootCmd.AddCommand(ginServerCmd)
	ginServerCmd.Flags().StringP("address", "a", "", "Listen on this address/port.")
	ginServerCmd.Flags().BoolP("verbose", "v", false, "Log additional info.")
}

func verifySignature() gin.HandlerFunc {
	keyHex := config.GetString("key.public")
	pubKey, err := hex.DecodeString(keyHex[2:])
	if err != nil {
		log.Fatalf("Invalid key: %v", err)
	}

	return func(c *gin.Context) {
		rMethod := c.Request.Method
		rPath := c.Request.URL.Path
		rDate := c.GetHeader("Date")
		rSignature := c.GetHeader("X-Signature")

		parsedTime, err := time.Parse(time.RFC1123, rDate)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Error parsing Date",
			})
			return
		}
		now := time.Now().UTC()
		diffSeconds := int(math.Abs(float64(now.Sub(parsedTime).Seconds())))
		if diffSeconds > 10 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Date diff more than 10 seconds",
			})
			return
		}
		signatureBytes, err := hex.DecodeString(rSignature)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Error decoding signature",
			})
			return
		}
		signedData := []byte(rMethod + "\n" + rPath + "\n" + rDate)
		if isValidSig := ed25519.Verify(pubKey, signedData, signatureBytes); !isValidSig {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "X-Signature is not valid",
			})
			return
		}

		c.Next()
	}
}

func ginServer(cmd *cobra.Command, args []string) {
	config.Load()
	if config.FARMA_VERSION != "" {
		gin.SetMode(gin.ReleaseMode)
	}
	err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	hub := fctools.NewFarcasterHub()
	defer hub.Close()

	serverAddr := config.GetString("host.addr")
	if a, _ := cmd.Flags().GetString("address"); a != "" {
		serverAddr = a
	}

	router := gin.Default()

	apiv1 := router.Group("/api/v1", verifySignature())
	{
		apiv1.GET("/frames/*id", api.H_FramesGet)
		apiv1.POST("/frames/", api.H_FrameAdd)

		apiv1.GET("/subscriptions/*frameId", api.H_SubscriptionsGet)
		apiv1.GET("/logs/*userId", api.H_LogsGet)

		apiv1.GET("/dbkeys/*prefix", api.H_DbKeysGet)

		apiv1.POST("/notifications/", api.H_Notify)
	}
	router.GET("/api/v1/version", api.H_Version)
	router.POST("/f/:id", api.WebhookHandler(hub))

	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}
	go func() {
		log.Printf("Starting farma %s\n", config.FARMA_VERSION)
		log.Printf("Listening and serving HTTP on %s", serverAddr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // block until signal received

	log.Println("Shutting down Farma")
	server.Shutdown(nil)
}
