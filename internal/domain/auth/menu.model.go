package auth

import (
	"time"

	"github.com/gofrs/uuid"
)

type Menu struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	Name            string     `db:"name" json:"name"`
	Link            string     `db:"link" json:"link"`
	Icon            *string    `db:"icon" json:"icon"`
	Description     *string    `db:"description" json:"description"`
	PermissionLabel string     `db:"permission_label" json:"permissionLabel"`
	Action          string     `db:"action" json:"action"`
	Level           int        `db:"level" json:"level"`
	Seq             int        `db:"seq" json:"seq"`
	ParentId        *string    `db:"parent_id" json:"parentId"`
	CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy       *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy       *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted       bool       `db:"is_deleted" json:"isDeleted"`
}

type MenuDTO struct {
	ID              *uuid.UUID `db:"id" json:"id"`
	Name            *string    `db:"name" json:"name"`
	Link            *string    `db:"link" json:"link"`
	Icon            *string    `db:"icon" json:"icon"`
	Description     *string    `db:"description" json:"description"`
	PermissionLabel *string    `db:"permission_label" json:"permissionLabel"`
	Action          *string    `db:"action" json:"action"`
	Level           *int       `db:"level" json:"level"`
	Seq             *int       `db:"seq" json:"seq"`
	ParentId        *string    `db:"parent_id" json:"parentId"`
	ParentMenu      *string    `db:"parent_menu" json:"parentMenu"`
	CreatedAt       *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy       *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy       *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted       bool       `db:"is_deleted" json:"isDeleted"`
}
type RequestMenuFormat struct {
	Name            string    `db:"name" json:"name"`
	Link            string    `db:"link" json:"link"`
	Icon            *string   `db:"icon" json:"icon"`
	Description     *string   `db:"description" json:"description"`
	PermissionLabel string    `db:"permission_label" json:"permissionLabel"`
	Action          string    `db:"action" json:"action"`
	Level           int       `db:"level" json:"level"`
	Seq             int       `db:"seq" json:"seq"`
	ParentId        *string   `db:"parent_id" json:"parentId"`
	UserID          uuid.UUID `json:"-"`
}

type RequestMenuSortFormat struct {
	Id        uuid.UUID `db:"id" json:"id"`
	JenisSort string    `json:"jenisSort"`
}
type Urutan struct {
	Urutan int `db:"urutan" json:"urutan"`
}
type UrutanRequest struct {
	Level    string `json:"level"`
	IdParent string `json:"idParent"`
}

func (menu *Menu) NewMenuFormat(reqFormat RequestMenuFormat) (newMenu Menu, err error) {
	newID, _ := uuid.NewV4()

	newMenu = Menu{
		ID:              newID,
		Name:            reqFormat.Name,
		Link:            reqFormat.Link,
		Description:     reqFormat.Description,
		Icon:            reqFormat.Icon,
		PermissionLabel: reqFormat.PermissionLabel,
		Action:          reqFormat.Action,
		Level:           reqFormat.Level,
		ParentId:        reqFormat.ParentId,
		Seq:             reqFormat.Seq,
		CreatedAt:       time.Now(),
		CreatedBy:       &reqFormat.UserID,
	}
	return
}

var ColumnMappMenu = map[string]interface{}{
	"id":              "m.id",
	"name":            "m.name",
	"link":            "m.link",
	"description":     "m.description",
	"icon":            "m.icon",
	"permissionLabel": "m.permission_label",
	"action":          "m.action",
	"level":           "m.level",
	"seq":             "m.seq",
	"parentMenu":      "cm.name",
	"createdAt":       "m.created_at",
	"updatedAt":       "m.updated_at",
	"createdBy":       "m.created_by",
	"updatedBy":       "m.updated_by",
}

func (menu *Menu) NewFormatUpdate(reqFormat RequestMenuFormat) (err error) {
	now := time.Now()
	menu.Name = reqFormat.Name
	menu.Link = reqFormat.Link
	menu.Description = reqFormat.Description
	menu.Icon = reqFormat.Icon
	menu.PermissionLabel = reqFormat.PermissionLabel
	menu.Action = reqFormat.Action
	menu.Level = reqFormat.Level
	menu.Seq = reqFormat.Seq
	menu.ParentId = reqFormat.ParentId
	menu.UpdatedAt = &now
	menu.UpdatedBy = &reqFormat.UserID

	return nil
}

func (menu *Menu) SoftDelete() {
	now := time.Now()
	menu.IsDeleted = true
	menu.UpdatedAt = &now
}

type MenuRole struct {
	ID              string    `db:"id" json:"id"`
	MenuId          string    `db:"menu_id" json:"menuId"`
	RoleId          string    `db:"role_id" json:"roleId"`
	Permission      *string   `db:"permission" json:"permission"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
	PermissionList  []string  `json:"permissionList"`
	MenuPermissions []string  `json:"menuPermissionList"`
}
type MenuResponse struct {
	ID              string         `db:"id" json:"id"`
	MenuID          string         `db:"menu_id" json:"menuId"`
	RoleID          string         `db:"role_id" json:"roleId"`
	Name            string         `db:"name" json:"name"`
	Link            string         `db:"link" json:"link"`
	Description     *string        `db:"description" json:"description"`
	Icon            *string        `db:"icon" json:"icon"`
	Level           int            `db:"level" json:"level"`
	Seq             int            `db:"seq" json:"seq"`
	PermissionLabel *string        `db:"permission_label" json:"permissionLabel"`
	Permission      *string        `db:"permission" json:"permission"`
	PermissionList  []string       `json:"permissionList"`
	Children        []MenuResponse `json:"children"`
}
type MenuResponseTrx struct {
	ID              *string           `db:"id" json:"id"`
	MenuID          string            `db:"menu_id" json:"menuId"`
	RoleID          string            `db:"role_id" json:"roleId"`
	Name            string            `db:"name" json:"name"`
	Link            string            `db:"link" json:"link"`
	Description     *string           `db:"description" json:"description"`
	Icon            *string           `db:"icon" json:"icon"`
	Level           int               `db:"level" json:"level"`
	Seq             int               `db:"seq" json:"seq"`
	PermissionLabel *string           `db:"permission_label" json:"permissionLabel"`
	Permission      *string           `db:"permission" json:"permission"`
	Action          *string           `db:"action" json:"action"`
	PermissionList  []string          `json:"permissionList"`
	ActionList      []string          `json:"actionList"`
	Children        []MenuResponseTrx `json:"children"`
}

type RequestMenuRoleFormat struct {
	MenuId   []string `db:"menu_id" json:"menuId"`
	Level    int      `db:"level" json:"level"`
	ParentId *string  `db:"parent_id" json:"parentId"`
	RoleId   string   `db:"role_id" json:"roleId"`
}

func (menuRole *MenuRole) NewMenuUserFormat(reqFormat RequestMenuRoleFormat) (newMenuUser []MenuRole, err error) {
	viewPermission := "VIEW"
	for i := 0; i < len(reqFormat.MenuId); i++ {
		newID, _ := uuid.NewV4()
		newMenu := MenuRole{
			ID:         newID.String(),
			MenuId:     reqFormat.MenuId[i],
			RoleId:     reqFormat.RoleId,
			Permission: &viewPermission,
			CreatedAt:  time.Now(),
		}

		newMenuUser = append(newMenuUser, newMenu)
	}

	return
}

type RequestMenuPermissionFormat struct {
	Data []PermissionFormat `json:"data"`
}

type PermissionFormat struct {
	Id         string  `db:"id" json:"id"`
	Permission *string `db:"permission" json:"permission"`
}

type MenuRequest struct {
	RoleId   string `json:"roleId"`
	ParentId string `json:"parentId"`
	Level    int    `json:"level"`
}

type RequestBulkMenuRole struct {
	Data []RequestBulkMenuRoleFormat `json:"data"`
}

type RequestBulkMenuRoleFormat struct {
	Id         string  `db:"id" json:"id"`
	MenuId     string  `db:"menu_id" json:"menuId"`
	RoleId     string  `db:"role_id" json:"roleId"`
	Permission *string `db:"permission" json:"permission"`
}

func (menuRole *MenuRole) NewMenuRoleFormatBulk(reqFormat []RequestBulkMenuRoleFormat) (newMenuRole []MenuRole, err error) {
	for _, v := range reqFormat {
		var detID string
		if v.Id == "" {
			newID, _ := uuid.NewV4()
			detID = newID.String()
		} else {
			detID = v.Id
		}

		newMenu := MenuRole{
			ID:         detID,
			MenuId:     v.MenuId,
			RoleId:     v.RoleId,
			Permission: v.Permission,
			CreatedAt:  time.Now(),
		}

		newMenuRole = append(newMenuRole, newMenu)
	}

	return
}
