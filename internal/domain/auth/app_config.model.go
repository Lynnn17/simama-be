package auth

import (
	"time"

	"github.com/gofrs/uuid"
)

type AppConfig struct {
	Id              int        `db:"id" json:"id"`
	AppName         *string    `db:"app_name" json:"appName"`
	AppLogo         *string    `db:"app_logo" json:"appLogo"`
	CompanyName     *string    `db:"company_name" json:"companyName"`
	CompanyEmail    *string    `db:"company_email" json:"companyEmail"`
	CompanyLogo     *string    `db:"company_logo" json:"companyLogo"`
	Address         *string    `db:"address" json:"address"`
	SmtpHost        *string    `db:"smtp_host" json:"smtpHost"`
	SmtpPort        *int       `db:"smtp_port" json:"smtpPort"`
	SmtpEmail       *string    `db:"smtp_email" json:"smtpEmail"`
	SmtpPassword    *string    `db:"smtp_password" json:"smtpPassword"`
	WorkhourPO      *float64   `db:"workhour_po" json:"workhourPO"`
	MonthlyWorkhour *float64   `db:"monthly_workhour" json:"monthlyWorkhour"`
	CreatedAt       *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy       *uuid.UUID `db:"updated_by" json:"updatedBy"`
}

type AppConfigDTO struct {
	Id              int      `db:"id" json:"id"`
	AppName         *string  `db:"app_name" json:"appName"`
	AppLogo         *string  `db:"app_logo" json:"appLogo"`
	CompanyName     *string  `db:"company_name" json:"companyName"`
	CompanyEmail    *string  `db:"company_email" json:"companyEmail"`
	CompanyLogo     *string  `db:"company_logo" json:"companyLogo"`
	Address         *string  `db:"address" json:"address"`
	SmtpHost        *string  `db:"smtp_host" json:"smtpHost"`
	SmtpPort        *int     `db:"smtp_port" json:"smtpPort"`
	SmtpEmail       *string  `db:"smtp_email" json:"smtpEmail"`
	SmtpPassword    *string  `db:"smtp_password" json:"smtpPassword"`
	WorkhourPO      *float64 `db:"workhour_po" json:"workhourPO"`
	MonthlyWorkhour *float64 `db:"monthly_workhour" json:"monthlyWorkhour"`
}

type AppConfigDTOByID struct {
	Id           int     `db:"id" json:"id"`
	AppName      *string `db:"app_name" json:"appName"`
	AppLogo      *string `db:"app_logo" json:"appLogo"`
	CompanyName  *string `db:"company_name" json:"companyName"`
	CompanyEmail *string `db:"company_email" json:"companyEmail"`
	CompanyLogo  *string `db:"company_logo" json:"companyLogo"`
	Address      *string `db:"address" json:"address"`
}

type RequestAppConfigFormat struct {
	AppName         string  `db:"app_name" json:"appName"`
	AppLogo         string  `db:"app_logo" json:"appLogo"`
	CompanyName     string  `db:"company_name" json:"companyName"`
	CompanyEmail    string  `db:"company_email" json:"companyEmail"`
	CompanyLogo     string  `db:"company_logo" json:"companyLogo"`
	Address         string  `db:"address" json:"address"`
	SmtpHost        string  `db:"smtp_host" json:"smtpHost"`
	SmtpPort        int     `db:"smtp_port" json:"smtpPort"`
	SmtpEmail       string  `db:"smtp_email" json:"smtpEmail"`
	SmtpPassword    string  `db:"smtp_password" json:"smtpPassword"`
	WorkhourPO      float64 `db:"workhour_po" json:"workhourPO"`
	MonthlyWorkhour float64 `db:"monthly_workhour" json:"monthlyWorkhour"`
}

var ColumnAppConfigDto = map[string]interface{}{
	"id":           "ac.id",
	"appName":      "ac.app_name",
	"appLogo":      "ac.app_logo",
	"companyName":  "ac.company_name",
	"companyEmail": "ac.company_email",
	"companyLogo":  "ac.company_logo",
	"address":      "ac.address",
	"smtpHost":     "ac.smtp_host",
	"smtpPort":     "ac.smtp_port",
	"smtpEmail":    "ac.smtp_email",
	"smtpPassword": "ac.smtp_password",
	"createdAt":    "ac.created_at",
	"updatedAt":    "ac.updated_at",
	"updatedBy":    "ac.updated_by",
}

func (m *AppConfig) AppConfigFormat(req RequestAppConfigFormat, userId uuid.UUID) {
	now := time.Now()

	m.AppName = &req.AppName
	m.AppLogo = &req.AppLogo
	m.CompanyName = &req.CompanyName
	m.CompanyEmail = &req.CompanyEmail
	m.CompanyLogo = &req.CompanyLogo
	m.Address = &req.Address
	m.SmtpHost = &req.SmtpHost
	m.SmtpPort = &req.SmtpPort
	m.SmtpEmail = &req.SmtpEmail
	if req.SmtpPassword != "__**password_not_changed**__" {
		m.SmtpPassword = &req.SmtpPassword
	}
	m.WorkhourPO = &req.WorkhourPO
	m.MonthlyWorkhour = &req.MonthlyWorkhour
	m.UpdatedAt = &now
	m.UpdatedBy = &userId
}
