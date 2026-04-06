package pagination

import (
	"math"
)

// Response is a standard list data
type Response struct {
	Items []interface{} `json:"items"`
	Meta  Metadata      `json:"meta"`
}

type ResponseSummaryMonitoringPO struct {
	Items []interface{}               `json:"items"`
	Meta  MetadataSummaryMonitoringPO `json:"meta"`
}
type ResponseSummaryMonitoringDO struct {
	Items []interface{}               `json:"items"`
	Meta  MetadataSummaryMonitoringDO `json:"meta"`
}

// Metadata is a additional info for list data
type Metadata struct {
	TotalItems   int `json:"totalItems"`
	TotalPage    int `json:"totalPage"`
	PreviousPage int `json:"previousPage"`
	CurrentPage  int `json:"currentPage"`
	NextPage     int `json:"nextPage"`
	LimitPerPage int `json:"limitPerPage"`
}

type MetadataSummaryMonitoringPO struct {
	TotalItems          int     `json:"totalItems"`
	TotalPage           int     `json:"totalPage"`
	PreviousPage        int     `json:"previousPage"`
	CurrentPage         int     `json:"currentPage"`
	NextPage            int     `json:"nextPage"`
	LimitPerPage        int     `json:"limitPerPage"`
	SummaryOrderQTY     float64 `json:"summaryOrderQTY"`
	SummaryDeliveredQTY float64 `json:"summaryDeliveredQTY"`
	SummaryReceivedQTY  float64 `json:"summaryReceivedQTY"`
	SummaryRemainingQTY float64 `json:"summaryRemainingQTY"`
}

type MetadataSummaryMonitoringDO struct {
	TotalItems         int     `json:"totalItems"`
	TotalPage          int     `json:"totalPage"`
	PreviousPage       int     `json:"previousPage"`
	CurrentPage        int     `json:"currentPage"`
	NextPage           int     `json:"nextPage"`
	LimitPerPage       int     `json:"limitPerPage"`
	SummaryQTY         float64 `json:"summaryQTY"`
	SummaryQTYReceived float64 `json:"summaryQTYReceived"`
}

// CreateMeta is a metadata creator
func CreateMeta(totalItems int, dataPerPage int, pageNumber int) (meta Metadata) {
	totalPageRaw := float64(totalItems) / float64(dataPerPage)
	maxPage := int(math.Ceil(totalPageRaw))
	minPage := 1

	if maxPage < minPage {
		maxPage = minPage
	}

	nextPage := pageNumber + 1
	if nextPage > maxPage {
		nextPage = maxPage
	}

	prevPage := pageNumber - 1
	if prevPage < minPage {
		prevPage = minPage
	}

	return Metadata{
		TotalItems:   totalItems,
		TotalPage:    maxPage,
		PreviousPage: prevPage,
		CurrentPage:  pageNumber,
		NextPage:     nextPage,
		LimitPerPage: dataPerPage,
	}
}

// Meta Summary Monitoring PO
func CreateMetaSummaryMonitoringPO(totalItems int, dataPerPage int, pageNumber int, orderQty, deliveredQty, receivedQty, remainingQty float64) (meta MetadataSummaryMonitoringPO) {
	metaData := CreateMeta(totalItems, dataPerPage, pageNumber)
	return MetadataSummaryMonitoringPO{
		TotalItems:          metaData.TotalItems,
		TotalPage:           metaData.TotalPage,
		PreviousPage:        metaData.PreviousPage,
		CurrentPage:         metaData.CurrentPage,
		NextPage:            metaData.NextPage,
		LimitPerPage:        metaData.LimitPerPage,
		SummaryOrderQTY:     orderQty,
		SummaryDeliveredQTY: deliveredQty,
		SummaryReceivedQTY:  receivedQty,
		SummaryRemainingQTY: remainingQty,
	}
}

// Meta Summary Monitoring DO
func CreateMetaSummaryMonitoringDO(totalItems int, dataPerPage int, pageNumber int, quantity, quantityReceived float64) (meta MetadataSummaryMonitoringDO) {
	metaData := CreateMeta(totalItems, dataPerPage, pageNumber)
	return MetadataSummaryMonitoringDO{
		TotalItems:         metaData.TotalItems,
		TotalPage:          metaData.TotalPage,
		PreviousPage:       metaData.PreviousPage,
		CurrentPage:        metaData.CurrentPage,
		NextPage:           metaData.NextPage,
		LimitPerPage:       metaData.LimitPerPage,
		SummaryQTY:         quantity,
		SummaryQTYReceived: quantityReceived,
	}
}
