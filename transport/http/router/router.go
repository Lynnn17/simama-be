package router

import (
	"github.com/go-chi/chi"

	"lms-be/internal/handlers"
	"lms-be/transport/http/middleware"
)

// DomainHandlers is a struct that contains all domain-specific handlers.
type DomainHandlers struct {
	// Auth
	LogSystemHandler handlers.LogSystemHandler
	MenuHandler      handlers.MenuHandler
	RoleHandler      handlers.RoleHandler
	UserHandler      handlers.UserHandler
	AppConfigHandler handlers.AppConfigHandler
	// Master
	AcademicYearHandler     handlers.AcademicYearHandler
	CompanyHandler          handlers.CompanyHandler
	DepartmentHandler       handlers.DepartmentHandler
	RegistrationHandler     handlers.RegistrationHandler
	MentorAssignmentHandler handlers.MentorAssignmentHandler
	LogbookHandler          handlers.LogbookHandler
	TaskHandler             handlers.TaskHandler
	HRDHandler              handlers.HRDHandler
	// Transaction

	// File
	FileHandler handlers.FileHandler
	// Import
	ImportTemplateHandler handlers.ImportTemplateHandler
}

// Router is the router struct containing handlers.
type Router struct {
	JwtMiddleware  *middleware.JWT
	DomainHandlers DomainHandlers
}

// ProvideRouter is the provider function for this router.
func ProvideRouter(domainHandlers DomainHandlers, jwtMiddleware *middleware.JWT) Router {
	return Router{
		DomainHandlers: domainHandlers,
		JwtMiddleware:  jwtMiddleware,
	}
}

// SetupRoutes sets up all routing for this server.
func (r *Router) SetupRoutes(mux *chi.Mux) {
	mux.Route("/v1", func(rc chi.Router) {
		// Auth
		r.DomainHandlers.LogSystemHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.MenuHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.RoleHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.UserHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.AppConfigHandler.Router(rc, r.JwtMiddleware)
		// Master
		r.DomainHandlers.AcademicYearHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.CompanyHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.DepartmentHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.RegistrationHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.MentorAssignmentHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.LogbookHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.TaskHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.HRDHandler.Router(rc, r.JwtMiddleware)
		// Transaction

		// File
		r.DomainHandlers.FileHandler.Router(rc, r.JwtMiddleware)
		// Import
		r.DomainHandlers.ImportTemplateHandler.Router(rc, r.JwtMiddleware)
	})
}
