package internship

import (
	"context"
	"errors"
	"strings"

	"github.com/gofrs/uuid/v5"

	"lms-be/configs"
	sendemail "lms-be/external/email"
	"lms-be/internal/domain/auth"
	"lms-be/shared/logger"
	"lms-be/shared/pagination"
	"lms-be/shared/random"
)

type RegistrationService interface {
	Create(ctx context.Context, req RequestRegistrationFormat) (newRegistration Registration, err error)
	GetAll(ctx context.Context, req RequestRegistrationListFormat) (data []RegistrationDTO, err error)
	ResolveAll(ctx context.Context, req RequestRegistrationListFormat) (data pagination.Response, err error)
	UpdateStatus(ctx context.Context, id uuid.UUID, req RequestUpdateRegistrationStatusFormat) (newRegistration Registration, err error)
}

type RegistrationServiceImpl struct {
	RegistrationRepository RegistrationRepository
	Config                 *configs.Config
	UserRepository         auth.UserRepository
}

func ProvideRegistrationServiceImpl(repository RegistrationRepository, config *configs.Config, userRepository auth.UserRepository) *RegistrationServiceImpl {
	s := new(RegistrationServiceImpl)
	s.RegistrationRepository = repository
	s.Config = config
	s.UserRepository = userRepository
	return s
}

func (s *RegistrationServiceImpl) Create(ctx context.Context, req RequestRegistrationFormat) (newRegistration Registration, err error) {
	if req.UserID != nil && *req.UserID != uuid.Nil {
		exist, err := s.RegistrationRepository.ExistByUserID(ctx, *req.UserID)
		if err != nil {
			return Registration{}, err
		}
		if exist {
			return Registration{}, errors.New("registration already exists")
		}
	}

	newRegistration, _ = newRegistration.NewRegistrationFormat(req)
	err = s.RegistrationRepository.Create(ctx, &newRegistration)
	if err != nil {
		return Registration{}, err
	}

	// Send email notification
	go func() {
		subject := "Pendaftaran Magang Berhasil - SIMAMA"
		message := "Halo " + req.FullName + ",\n\n" +
			"Pendaftaran magang Anda di SIMAMA telah berhasil diterima dan sedang dalam proses peninjauan.\n" +
			"Silakan tunggu informasi selanjutnya melalui email.\n\n" +
			"Terima kasih."

		err = sendemail.PrimitiveSendMail(s.Config, []string{req.Email}, []string{}, subject, message)
		if err != nil {
			logger.ErrorWithStack(err)
		}
	}()

	return newRegistration, nil
}

func (s *RegistrationServiceImpl) GetAll(ctx context.Context, req RequestRegistrationListFormat) (data []RegistrationDTO, err error) {
	data, err = s.RegistrationRepository.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}

	// Auto-healing for missing userIds
	for i := range data {
		if data[i].Status == "accepted" && (data[i].UserID == nil || *data[i].UserID == uuid.Nil) {
			users, _ := s.UserRepository.ResolveUserByUsername(data[i].Email)
			if len(users) > 0 {
				newUserID := uuid.UUID(users[0].ID)
				_ = s.RegistrationRepository.UpdateUserID(ctx, data[i].ID, newUserID)
				data[i].UserID = &newUserID
			}
		}
	}

	return data, nil
}

func (s *RegistrationServiceImpl) ResolveAll(ctx context.Context, req RequestRegistrationListFormat) (data pagination.Response, err error) {
	data, err = s.RegistrationRepository.ResolveAll(ctx, req)
	if err != nil {
		return data, err
	}

	// Auto-healing for missing userIds in paginated list
	for i := range data.Items {
		reg, ok := data.Items[i].(RegistrationDTO)
		if ok && reg.Status == "accepted" && (reg.UserID == nil || *reg.UserID == uuid.Nil) {
			users, _ := s.UserRepository.ResolveUserByUsername(reg.Email)
			if len(users) > 0 {
				newUserID := uuid.UUID(users[0].ID)
				_ = s.RegistrationRepository.UpdateUserID(ctx, reg.ID, newUserID)
				reg.UserID = &newUserID
				data.Items[i] = reg
			}
		}
	}

	return data, nil
}

func (s *RegistrationServiceImpl) UpdateStatus(ctx context.Context, id uuid.UUID, req RequestUpdateRegistrationStatusFormat) (newRegistration Registration, err error) {
	if id == uuid.Nil {
		return Registration{}, errors.New("registration id is required")
	}
	if req.UserID == nil || *req.UserID == uuid.Nil {
		return Registration{}, errors.New("reviewer id is required")
	}
	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status != "accepted" && status != "rejected" {
		return Registration{}, errors.New("invalid registration status")
	}

	existing, err := s.RegistrationRepository.ResolveByID(ctx, id)
	if err != nil {
		return Registration{}, err
	}
	if existing.ID == uuid.Nil {
		return Registration{}, errors.New("registration not found")
	}
	req.Status = status
	newRegistration, _ = existing.UpdateStatusFormat(req)
	err = s.RegistrationRepository.UpdateStatus(ctx, newRegistration)
	if err != nil {
		return Registration{}, err
	}

	// Async status handling (Email & User Creation)
	go func() {
		if status == "accepted" {
			// Generate password
			password, _ := random.RandStringBytes(8, random.AlphanumericMixed)

			// Generate username from email
			username := existing.Email
			if atIndex := strings.Index(existing.Email, "@"); atIndex != -1 {
				username = existing.Email[:atIndex]
			}

			// Check if username exists, if so append random suffix
			exist, _ := s.UserRepository.ExistByUsername(username)
			if exist {
				randSuffix, _ := random.RandStringBytes(4, random.AlphanumericLowercase)
				username = username + randSuffix
			}

			// Prepare User data
			userStatus := "1"
			inputUser := auth.InputUser{
				Name:     existing.FullName,
				Username: username,
				Email:    &existing.Email,
				Password: password,
				RoleId:   "HA02", // Mahasiswa
				Status:   &userStatus, // Use "1" for varchar(1) column
				Active:   true,
			}

			user := inputUser.CreateUser()
			errCreate := s.UserRepository.TransactionCreateUser(user)
			if errCreate != nil {
				logger.ErrorWithStack(errCreate)
				return
			}

			// Link the user to the registration
			errUpdate := s.RegistrationRepository.UpdateUserID(ctx, existing.ID, uuid.UUID(user.ID))
			if errUpdate != nil {
				logger.ErrorWithStack(errUpdate)
			}

			// Send Email
			subject := "Selamat! Pendaftaran Magang Anda Diterima - SIMAMA"
			message := "Halo " + existing.FullName + ",\n\n" +
				"Selamat! Pendaftaran magang Anda di SIMAMA telah DITERIMA.\n" +
				"Berikut adalah akun untuk login ke portal SIMAMA:\n\n" +
				"Username: " + username + "\n" +
				"Password: " + password + "\n\n" +
				"Silakan login dan segera ubah password Anda di portal SIMAMA.\n\n" +
				"Terima kasih."

			_ = sendemail.PrimitiveSendMail(s.Config, []string{existing.Email}, []string{}, subject, message)

		} else if status == "rejected" {
			// Send rejection email
			subject := "Update Pendaftaran Magang - SIMAMA"
			message := "Halo " + existing.FullName + ",\n\n" +
				"Mohon maaf, pendaftaran magang Anda di SIMAMA belum dapat kami terima saat ini.\n" +
				"Terima kasih atas minat Anda.\n\n" +
				"Salam,\nTIM SIMAMA"

			_ = sendemail.PrimitiveSendMail(s.Config, []string{existing.Email}, []string{}, subject, message)
		}
	}()

	return newRegistration, nil
}
