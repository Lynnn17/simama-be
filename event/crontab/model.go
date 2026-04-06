package crontab

type HistoryItem struct {
	DateTime        string  `json:"dateTime" db:"dateTime"`
	Description     string  `json:"description" db:"description"`
	TransactionCode string  `json:"transactionCode" db:"transaction_code"`
	Amount          float64 `json:"amount" db:"amount"`
	Flag            string  `json:"flag" db:"flag"`
	Ccy             string  `json:"ccy" db:"ccy"`
	ReffNo          string  `json:"reffno" db:"reffno"`
}

type ResponseBankJatim struct {
	ResponseCode string        `json:"responseCode"`
	ResponseDesc string        `json:"responseDesc"`
	History      []HistoryItem `json:"history"`
}
