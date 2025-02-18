package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/spf13/cobra"
	apiv1 "github.com/vrypan/farma/api"
	"github.com/vrypan/farma/config"
	"github.com/vrypan/farma/fctools"
	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/utils"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A server that receives Frame notification events and dumps them in stdout",
	Run:   server,
}

var hub *fctools.FarcasterHub
var verboseFlag bool

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringP("address", "a", "", "Listen on this address.")
	serverCmd.Flags().BoolP("verbose", "v", false, "Log additional info.")
}

func serverLog(request *http.Request, response int, other string) {
	log.Printf(
		"%s - \"%s %s %s\" \"%s\" %d %d \"%s\"",
		request.RemoteAddr,
		request.Method,
		request.URL.Path,
		request.Proto,
		request.Header["User-Agent"],
		request.ContentLength,
		response,
		other,
	)
}

// Allow only "/", alphanumeric and "-" characters in the path.
func isValidPath(path string) bool {
	matched, _ := regexp.MatchString(`^[\w/-]*$`, path)
	return matched
}

func server(cmd *cobra.Command, args []string) {
	serverAddr := config.GetString("host.addr")
	if a, _ := cmd.Flags().GetString("address"); a != "" {
		serverAddr = a
	}
	verboseFlag, _ = cmd.Flags().GetBool("verbose")

	err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	hub = fctools.NewFarcasterHub()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/", ApiH)
	mux.HandleFunc("/", notificationsH)

	fmt.Println("Starting server on", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		fmt.Println("Server error:", err)
	}
}
func ApiH(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	prefix := "/api/v1/"
	if r.URL.Path != prefix {
		log.Printf("Invalid API endpoint: %s", r.URL.Path)
		http.Error(w, "Invalid API endpoint", http.StatusNotFound)
		serverLog(r, http.StatusNotFound, "")
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		serverLog(r, http.StatusInternalServerError, "")
		return
	}
	api := apiv1.New()
	api.AddKey(config.GetString("key.public"))

	if err = api.Prepare(string(body)); err != nil {
		msg := fmt.Sprintf("Invalid API Call: %s", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	ret, err := api.Execute()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ret))
	serverLog(r, http.StatusOK, "")

}
func notificationsH(w http.ResponseWriter, r *http.Request) {
	// These are public endpoints that can and will be abused.
	// Let's make sure that HTTP requests are within some reasonable limits.
	if r.ContentLength > 1024 {
		serverLog(r, http.StatusBadRequest, "Content Length > 1024")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if len(r.URL.Path) > 128 {
		serverLog(r, http.StatusBadRequest, "Path Length > 128")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if !isValidPath(r.URL.Path) {
		serverLog(r, http.StatusBadRequest, "Path contains invalid_characters")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	frame := utils.NewFrame()
	if err := frame.FromEndpoint(r.URL.Path); err != nil {
		serverLog(r, http.StatusNotFound, err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		serverLog(r, http.StatusNoContent, err.Error())
		w.WriteHeader(http.StatusNoContent)
		return
	}

	subscription, eventType := utils.NewSubscription().FromHttpEvent(body)
	subscription.VerifyAppId(hub)
	subscription.FrameId = frame.Id
	if err = subscription.Save(); err != nil {
		log.Println("Error updating db.", err)
		log.Println("Subscription details:", subscription.NiceString())
		serverLog(r, http.StatusInternalServerError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ulog := utils.UserLog{
		FrameId:      subscription.FrameId,
		UserId:       subscription.UserId,
		AppId:        subscription.AppId,
		EvtType:      eventType,
		EventContext: "",
	}
	err = ulog.Save()

	serverLog(r, http.StatusOK, "")
	if verboseFlag {
		log.Println(subscription.NiceString())
	}
	w.WriteHeader(http.StatusOK)
}
