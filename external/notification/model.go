package notification

// FCMMessage is the main structure for FCM request
type FCMMessage struct {
	Message Message `json:"message"`
}

type Message struct {
	Token        string            `json:"token,omitempty"`
	Notification *FcmNotification  `json:"notification,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
	Android      *Android          `json:"android,omitempty"`
	APNS         *APNS             `json:"apns,omitempty"`
}

type FcmNotification struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

type Android struct {
	Notification AndroidNotification `json:"notification,omitempty"`
}

type AndroidNotification struct {
	ClickAction string `json:"click_action,omitempty"`
	Sound       string `json:"sound,omitempty"`
}

type APNS struct {
	Payload APNSPayload `json:"payload,omitempty"`
}

type APNSPayload struct {
	APS APS `json:"aps,omitempty"`
}

type APS struct {
	MutableContent   int    `json:"mutable-content,omitempty"`
	ContentAvailable int    `json:"content-available,omitempty"`
	Sound            string `json:"sound,omitempty"`
}
