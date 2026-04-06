package notifpublisher

// PublishRequest is an entity for generic WhatsApp message publish request.
// TODO: Add media template fields.
// TODO: HSM is about to be deprecated
//       https://developers.facebook.com/docs/whatsapp/api/messages/message-templates#using-the-hsm-message-template-type
type PublishRequest struct {
	Source        string  `json:"source,omitempty"`
	Event         string  `json:"event,omitempty"`
	TransactionID string  `json:"transactionId,omitempty"`
	Identity      string   `json:"identity"`
	HSMParams     []string `json:"hsmParams"`
	Language      *string  `json:"language,omitempty"`
	Template      string   `json:"template"`
	To            string   `json:"to"`
	TTL           *int     `json:"ttl,omitempty"`
	Type          string   `json:"type"`
}

// PublishResponse represents the response to a publish request.
type PublishResponse struct {
	ErrorMessage string `json:"errorMessage,omitempty"`
	Success      bool    `json:"status"`
}
