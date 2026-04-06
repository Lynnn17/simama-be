package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"lms-be/infras"
	"lms-be/shared/logger"
)

type AppConfigRepository interface {
	ResolveByID(ctx context.Context, id int, isConfig bool) (data AppConfigDTO, err error)
	ResolveDTOByID(ctx context.Context, id int) (data AppConfigDTOByID, err error)
	GetByID(ctx context.Context, id int) (data AppConfig, err error)
	Update(ctx context.Context, req AppConfig) (err error)
}

type AppConfigRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideAppConfigRepositoryPostgreSQL(db *infras.PostgresqlConn) *AppConfigRepositoryPostgreSQL {
	s := new(AppConfigRepositoryPostgreSQL)
	s.DB = db
	return s
}

const selectDTOAppConfigQuery = `select ac.id, ac.app_name, ac.app_logo, ac.company_name, ac.company_email, ac.company_logo, ac.address, 
ac.smtp_host, ac.smtp_port, ac.smtp_email, '__**password_not_changed**__' as smtp_password, ac.workhour_po, ac.monthly_workhour
from app_config ac `
const selectDTOAppConfigQueryForSistem = `select ac.id, ac.app_name, ac.app_logo, ac.company_name, ac.company_email, ac.company_logo, ac.address, 
ac.smtp_host, ac.smtp_port, ac.smtp_email, ac.smtp_password, ac.workhour_po, ac.monthly_workhour 
from app_config ac `

func (r *AppConfigRepositoryPostgreSQL) ResolveByID(ctx context.Context, id int, isConfig bool) (data AppConfigDTO, err error) {
	queryText := selectDTOAppConfigQuery
	if isConfig {
		queryText = selectDTOAppConfigQueryForSistem
	}

	err = r.DB.Read.GetContext(ctx, &data, queryText+" WHERE ac.id=$1 ", id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New(fmt.Sprintf("Config dengan ID [%d] tidak ditemukan!", id))
		}
		logger.ErrorWithStack(err)
		return
	}

	return
}

const selectAppConfigQuery = `select ac.id, ac.app_name, ac.app_logo, ac.company_name, ac.company_email, ac.company_logo, ac.address, 
	ac.smtp_host, ac.smtp_port, ac.smtp_email, ac.smtp_password, ac.workhour_po, ac.monthly_workhour, ac.created_at, ac.updated_at, ac.updated_by 
	from app_config ac `

func (r *AppConfigRepositoryPostgreSQL) GetByID(ctx context.Context, id int) (data AppConfig, err error) {
	err = r.DB.Read.GetContext(ctx, &data, selectAppConfigQuery+" WHERE ac.id=$1 ", id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New(fmt.Sprintf("Config dengan ID [%d] tidak ditemukan!", id))
		}
		logger.ErrorWithStack(err)
		return
	}

	return
}

const updateAppConfigQuery = `update app_config set 
	id=:id, 
	app_name=:app_name, 
	app_logo=:app_logo, 
	company_name=:company_name, 
	company_email=:company_email, 
	company_logo=:company_logo, 
	address=:address, 
	smtp_host=:smtp_host, 
	smtp_port=:smtp_port, 
	smtp_email=:smtp_email, 
	smtp_password=:smtp_password, 
	workhour_po=:workhour_po,
	monthly_workhour=:monthly_workhour,
	updated_at=:updated_at, 
	updated_by=:updated_by 
`

func (r *AppConfigRepositoryPostgreSQL) Update(ctx context.Context, req AppConfig) (err error) {
	stmt, err := r.DB.Write.PrepareNamedContext(ctx, updateAppConfigQuery)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, req)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

const selectDTOByIDAppConfigQuery = `select ac.id, ac.app_name, ac.app_logo, ac.company_name, ac.company_email, ac.company_logo, ac.address 
from app_config ac `

func (r *AppConfigRepositoryPostgreSQL) ResolveDTOByID(ctx context.Context, id int) (data AppConfigDTOByID, err error) {
	err = r.DB.Read.GetContext(ctx, &data, selectDTOByIDAppConfigQuery+" WHERE ac.id=$1 ", id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New(fmt.Sprintf("Config dengan ID [%d] tidak ditemukan!", id))
		}
		logger.ErrorWithStack(err)
		return
	}

	return
}
