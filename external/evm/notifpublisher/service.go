package notifpublisher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"lms-be/shared/failure"

	"github.com/rs/zerolog/log"
)

// Client is the client utility for publishing notification
type NotifPublisher struct {
	Client    *http.Client
	BaseURL   string
	Endpoints Endpoints
}

type Endpoints struct {
	Whatsapp string
}

// SendWhatsapp send notif using whatsapp message
func (s *NotifPublisher) SendWhatsapp(request PublishRequest) (err error) {
	url := fmt.Sprintf("%s%s", s.BaseURL, s.Endpoints.Whatsapp)
	jsonBody, _ := json.Marshal(request)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return
	}

	// send request using HTTP client
	resp, err := s.Client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var response PublishResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Error().Err(err).Interface("Response", response).Msg("failed unmarshalling response from whatsapp notif publisher")
		return
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
			err = failure.BadRequestFromString(response.ErrorMessage)
			return
		}
		if resp.StatusCode == http.StatusUnauthorized {
			err = failure.Unauthorized("Not authorized")
			return
		}
		err = failure.InternalError(fmt.Errorf(response.ErrorMessage))
		return
	}

	return nil
}
