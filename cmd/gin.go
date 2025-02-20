package cmd

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/utils"
	"google.golang.org/protobuf/proto"
)

var ginServerCmd = &cobra.Command{
	Use:   "gin",
	Short: "",
	Run:   ginServer,
}

func init() {
	rootCmd.AddCommand(ginServerCmd)
	ginServerCmd.Flags().StringP("address", "a", "", "Listen on this address/port.")
	ginServerCmd.Flags().BoolP("verbose", "v", false, "Log additional info.")
}

func ginServer(cmd *cobra.Command, args []string) {
	err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	router := gin.Default()
	router.GET("/api/v1/frames/", H_FramesGetAll)
	router.GET("/api/v1/frames/:id", H_FramesGet)
	router.POST("/api/v1/frames", H_FrameAdd)

	router.GET("/api/v1/subscriptions/", H_SubscriptionsGet)
	router.GET("/api/v1/subscriptions/:frameId", H_SubscriptionsGet)

	router.GET("/api/v1/logs/", H_LogsGet)
	router.GET("/api/v1/logs/:userId", H_LogsGet)

	router.GET("/api/v1/notifications/", apiHandler)

	router.Run(":1234")
}

func apiHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"path": c.Request.URL.Path})
}

func H_FramesGetAll(c *gin.Context) {
	frames := utils.AllFrames()
	c.JSON(http.StatusOK, gin.H{"frames": frames})
}

func H_FramesGet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_ID"})
		return
	}

	frame := utils.NewFrame().FromId(id)
	if frame == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FRAME_NOT_FOUND"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"frame": frame})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "frame exists"})
	}
	if err != db.ERR_NOT_FOUND {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	}
	c.JSON(http.StatusCreated, gin.H{"frame": frame})
}

func H_SubscriptionsGet(c *gin.Context) {
	frameId := c.Param("frameId")
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
	c.JSON(http.StatusOK, gin.H{"subscriptions": list})
}

func H_LogsGet(c *gin.Context) {
	userId := c.Param("userId")
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
	list := make([]*utils.Subscription, len(data))
	for i, item := range data {
		var pb utils.Subscription
		proto.Unmarshal(item, &pb)
		list[i] = &pb
	}
	c.JSON(http.StatusOK, gin.H{"result": list})
}
