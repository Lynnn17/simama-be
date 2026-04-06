package auth

import (
	"errors"

	"lms-be/configs"
	"lms-be/shared/model"
	"lms-be/shared/pagination"
)

// RoleService adalah interface RoleService untuk entity Role
type RoleService interface {
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	GetAllData() ([]Role, error)
	CreateRole(role Role) (bool, error)
	UpdateRole(id string, role Role) error
	DeleteRole(id string) error
	ResolveByID(id string) (Role, error)
}

// RoleServiceImpl adalah implementasi dari service yang digunakan untuk entity Role
type RoleServiceImpl struct {
	RoleRepository RoleRepository
	Config         *configs.Config
}

// ProvideRoleServiceImpl adalah provider untuk service RoleService
func ProvideRoleServiceImpl(roleRepository RoleRepository, config *configs.Config) *RoleServiceImpl {
	s := new(RoleServiceImpl)
	s.RoleRepository = roleRepository
	s.Config = config
	return s
}

func (s *RoleServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.RoleRepository.ResolveAll(request)
}

// ResolveAll get all Role data
func (s *RoleServiceImpl) GetAllData() ([]Role, error) {
	return s.RoleRepository.GetAllData()
}

// CreateRole is the service to create Role entity
func (r *RoleServiceImpl) CreateRole(role Role) (bool, error) {

	exist, err := r.RoleRepository.ExistRoleByName(role.Name)
	if exist {
		return exist, errors.New("Nama Role sudah dipakai")
	}
	if err != nil {
		return exist, err
	}
	newRole, _ := role.NewRoleFormat(role)
	err = r.RoleRepository.CreateRole(newRole)
	if err != nil {
		return exist, err
	}
	return exist, nil
}

// UpdateRole aalah service yang digunakan untuk mengubah data Role
func (r *RoleServiceImpl) UpdateRole(id string, newRole Role) error {
	role, err := r.RoleRepository.ResolveRoleByID(id)
	if err != nil || (Role{}) == role {
		return errors.New("Role dengan nama :" + role.Name + " tidak ditemukan")
	}

	return r.RoleRepository.UpdateRole(newRole)
}

// DeleteRole adalah service yang digunakan untuk menghapus data Role
func (r *RoleServiceImpl) DeleteRole(id string) error {
	role, err := r.RoleRepository.ResolveRoleByID(id)
	if err != nil || (Role{}) == role {
		return errors.New("Role dengan nama :" + role.Name + " tidak ditemukan")
	}
	err = r.RoleRepository.DeleteRoleByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Role dengan nama: " + role.Name)
	}
	return nil
}

// ResolveByRoleID adalah service yang digunakan untuk mendapatkan Role berdasarkan RoleID
func (r *RoleServiceImpl) ResolveByID(roleId string) (roleMenu Role, err error) {
	role, err := r.RoleRepository.ResolveRoleByID(roleId)
	if err != nil {
		return Role{}, errors.New("Role dengan nama: " + role.Name + " tidak ditemukan")
	}
	return
}
