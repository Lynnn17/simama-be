//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"lms-be/configs"
	"lms-be/infras"
	"lms-be/internal/domain/auth"
	"lms-be/internal/domain/internship"
	"lms-be/internal/domain/master"
	"lms-be/internal/files"
	"lms-be/internal/handlers"
	"lms-be/transport/http"
	"lms-be/transport/http/middleware"
	"lms-be/transport/http/router"
)

// Wiring for configurations.
var configurations = wire.NewSet(
	configs.Get,
)

// Wiring for persistences.
var persistences = wire.NewSet(
	infras.ProvidePostgreSQLConn,
	infras.ProvideTxManager,
)

// Wiring for all domains.
var domains = wire.NewSet(
	domainAuth,
	domainMaster,
	domainInternship,
	domainTransaction,
	// domainReport,
)

// Wiring for domain Auth
var domainAuth = wire.NewSet(
	// Log System interface and implementation
	auth.ProvideLogSystemServiceImpl,
	wire.Bind(new(auth.LogSystemService), new(*auth.LogSystemServiceImpl)),
	// LogSystemRepository interface and implementation
	auth.ProvideLogSystemRepositoryPostgreSQL,
	wire.Bind(new(auth.LogSystemRepository), new(*auth.LogSystemRepositoryPostgreSQL)),

	// Menu interface and implementation
	auth.ProvideMenuServiceImpl,
	wire.Bind(new(auth.MenuService), new(*auth.MenuServiceImpl)),
	// MenuRepository interface and implementation
	auth.ProvideMenuRepositoryPostgreSQL,
	wire.Bind(new(auth.MenuRepository), new(*auth.MenuRepositoryPostgreSQL)),

	// Role interface and implementation
	auth.ProvideRoleServiceImpl,
	wire.Bind(new(auth.RoleService), new(*auth.RoleServiceImpl)),
	// RoleRepository interface and implementation
	auth.ProvideRoleRepositoryPostgreSQL,
	wire.Bind(new(auth.RoleRepository), new(*auth.RoleRepositoryPostgreSQL)),

	// UserService interface and implementation
	auth.ProvideUserServiceImpl,
	wire.Bind(new(auth.UserService), new(*auth.UserServiceImpl)),
	// UserRepository interface and implementation
	auth.ProvideUserRepositoryPostgreSQL,
	wire.Bind(new(auth.UserRepository), new(*auth.UserRepositoryPostgreSQL)),

	// AppConfigService interface and implementation
	auth.ProvideAppConfigServiceImpl,
	wire.Bind(new(auth.AppConfigService), new(*auth.AppConfigServiceImpl)),
	// AppConfigRepository interface and implementation
	auth.ProvideAppConfigRepositoryPostgreSQL,
	wire.Bind(new(auth.AppConfigRepository), new(*auth.AppConfigRepositoryPostgreSQL)),
)

// Wiring for domain Master
var domainMaster = wire.NewSet(
	///FileService and implementation
	files.ProvideFileServiceImpl,
	wire.Bind(new(files.FileService), new(*files.FileServiceImpl)),

	// AcademicYear interface and implementation
	master.ProvideAcademicYearServiceImpl,
	wire.Bind(new(master.AcademicYearService), new(*master.AcademicYearServiceImpl)),
	// AcademicYearRepository interface and implementation
	master.ProvideAcademicYearRepositoryPostgreSQL,
	wire.Bind(new(master.AcademicYearRepository), new(*master.AcademicYearRepositoryPostgreSQL)),

	// Company interface and implementation
	master.ProvideCompanyServiceImpl,
	wire.Bind(new(master.CompanyService), new(*master.CompanyServiceImpl)),
	// CompanyRepository interface and implementation
	master.ProvideCompanyRepositoryPostgreSQL,
	wire.Bind(new(master.CompanyRepository), new(*master.CompanyRepositoryPostgreSQL)),

	// Department interface and implementation
	master.ProvideDepartmentServiceImpl,
	wire.Bind(new(master.DepartmentService), new(*master.DepartmentServiceImpl)),
	// DepartmentRepository interface and implementation
	master.ProvideDepartmentRepositoryPostgreSQL,
	wire.Bind(new(master.DepartmentRepository), new(*master.DepartmentRepositoryPostgreSQL)),

	// Personnel interface and implementation
	master.ProvidePersonnelServiceImpl,
	wire.Bind(new(master.PersonnelService), new(*master.PersonnelServiceImpl)),
	// PersonnelRepository interface and implementation
	master.ProvidePersonnelRepositoryPostgreSQL,
	wire.Bind(new(master.PersonnelRepository), new(*master.PersonnelRepositoryPostgreSQL)),
)

var domainInternship = wire.NewSet(
	internship.ProvideRegistrationServiceImpl,
	wire.Bind(new(internship.RegistrationService), new(*internship.RegistrationServiceImpl)),
	internship.ProvideRegistrationRepositoryPostgreSQL,
	wire.Bind(new(internship.RegistrationRepository), new(*internship.RegistrationRepositoryPostgreSQL)),
)

var domainTransaction = wire.NewSet()

// var domainReport = wire.NewSet()

// Wiring for HTTP routing.
var routing = wire.NewSet(
	wire.Struct(new(router.DomainHandlers), "*"),
	// Auth
	handlers.ProvideLogSystemHandler,
	handlers.ProvideMenuHandler,
	handlers.ProvideRoleHandler,
	handlers.ProvideUserHandler,
	handlers.ProvideAppConfigHandler,
	// Master
	handlers.ProvideAcademicYearHandler,
	handlers.ProvideCompanyHandler,
	handlers.ProvideDepartmentHandler,
	handlers.ProvidePersonnelHandler,
	handlers.ProvideRegistrationHandler,
	// Transaction
	// File
	handlers.ProvideFileHandler,
	// Import
	handlers.ProvideImportTemplateHandler,
	// JWT
	middleware.ProvideJWTMiddleware,
	router.ProvideRouter,
)

// Wiring for everything.
func InitializeService() *App {
	wire.Build(
		// configurations
		configurations,
		// persistences
		persistences,
		// domains
		domains,
		// routing
		routing,
		// selected transport layer
		http.ProvideHTTP,
		wire.Struct(new(App), "*"),
	)
	return &App{}
}
