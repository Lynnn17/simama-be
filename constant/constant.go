package constant

// Label Transaksi
const (
	PO         string = "PURCHASE_ORDER"
	DO         string = "DELIVERY_ORDER"
	SchedulePO string = "SCHEDULE_PO"
)

// PO Status
type POStatus string

const (
	PO_DRAFT     POStatus = "DRAFT"
	PO_OFFICIAL  POStatus = "OFFICIAL"
	PO_CONFIRM   POStatus = "CONFIRM"
	PO_OPEN      POStatus = "OPEN"
	PO_NEGOTIATE POStatus = "NEGOTIATE"
	PO_DONE      POStatus = "DONE"
	PO_REJECT    POStatus = "REJECT"
)

// Tipe untuk Schedule PO Status
type SchedulePOStatus string

const (
	SCHEDULE_PO_DRAFT      SchedulePOStatus = "DRAFT"
	SCHEDULE_PO_OFFICIAL   SchedulePOStatus = "OFFICIAL"
	SCHEDULE_PO_NEGOTIATE  SchedulePOStatus = "NEGOTIATE"
	SCHEDULE_PO_CONFIRM    SchedulePOStatus = "CONFIRM"
	SCHEDULE_PO_RESCHEDULE SchedulePOStatus = "RESCHEDULE"
	SCHEDULE_PO_DONE       SchedulePOStatus = "DONE"
	SCHEDULE_PO_CLOSE      SchedulePOStatus = "CLOSE"
	SCHEDULE_PO_REJECT     SchedulePOStatus = "REJECT"
)

// Tipe untuk Delivery Order Status
type DOStatus string

const (
	DO_ON_PREPARE DOStatus = "ON_PREPARE"
	DO_SHIPPED    DOStatus = "SHIPPED"
	DO_RECEIVED   DOStatus = "RECEIVED"
	DO_REJECT     DOStatus = "CANCEL"
	DO_DONE       DOStatus = "DONE"
)

var (
	NoticationStatus = struct {
		READ   string
		UNREAD string
	}{
		READ:   "READ",
		UNREAD: "UNREAD",
	}
)

var (
	NoticationDeliveryMethod = struct {
		APPS     string
		EMAIL    string
		WHATSAPP string
	}{
		APPS:     "APPS",
		EMAIL:    "EMAIL",
		WHATSAPP: "WHATSAPP",
	}
)

var (
	NoticationEntityType = struct {
		PO       string
		DO       string
		SCHEDULE string
	}{
		PO:       "PO",
		DO:       "DO",
		SCHEDULE: "SCHEDULE",
	}
)

var (
	NoticationEventType = struct {
		PO_APPROVAL       string
		DO_SHIPPED        string
		SCHEDULE_APPROVAL string
		RESCHEDULE        string
		INSPECTION_RESULT string
	}{
		PO_APPROVAL:       "PO_APPROVAL",
		DO_SHIPPED:        "DO_SHIPPED",
		SCHEDULE_APPROVAL: "SCHEDULE_APPROVAL",
		RESCHEDULE:        "RESCHEDULE",
		INSPECTION_RESULT: "INSPECTION_RESULT",
	}
)

var (
	NoticationPriority = struct {
		LOW      string
		MEDIUM   string
		HIGHT    string
		CRITICAL string
	}{
		LOW:      "LOW",
		MEDIUM:   "MEDIUM",
		HIGHT:    "HIGHT",
		CRITICAL: "CRITICAL",
	}
)
