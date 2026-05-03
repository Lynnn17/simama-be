package auth

import (
	"bytes"
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"

	"lms-be/infras"
	"lms-be/shared/failure"
	"lms-be/shared/logger"
	"lms-be/shared/model"
	"lms-be/shared/pagination"
)

var (
	userQuery = struct {
		Insert,
		Exist,
		Select,
		SelectDTO,
		SelectVerifikasi,
		Count,
		Update,
		UpdatePassword,
		resetImei,
		UpdateFoto string
	}{
		Insert: `INSERT INTO auth_user (id, name, username, email, password, role_id, status, active, created_at, created_by) 
			VALUES (:id, :name, :username, :email, :password, :role_id, :status, :active, :created_at, :created_by) `,
		Exist: `SELECT COUNT(u.id) > 0 FROM auth_user u`,
		Select: `SELECT u.id, u.name, u.username, u.email, u.password, u.status, u.role_id, u.foto, u.active, u.mobile_fcm_token, u.web_fcm_token, u.created_by, u.updated_by, u.created_at, u.updated_at, u.deleted_at, u.is_deleted 
			FROM auth_user u `,
		SelectDTO: `SELECT u.id, u.name, u.username, u.email, u.password, u.status, r.name as role, r.id as role_id, u.active 
			FROM auth_user u
			left join auth_role r on r.id = u.role_id
			 `,
		Count: `select count(u.id) from auth_user u 
			left join auth_role r on r.id = u.role_id `,
		Update: `UPDATE auth_user SET 
		  	id=:id,
			name=:name,
			username=:username,
			email=:email,
			password=:password, 
			status=:status, 
			role_id=:role_id,
			mobile_fcm_token=:mobile_fcm_token,
			web_fcm_token=:web_fcm_token,
			is_deleted=:is_deleted,
			active=:active,
			deleted_at=:deleted_at,
			updated_by=:updated_by,
			updated_at=:updated_at `,
		UpdatePassword: `UPDATE auth_user SET
			password=:password,
			updated_at=:updated_at `,
		UpdateFoto: `UPDATE auth_user SET 
			id=:id, 
			foto=:foto,
			updated_at=:updated_at `,
	}

	loginActivityQuery = struct {
		Insert string
	}{
		Insert: `INSERT INTO log_activity (
			id,
			username,
			jam
		) VALUES (
			:id,
			:username,
			:jam
		)`,
	}
)

// UserRepositoryPostgreSQL digunakan untuk Repository User
type UserRepositoryPostgreSQL struct {
	DB             *infras.PostgresqlConn
	roleRepisitory RoleRepository
}

// ProvideUserRepositoryPostgreSQL is the provider for this repository.
func ProvideUserRepositoryPostgreSQL(db *infras.PostgresqlConn, rr RoleRepository) *UserRepositoryPostgreSQL {
	return &UserRepositoryPostgreSQL{
		DB:             db,
		roleRepisitory: rr,
	}
}

type UserRepository interface {
	GetAll(ctx context.Context, req model.StandardRequestUser) (data []User, err error)
	ResolveAll(ctx context.Context, req model.StandardRequestUser) (data pagination.Response, err error)
	CreateLoginActivity(loginActivity LoginActivity) error
	ExistByUsername(username string) (exist bool, err error)
	ExistByEmail(email string) (exist bool, err error)
	ExistByUsernameOrEmail(identifier string) (exist bool, err error)
	ResolveUserByUsername(username string) ([]User, error)
	ResolveUserByUsernameRole(username string) (UserDTO, error)
	ResolveUserByEmailRole(email string) (UserDTO, error)
	ResolveUserByUsernameOrEmailRole(identifier string) (UserDTO, error)
	ResolveUserByID(id uuid.UUID) (User, error)
	ResolveUserByRole(roleName string, idBidang string) (data []User, err error)
	ResolveUserByIDDTO(id uuid.UUID) (UserDTO, error)
	ResolveUserByName(name string) (UserDTO, error)
	TransactionCreateUser(user User) error
	TransactionUpdateUser(user User) error
	UpdateUser(id uuid.UUID, user User) error
	UpdateUserPassword(id uuid.UUID, user User) error
	UpdateFoto(data ModelUpdateFoto) error
}

// TransactionCreateUser digunakan untuk menambahkan user baru dalam blok transaction
func (u *UserRepositoryPostgreSQL) TransactionCreateUser(user User) error {
	return u.DB.WithTransaction(func(db *sqlx.Tx, errs chan error) {
		err := u.createUser(user)
		if err != nil {
			errs <- err
			return
		}

		errs <- nil
	})
}

func (u *UserRepositoryPostgreSQL) TransactionUpdateUser(user User) error {
	return u.DB.WithTransaction(func(db *sqlx.Tx, errs chan error) {
		err := u.UpdateUser(user.ID, user)
		if err != nil {
			errs <- err
			return
		}

		errs <- nil
	})
}

// createUser is method to create a new user
func (u *UserRepositoryPostgreSQL) createUser(user User) error {
	stmt, err := u.DB.Read.PrepareNamed(userQuery.Insert)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	_, err = stmt.Exec(user)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	return nil
}

// ExistByUsername is function to check that username exist or not
func (u *UserRepositoryPostgreSQL) ExistByUsername(username string) (exist bool, err error) {
	err = u.DB.Read.Get(&exist, userQuery.Exist+" WHERE username = $1 AND u.active is true AND u.is_deleted is false ", username)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return exist, err
}

// ExistByEmail is function to check that email exist or not
func (u *UserRepositoryPostgreSQL) ExistByEmail(email string) (exist bool, err error) {
	err = u.DB.Read.Get(&exist, userQuery.Exist+" WHERE email = $1 AND u.active is true AND u.is_deleted is false ", email)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return exist, err
}

// ExistByUsernameOrEmail is function to check that username or email exist or not
func (u *UserRepositoryPostgreSQL) ExistByUsernameOrEmail(identifier string) (exist bool, err error) {
	err = u.DB.Read.Get(&exist, userQuery.Exist+" WHERE (username = $1 OR email = $1) AND u.active is true AND u.is_deleted is false ", identifier)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return exist, err
}

// ResolveUserByUsername is function resolving user data by username
func (u *UserRepositoryPostgreSQL) ResolveUserByUsername(username string) (user []User, err error) {
	err = u.DB.Read.Select(&user, userQuery.Select+" WHERE u.username = $1 AND u.deleted_at is null", username)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return user, nil
}

// ResolveUserByUsernameRole is function resolving user data by username and role id
func (u *UserRepositoryPostgreSQL) ResolveUserByUsernameRole(username string) (UserDTO, error) {
	var user UserDTO
	err := u.DB.Read.Get(&user, userQuery.SelectDTO+" WHERE u.username = $1 and u.active = true and u.is_deleted = false  ", username)
	if err != nil {
		logger.ErrorWithStack(err)
		return UserDTO{}, err
	}

	return user, nil
}

// ResolveUserByEmailRole is function resolving user data by email and role id
func (u *UserRepositoryPostgreSQL) ResolveUserByEmailRole(email string) (UserDTO, error) {
	var user UserDTO
	err := u.DB.Read.Get(&user, userQuery.SelectDTO+" WHERE u.email = $1 and u.active = true and u.is_deleted = false  ", email)
	if err != nil {
		logger.ErrorWithStack(err)
		return UserDTO{}, err
	}

	return user, nil
}

// ResolveUserByUsernameOrEmailRole is function resolving user data by username or email and role id
func (u *UserRepositoryPostgreSQL) ResolveUserByUsernameOrEmailRole(identifier string) (UserDTO, error) {
	var user UserDTO
	err := u.DB.Read.Get(&user, userQuery.SelectDTO+" WHERE (u.username = $1 OR u.email = $1) and u.active = true and u.is_deleted = false  ", identifier)
	if err != nil {
		logger.ErrorWithStack(err)
		return UserDTO{}, err
	}

	return user, nil
}

// ResolveUserByID is function resolving user data by id
func (u *UserRepositoryPostgreSQL) ResolveUserByID(id uuid.UUID) (User, error) {
	var user User
	err := u.DB.Read.Get(&user, userQuery.Select+" WHERE u.id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, err
		}
		logger.ErrorWithStack(err)
		return User{}, err
	}
	return user, nil
}

// ResolveUserByRole is function resolving user data by role
func (u *UserRepositoryPostgreSQL) ResolveUserByRole(roleName string, idBidang string) (data []User, errr error) {
	rows, err := u.DB.Read.Queryx(userQuery.Select+" WHERE u.role_id = $1  AND u.active = true", roleName)
	if err == sql.ErrNoRows {
		errr = failure.NotFound("User")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var master User
		err = rows.StructScan(&master)
		if err != nil {
			return
		}
		data = append(data, master)
	}
	return
}

// ResolveUserByID is function resolving user data by email
func (u *UserRepositoryPostgreSQL) ResolveUserByIDDTO(id uuid.UUID) (UserDTO, error) {
	var user UserDTO
	err := u.DB.Read.Get(&user, userQuery.SelectDTO+" WHERE u.id = $1 AND u.deleted_at is null", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return UserDTO{}, err
		}
		logger.ErrorWithStack(err)
		return UserDTO{}, err
	}

	return user, nil
}

func (u *UserRepositoryPostgreSQL) ResolveUserByName(name string) (UserDTO, error) {
	var user UserDTO
	err := u.DB.Read.Get(&user, userQuery.SelectDTO+" WHERE u.name = $1 AND u.deleted_at is null", name)
	if err != nil {
		if err == sql.ErrNoRows {
			return UserDTO{}, err
		}
		logger.ErrorWithStack(err)
		return UserDTO{}, err
	}

	return user, nil
}

// UpdateUser is function to update the user entity
func (u *UserRepositoryPostgreSQL) UpdateUser(id uuid.UUID, user User) error {
	stmt, err := u.DB.Read.PrepareNamed(userQuery.Update + " WHERE id = :id")
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	_, err = stmt.Exec(user)
	defer stmt.Close()
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	return nil
}

// UpdateUserPassword is function to update the user password
func (u *UserRepositoryPostgreSQL) UpdateUserPassword(id uuid.UUID, user User) error {
	stmt, err := u.DB.Read.PrepareNamed(userQuery.UpdatePassword + " WHERE id = :id")
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	_, err = stmt.Exec(user)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	return nil
}

// CreateLoginActivity is function to create log from login activity
func (u *UserRepositoryPostgreSQL) CreateLoginActivity(loginActivity LoginActivity) error {
	stmt, err := u.DB.Read.PrepareNamed(loginActivityQuery.Insert)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	_, err = stmt.Exec(loginActivity)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	return nil
}

// ResolveAll digunakan untuk menampilkan semua data
func (r *UserRepositoryPostgreSQL) ResolveAll(ctx context.Context, req model.StandardRequestUser) (response pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(u.is_deleted, false)=false ")

	// if !req.Active {
	// searchRoleBuff.WriteString(" AND u.active = ? ")
	// searchParams = append(searchParams, req.Active)
	// }

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat (u.name, u.username, u.email, r.name, p.name) ilike ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	if req.RoleID != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" u.role_id = ?  ")
		searchParams = append(searchParams, req.RoleID)
	}

	query := r.DB.Read.Rebind("select count(*) from (" + userQuery.SelectDTO + searchRoleBuff.String() + ")s")
	var totalData int
	err = r.DB.Read.QueryRowxContext(ctx, query, searchParams...).Scan(&totalData)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if totalData < 1 {
		response.Items = make([]interface{}, 0)
		return
	}

	searchRoleBuff.WriteString("order by " + ColumnMappUser[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchUserQuery := searchRoleBuff.String()
	searchUserQuery = r.DB.Read.Rebind(userQuery.SelectDTO + searchUserQuery)
	rows, err := r.DB.Read.Queryx(searchUserQuery, searchParams...)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	for rows.Next() {
		var userDTO UserDTO
		err = rows.StructScan(&userDTO)
		if err != nil {
			return
		}

		response.Items = append(response.Items, userDTO)
	}

	response.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)

	return
}

func (r *UserRepositoryPostgreSQL) GetAll(ctx context.Context, req model.StandardRequestUser) (data []User, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(u.is_deleted, false) = false ")

	if req.RoleID != "" {
		searchRoleBuff.WriteString(" AND u.role_id = ? ")
		searchParams = append(searchParams, req.RoleID)
	}

	searchRoleBuff.WriteString(" ORDER BY name ASC ")
	query := r.DB.Read.Rebind(userQuery.Select + searchRoleBuff.String())
	rows, err := r.DB.Read.QueryxContext(ctx, query, searchParams...)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("User Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	for rows.Next() {
		var user User
		err = rows.StructScan(&user)

		if err != nil {
			return
		}
		data = append(data, user)
	}
	return
}

func (r *UserRepositoryPostgreSQL) UpdateFoto(data ModelUpdateFoto) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := txUpdateFoto(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func txUpdateFoto(tx *sqlx.Tx, data ModelUpdateFoto) (err error) {
	stmt, err := tx.PrepareNamed(userQuery.UpdateFoto + " where id=:id")
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(data)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}
