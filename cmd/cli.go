package cmd

import (
	"bufio"
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/vrypan/farma/config"
)

var cliCmd = &cobra.Command{
	Use:   "cli [json payload]",
	Short: "Sign and Send payload to the server API",
	Run:   cli,
}

func init() {
	rootCmd.AddCommand(cliCmd)
	cliCmd.Flags().StringP("url", "u", "config", "API endpoint. Defaults to host.addr/api/v1/ (from config file)")
	cliCmd.Flags().BoolP("send", "s", false, "Send generated JSON to server")
	cliCmd.Flags().BoolP("print", "p", true, "Print generated JSON request/response to stdout")
}

func cli(cmd *cobra.Command, args []string) {
	serverAddr, err := cmd.Flags().GetString("url")
	if err != nil || serverAddr == "config" {
		serverAddr = "http://" + config.GetString("host.addr") + "/api/v1/"
	}

	if a, _ := cmd.Flags().GetString("address"); a != "" {
		serverAddr = a
	}
	sendFlag, _ := cmd.Flags().GetBool("send")
	printFlag, _ := cmd.Flags().GetBool("print")

	config.Load()
	//keyPublic := config.GetString("key.public")
	keyPrivate := config.GetString("key.private")
	if keyPrivate == "" {
		fmt.Println("No private key. Use \n$ farma config set key.private <private_key>\nor the environment variable FARMA_KEY_PRIVATE")
		return
	}
	keyPrivateBytes, err := hex.DecodeString(keyPrivate[2:])
	if err != nil {
		fmt.Println("Error converting private key from hex:", err)
		return
	}

	keyPublic := hex.EncodeToString(
		ed25519.PrivateKey(keyPrivateBytes).Public().(ed25519.PublicKey),
	)
	if len(args) == 0 {
		fmt.Println("No payload")
		return
	}

	payload := args[0]
	if args[0] == "-" {
		var buffer bytes.Buffer
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			buffer.WriteString(scanner.Text())
		}
		payload = buffer.String()
	} else {
		fmt.Println("Error reading payload from stdin")
		os.Exit(1)
	}

	header := struct {
		Fid  int    `json:"fid"`
		Type string `json:"type"`
		Key  string `json:"key"`
	}{
		Fid:  0,
		Type: "shared",
		Key:  keyPublic,
	}
	headerBytes, err := json.Marshal(header)
	if err != nil {
		fmt.Println("Error marshaling header:", err)
		return
	}
	header64 := base64.RawURLEncoding.EncodeToString(headerBytes)
	payload64 := base64.RawURLEncoding.EncodeToString([]byte(payload))
	message := header64 + "." + payload64

	// This line is commented out because signerCombined and hash are not defined in this scope
	// You would need to find the correct values for these variables for the line to work
	signature64 := base64.RawURLEncoding.EncodeToString(ed25519.Sign(keyPrivateBytes, []byte(message)))
	response := struct {
		Header    string `json:"header"`
		Payload   string `json:"payload"`
		Signature string `json:"signature"`
	}{
		Header:    header64,
		Payload:   payload64,
		Signature: signature64,
	}
	requestJson, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling response:", err)
		return
	}
	//fmt.Println(string(ret))

	if printFlag {
		fmt.Println(string(requestJson))
	}

	if sendFlag {
		server := serverAddr
		httpClient := &http.Client{}            // creating an HTTP client
		reqBody := bytes.NewBuffer(requestJson) // converting ret into a buffer for use as the request body

		req, err := http.NewRequest("POST", server, reqBody) // creating a POST request
		if err != nil {
			fmt.Println("Error creating HTTP request:", err)
			return
		}

		req.Header.Set("Content-Type", "application/json") // setting the Content-Type header

		resp, err := httpClient.Do(req) // sending the HTTP request and receiving a response
		if err != nil {
			fmt.Println("Error sending request to server:", err)
			return
		}
		defer resp.Body.Close() // ensure that the response body will be closed

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Server returned non-200 status code:", resp.StatusCode)
			return
		}

		respBody, err := io.ReadAll(resp.Body) // reading the response body
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		if printFlag {
			fmt.Println(string(respBody))
		}
	}

}
