package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"lms-be/configs"
	"net/http"
	"path/filepath"
	"time"

	"golang.org/x/oauth2/google"
)

func getAccessToken(ctx context.Context, config *configs.Config) (string, string, error) {
	path := filepath.Join(config.App.File.Dir, config.App.Firebase.FCMCredensial)
	jsonKey, err := ioutil.ReadFile(path)
	if err != nil {
		return "", "", fmt.Errorf("error reading service account file: %v", err)
	}

	creds, err := google.CredentialsFromJSON(ctx, jsonKey, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return "", "", fmt.Errorf("error creating credentials: %v", err)
	}

	token, err := creds.TokenSource.Token()
	if err != nil {
		return "", "", fmt.Errorf("error getting token: %v", err)
	}

	return token.AccessToken, creds.ProjectID, nil
}

func SendNotification(config *configs.Config, title, body, token string) error {
	ctx := context.Background()
	accessToken, projectId, err := getAccessToken(ctx, config)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to get access token: %w", err)
	}

	msg := FCMMessage{
		Message: Message{
			Token: token,
			Notification: &FcmNotification{
				Title: title,
				Body:  body,
			},
			Data: map[string]string{
				"title": title,
				"body":  body,
			},
			Android: &Android{
				Notification: AndroidNotification{
					ClickAction: "FLUTTER_NOTIFICATION_CLICK",
				},
			},
			APNS: &APNS{
				Payload: APNSPayload{
					APS: APS{
						MutableContent:   1,
						ContentAvailable: 1,
					},
				},
			},
		},
	}

	messageJSON, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	fcmURL := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", projectId)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fcmURL, bytes.NewBuffer(messageJSON))
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyResp, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fcm send failed (%d): %s", resp.StatusCode, bodyResp)
	}

	return nil
}
