package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/fctools"
	db "github.com/vrypan/farma/localdb"
	"github.com/vrypan/farma/utils"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A server that receives Frame notification events and dumps them in stdout",
	Run:   server,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntP("port", "", 8080, "HTTP port")
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
	err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	hub := fctools.NewFarcasterHub()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

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

		update := utils.NewNotificationUpdate(frame.Id).FromHttpEvent(body).GetAppFid(hub)
		log.Println(update)

		if err = update.UpdateDb(); err != nil {
			log.Println("Error updating db.", err)
			serverLog(r, http.StatusInternalServerError, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		serverLog(r, http.StatusOK, "")
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
