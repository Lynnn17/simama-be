package auth

import (
	"time"

	"github.com/gofrs/uuid"
)

type LogSystem struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Actions    string    `json:"actions" db:"actions"`
	Jam        time.Time `json:"jam" db:"jam"`
	Keterangan string    `json:"keterangan" db:"keterangan"`
	IdUser     uuid.UUID `json:"idUser" db:"id_user"`
	Platform   string    `json:"platform" db:"platform"`
	IpAddress  string    `json:"ipAddress" db:"ip_address"`
	UserAgent  string    `json:"userAgent" db:"user_agent"`
	Kode       string    `json:"kode" db:"kode"`
}

type LogSystemDto struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Actions    string    `json:"actions" db:"actions"`
	Jam        time.Time `json:"jam" db:"jam"`
	Keterangan string    `json:"keterangan" db:"keterangan"`
	IdUser     uuid.UUID `json:"idUser" db:"id_user"`
	Platform   string    `json:"platform" db:"platform"`
	IpAddress  string    `json:"ipAddress" db:"ip_address"`
	UserAgent  string    `json:"userAgent" db:"user_agent"`
	Kode       string    `json:"kode" db:"kode"`
	NamaUser   *string   `json:"namaUser" db:"nama_user"`
	Username   *string   `json:"username" db:"username"`
	Email      *string   `json:"email" db:"email"`
	IdRole     *string   `json:"roleId" db:"role_id"`
	NamaRole   *string   `json:"namaRole" db:"nama_role"`
}

type RequestLogSystemFormat struct {
	Actions    string `json:"actions" db:"actions"`
	Keterangan string `json:"keterangan" db:"keterangan"`
	Kode       string `json:"kode" db:"kode"`
}

var ColumnLogSystemDto = map[string]interface{}{
	"id":         "id",
	"action":     "action",
	"jam":        "jam",
	"keterangan": "keterangan",
	"idUser":     "id_user",
	"platform":   "platform",
	"ipAddress":  "ip_address",
	"userAgent":  "user_agent",
	"kode":       "kode",
	"namaUser":   "nama_user",
	"username":   "username",
	"email":      "email",
	"roleId":     "role_id",
	"namaRole":   "nama_role",
}

func (logSystem *LogSystem) NewLogSystemFormat(reqFormat RequestLogSystemFormat, userId uuid.UUID, ipAddress string, userAgent string) (newLogSystem LogSystem, err error) {
	newID, _ := uuid.NewV4()

	newLogSystem = LogSystem{
		ID:         newID,
		Actions:    reqFormat.Actions,
		Jam:        time.Now(),
		Keterangan: reqFormat.Keterangan,
		IdUser:     userId,
		Platform:   "WEB",
		IpAddress:  ipAddress,
		UserAgent:  userAgent,
		Kode:       reqFormat.Kode,
	}

	return
}
