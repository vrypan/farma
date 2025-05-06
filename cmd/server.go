package cmd

import (
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	apiv2 "github.com/vrypan/farma/apiv2"
	"github.com/vrypan/farma/config"
	"github.com/vrypan/farma/fctools"

	db "github.com/vrypan/farma/localdb"
)

var StaticFiles embed.FS

var ginServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the farma server.",
	Run:   ginServer,
}

func init() {
	rootCmd.AddCommand(ginServerCmd)
	ginServerCmd.Flags().StringP("address", "a", "", "Listen on this address/port.")
	ginServerCmd.Flags().BoolP("verbose", "v", false, "Log additional info.")
	ginServerCmd.Flags().StringP("test-frame", "t", "", "Path to a directory with a static test frame")
}

func ginServer(cmd *cobra.Command, args []string) {
	testFrame, _ := cmd.Flags().GetString("test-frame")

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

	allowedHosts := config.GetStringSlice("host.cors")
	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedHosts,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Signature", "X-Public-Key", "X-Date"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	frameOrAdminGroup := router.Group("/api/v2", apiv2.VerifySignature(apiv2.ACL_FRAME_OR_ADMIN))
	{
		frameOrAdminGroup.GET("/frame/:frameId", apiv2.H_FramesGet)
		frameOrAdminGroup.POST("/frame/:frameId", apiv2.H_FrameUpdate)
		frameOrAdminGroup.GET("/subscription/*frameId", apiv2.H_SubscriptionsGet)
		frameOrAdminGroup.POST("/subscription-import/:frameId", apiv2.H_SubscriptionsImportCSV)
		frameOrAdminGroup.GET("/logs/:frameId/*userId", apiv2.H_LogsGet)
		frameOrAdminGroup.GET("/notification/:frameId", apiv2.H_NotificationsGet)
		frameOrAdminGroup.GET("/notification/:frameId/:notificationId", apiv2.H_NotificationsGet)
		frameOrAdminGroup.POST("/notification/:frameId", apiv2.H_Notify)
	}
	onlyAdminGroup := router.Group("/api/v2", apiv2.VerifySignature(apiv2.ACL_ADMIN))
	{
		onlyAdminGroup.GET("/frame/", apiv2.H_FramesGetAll)
		onlyAdminGroup.POST("/frame/", apiv2.VerifySignature(apiv2.ACL_ADMIN), apiv2.H_FrameAdd)
		onlyAdminGroup.GET("/dbkeys/*prefix", apiv2.VerifySignature(apiv2.ACL_ADMIN), apiv2.H_DbKeysGet)

	}
	router.GET("/api/v2/version", apiv2.H_Version)
	router.GET("/api/v2/new_keypair/:frameId", apiv2.H_NewKeypair)
	router.POST("/f/:id", apiv2.WebhookHandler(hub))

	if testFrame != "" {
		router.Static("/test", testFrame)
		router.Static("/.well-known", testFrame+"/.well-known")
	}

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
