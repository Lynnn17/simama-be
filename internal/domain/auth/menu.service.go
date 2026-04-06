package auth

import (
	"errors"
	"fmt"
	"strings"

	"lms-be/shared/logger"
	"lms-be/shared/model"
	"lms-be/shared/pagination"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

// MenuService adalah interface MenuService untuk entity menu
type MenuService interface {
	GetAllMenu() ([]Menu, error)
	ResolveAll(request model.StandardRequestMenu) (orders pagination.Response, err error)
	ResolveMenuByRoleID(req MenuRequest) (menu []MenuResponse, err error)
	ResolveMenuByRoleIDTrx(req MenuRequest) (menu []MenuResponseTrx, err error)
	CreateMenu(reqFormat RequestMenuFormat) (menu Menu, error error)
	UpdateMenu(id uuid.UUID, newMenu RequestMenuFormat) (menu Menu, error error)
	ResolveMenuByID(id uuid.UUID) (menu Menu, error error)
	DeleteMenuByID(id uuid.UUID) error
	ResolveMenuRoleByID(id uuid.UUID) (menuRole MenuRole, error error)
	UpdateMenuPermission(req RequestMenuPermissionFormat) (err error)
	CreateBulkMenuRole(reqFormat RequestBulkMenuRole) (newMenuRole []MenuRole, err error)
}

// MenuServiceImpl adalah implementasi dari service yang digunakan untuk entity Role
type MenuServiceImpl struct {
	MenuRepository MenuRepository
}

// ProvideServiceImpl adalah provider untuk service MenuService
func ProvideMenuServiceImpl(MenuRepository MenuRepository) *MenuServiceImpl {
	s := new(MenuServiceImpl)
	s.MenuRepository = MenuRepository
	return s
}

// ResolveByMenuRoleID adalah service yang digunakan untuk mendapatkan menu berdasarkan RoleID
func (r *MenuServiceImpl) ResolveMenuByRoleID(req MenuRequest) (menu []MenuResponse, err error) {
	menu, err = r.MenuRepository.ResolveMenuByRoleID(req)
	if err != nil {
		return []MenuResponse{}, errors.New("Ada kesalahan waktu get menu berdasarkan roleID: " + req.RoleId)
	}

	for i := 0; i < len(menu); i++ {
		d := menu[i]
		if d.Permission != nil {
			permissionArr := strings.Split(*d.Permission, ",")
			if len(permissionArr) > 0 {
				permissionList := make([]string, 0)
				for _, v := range permissionArr {
					permissionList = append(permissionList, model.ParseString(d.PermissionLabel)+"."+v)
				}
				menu[i].PermissionList = permissionList
			}
		}

		var CMenu []MenuResponse
		reqChild := MenuRequest{
			RoleId:   req.RoleId,
			ParentId: d.MenuID,
		}
		CMenu, err = r.MenuRepository.ResolveMenuByParentID(reqChild)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}

		for j := 0; j < len(CMenu); j++ {
			v := CMenu[j]
			if v.Permission != nil {
				permissionArr := strings.Split(*v.Permission, ",")
				if len(permissionArr) > 0 {
					permissionList := make([]string, 0)
					for _, t := range permissionArr {
						permissionList = append(permissionList, model.ParseString(v.PermissionLabel)+"."+t)
					}
					CMenu[j].PermissionList = permissionList
				}
			}
		}

		menu[i].Children = CMenu
	}

	return
}

func (r *MenuServiceImpl) ResolveMenuByRoleIDTrx(req MenuRequest) (menu []MenuResponseTrx, err error) {
	menu, err = r.MenuRepository.ResolveMenuByRoleIDTrx(req)
	if err != nil {
		return []MenuResponseTrx{}, errors.New("Ada kesalahan waktu get menu berdasarkan roleID: " + req.RoleId)
	}

	for i := 0; i < len(menu); i++ {
		d := menu[i]
		if d.Permission != nil {
			permissionArr := strings.Split(*d.Permission, ",")
			menu[i].PermissionList = permissionArr
		}

		if d.Action != nil {
			actionArr := strings.Split(*d.Action, ",")
			menu[i].ActionList = actionArr
		}

		var CMenu []MenuResponseTrx
		reqChild := MenuRequest{
			RoleId:   req.RoleId,
			ParentId: d.MenuID,
		}
		CMenu, err = r.MenuRepository.ResolveMenuByParentIDTrx(reqChild)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}

		for j := 0; j < len(CMenu); j++ {
			v := CMenu[j]
			if v.Permission != nil {
				permissionArr := strings.Split(*v.Permission, ",")
				CMenu[j].PermissionList = permissionArr
			}

			if v.Action != nil {
				actionArr := strings.Split(*v.Action, ",")
				CMenu[j].ActionList = actionArr
			}
		}

		menu[i].Children = CMenu
	}

	return
}

func (s *MenuServiceImpl) GetAllMenu() (data []Menu, err error) {
	return s.MenuRepository.GetAllMenu()
}

func (s *MenuServiceImpl) ResolveAll(request model.StandardRequestMenu) (orders pagination.Response, err error) {
	return s.MenuRepository.ResolveAll(request)
}

// CreateMenu is the service to create Menu entity
func (s *MenuServiceImpl) CreateMenu(reqFormat RequestMenuFormat) (newMenu Menu, err error) {
	if err != nil {
		return Menu{}, err
	}
	newMenu, err = newMenu.NewMenuFormat(reqFormat)

	err = s.MenuRepository.CreateMenu(newMenu)
	if err != nil {
		return Menu{}, err
	}
	return newMenu, nil
}

func (s *MenuServiceImpl) ResolveMenuByID(id uuid.UUID) (menu Menu, err error) {
	menu, err = s.MenuRepository.ResolveMenuByID(id)
	if err != nil {
		return
	}
	return
}

// UpdateMenu adalah service yang digunakan untuk mengupdate Menu
func (s *MenuServiceImpl) UpdateMenu(id uuid.UUID, newMenu RequestMenuFormat) (menu Menu, err error) {
	menu, err = s.MenuRepository.ResolveMenuByID(id)
	if err != nil {
		return Menu{}, errors.New("Menu dengan ID :" + id.String() + " tidak ditemukan")
	}
	menu.NewFormatUpdate(newMenu)
	log.Info().Msgf("service.UpdateMenu %s", menu)

	err = s.MenuRepository.UpdateMenu(menu)

	if err != nil {
		log.Error().Msgf("service.UpdateMenu error", err)
	}
	return
}

func (s *MenuServiceImpl) DeleteMenuByID(id uuid.UUID) error {
	menu, err := s.MenuRepository.ResolveMenuByID(id)
	if err != nil || (Menu{}) == menu {
		return errors.New("menu dengan ID :" + id.String() + " tidak ditemukan")
	}

	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data menu dengan ID: " + id.String())
	}
	menu.SoftDelete()

	err = s.MenuRepository.UpdateMenu(menu)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data menu dengan ID: " + id.String())
	}
	return nil
}

/* ================ FUNCTION MENU ROLE ACCESS ================== */

func (s *MenuServiceImpl) ResolveMenuRoleByID(id uuid.UUID) (menuRole MenuRole, err error) {
	menuRole, err = s.MenuRepository.ResolveMenuRoleByID(id)
	if err != nil {
		return
	}
	return
}

// UpdateMenu adalah service yang digunakan untuk mengupdate Menu
func (s *MenuServiceImpl) UpdateMenuPermission(req RequestMenuPermissionFormat) (err error) {
	for _, v := range req.Data {
		menuUser := MenuRole{
			ID:         v.Id,
			Permission: v.Permission,
		}

		err = s.MenuRepository.UpdatePermission(menuUser)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (s *MenuServiceImpl) CreateBulkMenuRole(reqFormat RequestBulkMenuRole) (newMenuRole []MenuRole, err error) {
	var menuRole MenuRole
	newMenuRole, _ = menuRole.NewMenuRoleFormatBulk(reqFormat.Data)
	err = s.MenuRepository.CreateBulkMenuRole(newMenuRole)
	if err != nil {
		return []MenuRole{}, err
	}

	return newMenuRole, nil
}
