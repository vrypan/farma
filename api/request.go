package api

import (
	"bytes"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

type Request struct {
	Method     string
	Path       string
	Date       string
	Signature  string
	Body       []byte
	Query      string
	privateKey []byte
}

type ApiResult struct {
	Error  string            `json:"error"`
	Result []json.RawMessage `json:"result"`
	Next   string            `json:"next"`
}

func (r *Request) SignEd25519(privateKey []byte) *Request {
	r.Date = time.Now().UTC().Format(time.RFC1123)
	signature := ed25519.Sign(privateKey,
		[]byte(r.Method+"\n"+r.Path+"\n"+r.Date),
	)
	sig := hex.EncodeToString(signature)
	r.Signature = sig
	r.privateKey = privateKey
	return r
}

func (r *Request) VerifyEd25519(pubKey []byte) error {
	parsedTime, err := time.Parse(time.RFC1123, r.Date)
	if err != nil {
		return fmt.Errorf("Error parsing Date: %v", err)
	}
	now := time.Now().UTC()
	diffSeconds := int(math.Abs(float64(now.Sub(parsedTime).Seconds())))
	if diffSeconds > 10 {
		return fmt.Errorf("Date diff more than 10 seconds")
	}
	signatureBytes, err := hex.DecodeString(r.Signature)
	if err != nil {
		return fmt.Errorf("Error decoding signature: %v", err)
	}
	signedData := []byte(r.Method + "\n" + r.Path + "\n" + r.Date)
	if isValidSig := ed25519.Verify(pubKey, signedData, signatureBytes); !isValidSig {
		return fmt.Errorf("X-Signature is not valid")
	}
	return nil
}

func (r *Request) SignHmac(key []byte) *Request {
	mac := hmac.New(sha512.New, key)
	mac.Reset()

	r.Date = time.Now().UTC().Format(time.RFC1123)

	mac.Write([]byte(r.Method + "\n" + r.Path + "\n" + r.Date))

	signature := mac.Sum(nil)
	sig := hex.EncodeToString(signature)

	r.Signature = sig
	return r
}

func (r *Request) VerifyHmac(key []byte) error {
	mac := hmac.New(sha512.New, key)
	mac.Reset()

	parsedTime, err := time.Parse(time.RFC1123, r.Date)
	if err != nil {
		return fmt.Errorf("Error parsing Date: %v", err)
	}
	now := time.Now().UTC()
	diffSeconds := int(math.Abs(float64(now.Sub(parsedTime).Seconds())))
	if diffSeconds > 10 {
		return fmt.Errorf("Date diff more than 10 seconds")
	}

	mac.Write([]byte(r.Method + "\n" + r.Path + "\n" + r.Date))

	signature := mac.Sum(nil)
	sig := hex.EncodeToString(signature)

	if r.Signature != sig {
		return fmt.Errorf("Calculated signature does not match X-Signature")
	}
	return nil
}

func (r *Request) Send(server string) ([]byte, error) {
	if r.Method == "" || r.Path == "" || r.Date == "" || r.Signature == "" {
		return nil, fmt.Errorf("Request is not signed")
	}

	client := &http.Client{}
	req, err := http.NewRequest(r.Method, fmt.Sprintf("%s%s?%s", server, r.Path, r.Query), bytes.NewBuffer(r.Body))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Date", r.Date)
	req.Header.Set("X-Signature", r.Signature)

	if r.privateKey != nil {
		pubKeyBytes := ed25519.PublicKey(r.privateKey)
		pubKey := string(pubKeyBytes)
		req.Header.Set("X-Public_key", pubKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	var retErr error
	if (resp.StatusCode != http.StatusOK) && (resp.StatusCode != http.StatusCreated) {
		retErr = fmt.Errorf("Server returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}

	return body, retErr
}
