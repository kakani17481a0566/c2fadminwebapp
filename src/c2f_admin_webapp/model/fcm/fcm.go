package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	// "google.golang.org/api/option"
	// "google.golang.org/api/transport"

	"c2f_admin_webapp/model/util/log"

	"github.com/astaxie/beego"
)

// Notification represents the notification fields
type Notification struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
	Image string `json:"image,omitempty"`
	Color string `json:"color,omitempty"`
}

// FCMMessage wraps the message to send via HTTP v1 API
type FCMMessage struct {
	Message struct {
		Token        string       `json:"token,omitempty"`  // single device token
		Tokens       []string     `json:"tokens,omitempty"` // batch send (optional; can send only one message per call for tokens array)
		Notification Notification `json:"notification"`
	} `json:"message"`
}

// SendNotificationToAll sends push notification using FCM HTTP v1 API
func SendNotificationToAll(tokens []string, title string, notificationBody string, image string) error {

	fmt.Println("tokens", tokens)
	fmt.Println("title", title)
	fmt.Println("notificationBody", notificationBody)
	fmt.Println("image", image)

	if len(tokens) == 0 {
		return fmt.Errorf("no tokens to send")
	}

	// Filter tokens: remove duplicates and empty
	tokenSet := make(map[string]struct{})
	cleanTokens := []string{}
	for _, token := range tokens {
		t := strings.TrimSpace(token)
		if t != "" {
			if _, exists := tokenSet[t]; !exists {
				tokenSet[t] = struct{}{}
				cleanTokens = append(cleanTokens, t)
			}
		}
	}

	if len(cleanTokens) == 0 {
		return fmt.Errorf("no valid tokens after filtering")
	}

	// Load project ID from config/env
	projectID := beego.AppConfig.String("FCM_PROJECT_ID")
	if projectID == "" {
		return fmt.Errorf("FCM_PROJECT_ID not configured")
	}

	// Load Service Account JSON file path from config/env
	serviceAccountPath := beego.AppConfig.String("FCM_SERVICE_ACCOUNT_FILE")
	if serviceAccountPath == "" {
		return fmt.Errorf("FCM_SERVICE_ACCOUNT_FILE not configured")
	}

	// Get OAuth2 token source from Service Account file
	ctx := context.Background()
	creds, err := google.CredentialsFromJSON(ctx, []byte(readFile(serviceAccountPath)), "https://www.googleapis.com/auth/firebase.messaging")
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error getting credentials", err)
		return err
	}

	// Compose request payload for each token (HTTP v1 doesn't support multiple tokens in one call; send multiple calls)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, token := range cleanTokens {

		payload := map[string]interface{}{
			"message": map[string]interface{}{
				"token": token,
				"notification": map[string]interface{}{
					"title": title,
					"body":  notificationBody,
					"image": image,
				},
				// Optional: you can add "android", "apns", or "webpush" keys here for platform-specific options
			},
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error marshaling JSON", err)
			return err
		}

		url := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/krishivala-397fc/messages:send")
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error creating HTTP request", err)
			return err
		}

		tokenSource, err := creds.TokenSource.Token()
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error fetching OAuth2 token", err)
			return err
		}

		req.Header.Set("Authorization", "Bearer "+tokenSource.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error sending FCM request", err)
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error reading response body", err)
			return err
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "FCM Response: ", string(body))

		if resp.StatusCode != http.StatusOK {
			var fcmResp map[string]interface{}
			if err := json.Unmarshal(body, &fcmResp); err == nil {
				if errObj, ok := fcmResp["error"].(map[string]interface{}); ok {
					// Log the error message (not hardcoding status)
					if msg, ok := errObj["message"].(string); ok {
						log.Println("Info", "Skipping token:", token, "Reason:", msg)
						continue
					}
				}
			}
			// If response is not ok but also not parseable, still continue
			log.Println("Info", "Skipping token due to unknown error:", token, "Response:", string(body))
			continue
		}

	}

	return nil
}

// Helper function to read file contents
func readFile(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error reading file", err)
		return []byte{}
	}
	return content
}
