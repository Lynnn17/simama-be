package internship

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"

	"lms-be/configs"
	sendemail "lms-be/external/email"
	"lms-be/internal/domain/auth"
	"lms-be/shared/logger"
	"lms-be/shared/pagination"
	"lms-be/shared/random"
	"lms-be/shared/socket"
	"time"
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

	// Check if email already exists in users table
	existUser, _ := s.UserRepository.ExistByEmail(req.Email)
	if existUser {
		return Registration{}, errors.New("Email ini sudah terdaftar dalam sistem")
	}

	// Check if email already exists in registrations table (active/pending)
	existReg, _ := s.RegistrationRepository.ExistByEmail(ctx, req.Email)
	if existReg {
		return Registration{}, errors.New("Email ini sudah memiliki pendaftaran yang sedang diproses")
	}

	newRegistration, _ = newRegistration.NewRegistrationFormat(req)
	err = s.RegistrationRepository.Create(ctx, &newRegistration)
	if err != nil {
		return Registration{}, err
	}

	// Send Real-time Notification to HRD (Role HA01)
	hub := socket.GetInstance()
	notificationMsg := map[string]interface{}{
		"title":   "Pendaftar Baru",
		"message": "Ada lamaran magang baru dari " + req.FullName,
		"type":    "registration",
	}
	hub.BroadcastToRole("HA01", "new_notification", notificationMsg)
	hub.BroadcastToRole("HA01", "refresh_registrations", nil)

	// Send email notification with Retry Logic
	go func() {
		subject := "Konfirmasi Penerimaan Lamaran Magang - PT Greatsoft Solusi Indonesia"
		now := time.Now().Format("02-01-2006 15:04")
		message := "Halo " + req.FullName + ",\n\n" +
			"Lamaran magang Anda telah berhasil masuk ke dalam sistem kami.\n\n" +
			"Berikut adalah pembaruan informasi Anda:\n" +
			"Status Terkini: Diproses (Pending)\n" +
			"Tanggal Penetapan Status: " + now + " WIB\n\n" +
			"Tim HRD kami akan segera mereview dokumen dan CV yang telah Anda unggah. Mohon menunggu informasi selanjutnya yang akan kami sampaikan melalui email ini.\n\n" +
			"Salam, Tim HRD PT Greatsoft Solusi Indonesia"

		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			errMail := sendemail.PrimitiveSendMail(s.Config, []string{req.Email}, []string{}, subject, message)
			if errMail == nil {
				break
			}
			logger.Error(fmt.Sprintf("Failed to send initial email to %s (attempt %d): %v", req.Email, i+1, errMail))
			if i < maxRetries-1 {
				time.Sleep(2 * time.Second)
			}
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
		// New context for background process
		bgCtx := context.Background()

		if status == "accepted" {
			// Generate password
			password, _ := random.RandStringBytes(8, random.AlphanumericMixed)

			// Generate username from email
			username := existing.Email
			if atIndex := strings.Index(existing.Email, "@"); atIndex != -1 {
				username = existing.Email[:atIndex]
			}

			// Check if user with this email already exists
			existEmail, _ := s.UserRepository.ExistByEmail(existing.Email)
			if existEmail {
				// Resolve the existing user to get their ID
				existingUser, errResolve := s.UserRepository.ResolveUserByEmailRole(existing.Email)
				if errResolve == nil {
					// Link the existing user to the registration
					newID, _ := uuid.FromString(existingUser.ID.String())
					_ = s.RegistrationRepository.UpdateUserID(bgCtx, existing.ID, newID)
				}
			} else {
				// Check if username exists
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
					RoleId:   "HA02",
					Status:   &userStatus,
					Active:   false,
				}

				user := inputUser.CreateUser()
				errCreate := s.UserRepository.TransactionCreateUser(user)
				if errCreate == nil {
					_ = s.RegistrationRepository.UpdateUserID(bgCtx, existing.ID, uuid.UUID(user.ID))
				}
			}

			// Send Email with Retry Logic
			subject := "Pengumuman Seleksi Program Magang - PT Greatsoft Solusi Indonesia"
			now := time.Now().Format("02-01-2006 15:04")
			message := "Halo " + existing.FullName + ",\n\n" +
				"Selamat! Kami dengan senang hati menginformasikan bahwa Anda telah lulus seleksi program magang di PT Greatsoft Solusi Indonesia.\n\n" +
				"Berikut adalah pembaruan informasi Anda:\n" +
				"Status Terkini: Diterima (Accepted)\n" +
				"Tanggal Penetapan Status: " + now + " WIB\n\n" +
				"Informasi lebih lanjut mengenai persiapan hari pertama, pembuatan akun sistem SIMAMA, serta penugasan Mentor Anda akan segera kami informasikan lebih lanjut.\n\n" +
				"Salam, Tim HRD PT Greatsoft Solusi Indonesia"

			maxRetries := 3
			isSent := false
			for i := 0; i < maxRetries; i++ {
				errMail := sendemail.PrimitiveSendMail(s.Config, []string{existing.Email}, []string{}, subject, message)
				if errMail == nil {
					isSent = true
					break
				}
				logger.Error(fmt.Sprintf("Failed to send email to %s (attempt %d): %v", existing.Email, i+1, errMail))
				if i < maxRetries-1 {
					time.Sleep(2 * time.Second)
				}
			}

			if !isSent {
				_ = s.RegistrationRepository.UpdateEmailStatus(bgCtx, existing.ID, false)
			}

		} else if status == "rejected" {
			// Send rejection email with similar retry logic
			subject := "Pengumuman Seleksi Program Magang - PT Greatsoft Solusi Indonesia"
			now := time.Now().Format("02-01-2006 15:04")
			message := "Halo " + existing.FullName + ",\n\n" +
				"Terima kasih atas ketertarikan dan antusiasme Anda untuk bergabung dalam program magang di PT Greatsoft Solusi Indonesia, serta atas waktu yang Anda luangkan untuk melengkapi proses pendaftaran.\n\n" +
				"Berikut adalah pembaruan informasi Anda:\n" +
				"Status Terkini: Ditolak (Rejected)\n" +
				"Tanggal Penetapan Status: " + now + " WIB\n\n" +
				"Setelah melalui proses review dengan saksama, mohon maaf kami belum dapat menerima Anda pada periode magang kali ini. Kami berharap Anda sukses dalam perjalanan karier dan pendidikan ke depannya.\n\n" +
				"Salam, Tim HRD PT Greatsoft Solusi Indonesia"

			maxRetries := 3
			isSent := false
			for i := 0; i < maxRetries; i++ {
				errMail := sendemail.PrimitiveSendMail(s.Config, []string{existing.Email}, []string{}, subject, message)
				if errMail == nil {
					isSent = true
					break
				}
				if i < maxRetries-1 {
					time.Sleep(2 * time.Second)
				}
			}

			if !isSent {
				_ = s.RegistrationRepository.UpdateEmailStatus(bgCtx, existing.ID, false)
			}
		}
	}()

	return newRegistration, nil
}
