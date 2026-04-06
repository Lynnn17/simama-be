package auth

import (
	"bytes"
	"fmt"

	"github.com/jmoiron/sqlx"

	"lms-be/infras"
	"lms-be/shared/logger"
	"lms-be/shared/model"
	"lms-be/shared/pagination"
)

var (
	roleQuery = struct {
		Select                        string
		Count                         string
		Insert                        string
		Update                        string
		Delete                        string
		Exist                         string
		DeleteMenuByRole              string
		InsertBulkRoleMenu            string
		InsertBulkRoleMenuPlaceHolder string
	}{
		Select: `select id, name, description from auth_role `,
		Count:  `select count(id) from auth_role `,
		Insert: `insert into auth_role(id, name, description)
							values(:id, :name, :description) `,
		Update: `Update auth_role set
						name=:name,
						description=:description
						where id=:id `,
		Delete: `delete from auth_role `,
		Exist:  `select count(id)>0 from auth_role `,
	}
)

type RoleRepository interface {
	ResolveAll(req model.StandardRequest) (response pagination.Response, err error)
	GetAllData() (roles []Role, err error)
	ExistRoleByID(id string) (bool, error)
	ExistRoleByName(nama string) (bool, error)
	CreateRole(role Role) error
	ResolveRoleByID(id string) (Role, error)
	UpdateRole(role Role) error
	DeleteRoleByID(id string) (err error)
}

type RoleRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideRoleRepositoryPostgreSQL(db *infras.PostgresqlConn) *RoleRepositoryPostgreSQL {
	s := new(RoleRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *RoleRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (response pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" WHERE ")
		searchRoleBuff.WriteString("  concat (name, description) like ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := roleQuery.Count + searchRoleBuff.String()
	queryDTO := roleQuery.Select + searchRoleBuff.String()

	query = r.DB.Read.Rebind(query)
	var totalData int
	err = r.DB.Read.QueryRow(query, searchParams...).Scan(&totalData)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if totalData < 1 {
		response.Items = make([]interface{}, 0)
		return
	}

	offset := (req.PageNumber - 1) * req.PageSize
	queryDTO = r.DB.Read.Rebind(queryDTO + fmt.Sprintf("order by %s %s  limit %d offset %d", req.SortBy, req.SortType, req.PageSize, offset))

	rows, err := r.DB.Read.Queryx(queryDTO, searchParams...)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	for rows.Next() {
		var role Role
		err = rows.StructScan(&role)
		if err != nil {
			return
		}

		response.Items = append(response.Items, role)
	}

	response.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)

	return
}

func (r *RoleRepositoryPostgreSQL) GetAllData() (roles []Role, err error) {
	err = r.DB.Read.Select(&roles, roleQuery.Select+" order by id")
	return
}

func (r *RoleRepositoryPostgreSQL) ExistRoleByID(id string) (bool, error) {
	var exist bool

	err := r.DB.Read.Get(&exist, roleQuery.Exist+" where id = $1", id)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *RoleRepositoryPostgreSQL) ExistRoleByName(nama string) (bool, error) {
	var exist bool

	err := r.DB.Read.Get(&exist, roleQuery.Exist+" where name = $1", nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

// CreateRole adalah method untuk menambahkan Role Baru
func (r *RoleRepositoryPostgreSQL) CreateRole(role Role) error {
	stmt, err := r.DB.Read.PrepareNamed(roleQuery.Insert)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(role)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}

// ResolveByID adalah method yang digunakan untuk mendapatkan data Role berdasarkan ID
func (r *RoleRepositoryPostgreSQL) ResolveRoleByID(id string) (Role, error) {
	var role Role
	err := r.DB.Read.Get(&role, roleQuery.Select+" where id=$1", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return Role{}, err
	}
	return role, nil
}

// UpdateRole adalah method untuk mengubah Role yang sudah ada
func (r *RoleRepositoryPostgreSQL) UpdateRole(role Role) error {
	stmt, err := r.DB.Write.PrepareNamed(roleQuery.Update)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(role)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}

// DeleteRoleByID adalah method untuk menghapus Role berdasarkan ID yang dikirimkan
func (r *RoleRepositoryPostgreSQL) DeleteRoleByID(id string) (err error) {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		stmt, err := r.DB.Read.PrepareNamed(roleQuery.Delete + " where id=:id ")
		if err != nil {
			logger.ErrorWithStack(err)
			e <- err
			return
		}
		var params = map[string]interface{}{
			"id": id,
		}
		defer stmt.Close()
		res, err := stmt.Exec(params)
		_, err = res.RowsAffected()
		if err != nil {
			logger.ErrorWithStack(err)
			e <- err
			return
		}
		e <- nil
	})
}
