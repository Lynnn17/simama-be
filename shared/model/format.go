package model

import (
	"strings"
	"time"

	"github.com/bojanz/currency"
	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
)

const (
	// DefaultDateFormat digunakan untuk standar format input dari string ke format tanggal
	DateDMYFormat = "02/01/2006"

	// DefaultDateFormat digunakan untuk standar format input dari string ke format tanggal
	DefaultDateFormat = "2006-01-02"

	DefaultDateYearFormat = "2006-01"

	// DefaultDateFormat digunakan untuk standar format input dari string ke format tanggal dan jam
	DefaultDateTimeFormat = "2006-01-02 15:04:05"

	// DefaultTimeFormat digunakan untuk standar format input dari string ke format jam
	DefaultTimeFormat = "15:04:05"

	ISOFormat = "2006-01-02T15:04:05"
)

// DecimalToRupiah parse decimal to rupiah format
func DecimalToRupiah(value decimal.Decimal) string {
	ac := accounting.Accounting{Symbol: "Rp. ", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(value)
}

// RupiahToDecimal parse rupiah format to decimal
func RupiahToDecimal(rupiah string) (value decimal.Decimal, err error) {
	locale := currency.NewLocale("id")
	formatter := currency.NewFormatter(locale)
	amount, err := formatter.Parse(rupiah, "IDR")
	if err != nil {
		return
	}
	value, err = decimal.NewFromString(amount.Number())
	if err != nil {
		return
	}
	return
}

func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func ParseSplitString(val string) string {
	var ids string
	resSplit := strings.Split(val, ",")
	lengId := len(resSplit)
	for i, s := range resSplit {
		join := "'" + strings.TrimSpace(s) + "'"
		if i == lengId-1 {
			ids += join
		} else {
			ids += join + ","
		}
	}

	return ids
}

func StringJoin(arr []string) string {
	var ids string
	lengId := len(arr)
	for i, s := range arr {
		join := "'" + strings.TrimSpace(s) + "'"
		if i == lengId-1 {
			ids += join
		} else {
			ids += join + ","
		}
	}

	return ids
}

func SplitString(val string, demiliter string) []string {
	resSplit := strings.Split(val, demiliter)
	return resSplit
}

func ParseString(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

func ParseInt(s *int) int {
	if s != nil {
		return *s
	}

	return 0
}

func ParseFloat64(s *float64) float64 {
	if s != nil {
		return *s
	}

	return 0
}

func NullString(val string) *string {
	if val != "" {
		return &val
	}

	return nil
}

func Ternary(statement bool, a, b interface{}) interface{} {
	if statement {
		return a
	}
	return b
}

func ParseBool(s *bool) bool {
	if s != nil {
		return *s
	}

	return false
}

// JSONDate is a custom type for handling YYYY-MM-DD format in JSON
type JSONDate time.Time

func (j *JSONDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		// Try fallback to RFC3339
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}
	*j = JSONDate(t)
	return nil
}

func (j JSONDate) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(j).Format("2006-01-02") + "\""), nil
}

func (j JSONDate) Time() time.Time {
	return time.Time(j)
}
