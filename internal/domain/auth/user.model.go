package auth

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"

	"lms-be/shared"
	"lms-be/shared/failure"
)

type User struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	Name           string     `json:"name" db:"name"`
	Username       string     `json:"username" db:"username"`
	Email          *string    `json:"email" db:"email"`
	Password       string     `json:"password" db:"password"`
	RoleId         string     `json:"roleId" db:"role_id"`
	Status         *string    `json:"status" db:"status"`
	Foto           *string    `json:"foto" db:"foto"`
	Active         bool       `db:"active" json:"active"`
	MobileFcmToken string     `db:"mobile_fcm_token" json:"mobileFcmToken"`
	WebFcmToken    string     `db:"web_fcm_token" json:"webFcmToken"`
	CreatedBy      *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedBy      *uuid.UUID `db:"updated_by" json:"updatedBy"`
	CreatedAt      *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deletedAt"`
	IsDeleted      bool       `db:"is_deleted" json:"isDeleted"`
}

// UserUpdateFormat
type UserUpdateFormat struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	Username       string    `json:"username" db:"username"`
	Email          *string   `json:"email" db:"email"`
	Password       string    `json:"password" db:"password"`
	RoleId         string    `json:"roleId" db:"roleId"`
	Status         *string   `json:"status" db:"status"`
	Foto           *string   `json:"foto" db:"foto"`
	MobileFcmToken string    `db:"mobile_fcm_token" json:"mobileFcmToken"`
	WebFcmToken    string    `db:"web_fcm_token" json:"webFcmToken"`
	Active         bool      `json:"active" db:"active"`
	UserID         uuid.UUID `json:"-"`
}

// UserUpdateFcmTokenFormat
type UserUpdateFcmTokenFormat struct {
	ID       uuid.UUID `json:"id"`
	Device   string    `json:"device" validate:"required,oneof=MOBILE WEB"`
	FcmToken string    `json:"fcmToken"`
}

// UserDTO digunakan untuk model join ke Role
type UserDTO struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Username string    `json:"username" db:"username"`
	Email    *string   `json:"email" db:"email"`
	Password string    `json:"password" db:"password"`
	RoleId   string    `json:"roleId" db:"role_id"`
	Role     *string   `json:"role" db:"role"`
	Status   *string   `json:"status" db:"status"`
	Foto     *string   `json:"foto" db:"foto"`
	Active   bool      `db:"active" json:"active"`
}

// UserDTO digunakan untuk model join ke Role
type UserRoleDTO struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Username      string    `json:"username" db:"username"`
	Email         string    `json:"email" db:"email"`
	Status        *string   `json:"status" db:"status"`
	FirebaseToken *string   `json:"firebaseToken" db:"firebase_token"`
	IsDeleted     bool      `json:"isDeleted" db:"is_deleted"`
	RoleID        *string   `json:"roleId" db:"role_id"`
	Role          *string   `json:"role" db:"name"`
}

type StatusLogin string

const (
	SuccessLogin StatusLogin = "success"
	FailedLogin  StatusLogin = "failed"
)

type LoginActivity struct {
	ID       uuid.UUID   `json:"id" db:"id"`
	Username string      `json:"username" db:"username"`
	Status   StatusLogin `json:"status" db:"status"`
	Jam      time.Time   `json:"jam" db:"jam"`
}

func NewCreateActivityLogin(username string, status StatusLogin) LoginActivity {
	loginActivityID, _ := uuid.NewV4()
	return LoginActivity{
		ID:       loginActivityID,
		Username: username,
		Status:   status,
		Jam:      time.Now(),
	}
}

// InputUser is struct as register json body
type InputUser struct {
	Name     string  `json:"name" db:"name"`
	Username string  `json:"username" db:"username"`
	Email    *string `json:"email" db:"email"`
	Password string  `json:"password" db:"password"`
	RoleId   string  `json:"roleId" db:"roleId"`
	Status   *string `json:"status" db:"status"`
	Active   bool    `db:"active" json:"active"`
}

// Validate digunakan untuk memvalidasi inputan user
func (i InputUser) Validate() error {
	v := shared.GetValidator()
	v.RegisterValidation("alphaspace", shared.AlphaSpace)
	v.RegisterValidation("alphanumspace", shared.AlphaNumSpace)

	return v.Struct(i)
}

// CreateUser is function to parse from user input to user struct
func (i InputUser) CreateUser() User {
	userID, _ := uuid.NewV4()
	now := time.Now()
	hash, _ := bcrypt.GenerateFromPassword([]byte(i.Password), bcrypt.DefaultCost)
	var user = User{
		ID:        userID,
		RoleId:    i.RoleId,
		Username:  i.Username,
		Name:      i.Name,
		Status:    i.Status,
		Password:  string(hash),
		Email:     i.Email,
		Active:    i.Active,
		CreatedAt: &now,
	}

	return user
}

// CreateUser is function to parse from user input to user struct
func (i InputUser) Registrasi() User {
	userID, _ := uuid.NewV4()

	hash, _ := bcrypt.GenerateFromPassword([]byte(i.Password), bcrypt.DefaultCost)
	return User{
		ID:       userID,
		RoleId:   i.RoleId,
		Username: i.Username,
		Password: string(hash),
		Active:   i.Active,
	}
}

type InputChangePassword struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// Update is function to transform into to User entity
func (i InputChangePassword) Update(user User) (User, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(i.OldPassword))
	if err != nil {
		return User{}, failure.Conflict("update password", "password", "old password does not match with the current password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(i.NewPassword))
	if err == nil {
		return User{}, failure.Conflict("update password", "password", "new password match with the current password")
	}

	newPassword, _ := bcrypt.GenerateFromPassword([]byte(i.NewPassword), bcrypt.DefaultCost)
	now := time.Now()
	return User{
		ID:        user.ID,
		RoleId:    user.RoleId,
		Username:  user.Username,
		Password:  string(newPassword),
		Email:     user.Email,
		Active:    user.Active,
		UpdatedAt: &now,
	}, nil
}

// ResetPasswdUpdate is function to transform into to User entity
func (i InputChangePassword) ResetPasswdUpdate(user User) (User, error) {
	newPassword, _ := bcrypt.GenerateFromPassword([]byte(i.NewPassword), bcrypt.DefaultCost)
	return User{
		ID:       user.ID,
		RoleId:   user.RoleId,
		Username: user.Username,
		Password: string(newPassword),
	}, nil
}

// Update is function to transform into to User entity
func (i UserUpdateFormat) Update(user UserUpdateFormat) (User, error) {
	// newPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	now := time.Now()
	return User{
		ID:        user.ID,
		RoleId:    user.RoleId,
		Username:  user.Username,
		Name:      user.Name,
		Status:    user.Status,
		Email:     user.Email,
		Active:    user.Active,
		UpdatedAt: &now,
	}, nil
}

func (u *User) UpdateUserFormat(id uuid.UUID, user UserUpdateFormat) {
	now := time.Now()
	u.ID = user.ID
	u.RoleId = user.RoleId
	u.Username = user.Username
	u.Name = user.Name
	u.Email = user.Email
	u.Status = user.Status
	u.UpdatedAt = &now
	u.UpdatedBy = &user.UserID

	// generate new password
	if user.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		u.Password = string(hash)
	}
}

// Update is function to transform into to User entity untuk update token fcm
func (u *User) UpdateFcmToken(user UserUpdateFcmTokenFormat) {
	u.ID = user.ID
	if user.Device == "MOBILE" {
		u.MobileFcmToken = user.FcmToken
	} else if user.Device == "WEB" {
		u.WebFcmToken = user.FcmToken
	}
}

// InputLogin is struct as login json body
type InputLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	// RoleID   string `json:"roleId"`
}

type InputLoginWeb struct {
	Username      string `json:"username" validate:"required"`
	Password      string `json:"password" validate:"required"`
	CaptchaID     string `json:"captchaId" validate:"required"`
	CaptchaAnswer string `json:"captchaAnswer" validate:"required"`
	// RoleID   string `json:"roleId"`
}

// Response is represent respond login
func (r *InputLogin) Response(user UserDTO, role Role, accessToken string) ResponseLogin {
	return ResponseLogin{
		Token: ResponseLoginToken{
			AccessToken: accessToken,
		},
		User: ResponseLoginUser{
			ID:       user.ID,
			RoleId:   user.RoleId,
			Username: user.Username,
			Name:     user.Name,
			Status:   user.Status,
			Email:    user.Email,
			Foto:     user.Foto,
			Role:     role,
		},
	}
}

// ResponseLogin is result processing from login process
type ResponseLogin struct {
	Token ResponseLoginToken `json:"token"`
	User  ResponseLoginUser  `json:"user"`
}

// ResponseLoginUser deliver result of user entity
type ResponseLoginUser struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name" db:"name"`
	Username      string    `json:"username" db:"username"`
	Email         *string   `json:"email" db:"email"`
	Status        *string   `json:"status" db:"status"`
	RoleId        string    `json:"roleId" db:"role_id"`
	FirebaseToken *string   `json:"firebaseToken"`
	Foto          *string   `json:"foto" db:"foto"`
	Role          Role      `json:"role"`
}

// ResponseLoginToken deliver result of user token
type ResponseLoginToken struct {
	AccessToken string
}

// NewUserLoginClaims digunakan untuk mengeset nilai dari JWT
func NewUserLoginClaims(user UserDTO, expiredIn int) jwt.MapClaims {
	claims := jwt.MapClaims{}
	claims["userId"] = user.ID
	claims["roleId"] = user.RoleId
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Duration(expiredIn) * time.Hour).Unix()

	return claims
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// SoftDelete untuk mengeset flag isDeleted
func (u *User) SoftDelete(userID uuid.UUID) {
	now := time.Now()
	u.Active = false
	u.IsDeleted = true
	u.DeletedAt = &now
}

func (u *User) SoftActive(userID uuid.UUID, active bool) {
	now := time.Now()
	u.Active = active
	u.IsDeleted = false
	u.UpdatedAt = &now
	u.UpdatedBy = &userID
}

type ModelUpdateFoto struct {
	Id        uuid.UUID  `db:"id" json:"id"`
	Foto      string     `json:"foto" db:"foto"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *uuid.UUID `db:"updated_by" json:"updatedBy"`
}

type UpdateFotoRequest struct {
	Id       uuid.UUID `db:"id" json:"id"`
	Foto     string    `json:"file" db:"foto"`
	FotoLama string    `db:"foto_lama" json:"fotoLama"`
}

var ColumnMappUser = map[string]interface{}{
	"id":         "u.id",
	"name":       "u.name",
	"username":   "u.username",
	"email":      "u.email",
	"password":   "u.password",
	"roleId":     "u.role_id",
	"role":       "r.name",
	"personName": "p.name",
	"createdAt":  "u.created_at",
	"updatedAt":  "u.updated_at",
	"createdBy":  "u.created_by",
	"updatedBy":  "u.updated_by",
}
