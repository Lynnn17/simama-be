package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/xuri/excelize/v2"

	"lms-be/configs"
	"lms-be/shared/failure"
	"lms-be/transport/http/middleware"
	"lms-be/transport/http/response"
)

// ImportTemplate defines the structure for an import template
type ImportTemplate struct {
	Headers  []string
	Examples [][]interface{}
}

// ImportTemplates contains all import template configurations
var ImportTemplates = map[string]ImportTemplate{
	"company": {
		Headers:  []string{"Name", "Is Registered Partner", "PIC Contact"},
		Examples: [][]interface{}{{"PT. ABC", true, "John Doe"}, {"PT. XYZ", false, "Jane Smith"}},
	},
	"academic-year": {
		Headers:  []string{"Code", "Name", "Semester Type", "Active"},
		Examples: [][]interface{}{{"2024/2025-1", "Tahun Ajaran 2024/2025 Ganjil", "semester", "true"}, {"2024/2025", "Tahun Ajaran 2024/2025", "year", "true"}},
	},
	"grade": {
		Headers:  []string{"Name", "Level"},
		Examples: [][]interface{}{{"10", "sma"}, {"s1", "s1"}},
	},
	"organization": {
		Headers:  []string{"Name", "Email", "Phone", "Address", "Active"},
		Examples: [][]interface{}{{"Organization A", "contact@org.com", "08123456789", "Jakarta", "true"}},
	},
	"institution": {
		Headers:  []string{"Code", "Name", "Type", "Email", "Phone", "Address", "Active"},
		Examples: [][]interface{}{{"INST001", "SMA Negeri 1", "school", "sma1@email.com", "08123456789", "Jl. Sudirman No. 1", "true"}, {"INST002", "Universitas Negeri 1", "university", "univ1@email.com", "08123456789", "Jl. Sudirman No. 2", "true"}},
	},
	"study-program": {
		Headers:  []string{"Code", "Name", "Degree"},
		Examples: [][]interface{}{{"RPL", "Rekaya Perangkat Lunak", "smk"}, {"TI", "Teknik Informatika", "s1"}},
	},
	"class": {
		Headers:  []string{"Code", "Name", "Grade", "Academic Year Code", "Study Program Code"},
		Examples: [][]interface{}{{"X-RPL1", "10 RPL 1", "10", "2024/2025-1", "RPL"}, {"TI-2A", "TI 2A", "s1", "2024/2025", "TI"}},
	},
	"course-subject": {
		Headers:  []string{"Code", "Name"},
		Examples: [][]interface{}{{"MTK001", "Matematika"}, {"BD001", "Basis Data"}},
	},
	"student": {
		Headers:  []string{"No. ID", "Name", "Email", "Phone", "Address", "Entry Year"},
		Examples: [][]interface{}{{"STU2024001", "John Doe", "john@example.com", "081234567890", "Jl. Merdeka No. 10", 2024}},
	},
	"teacher": {
		Headers:  []string{"No. ID", "Name", "Email", "Phone", "Address", "Entry Year"},
		Examples: [][]interface{}{{"TCH2020001", "Jane Smith", "jane@example.com", "081987654321", "Jl. Pendidikan No. 5", 2020}},
	},
}

type ImportTemplateHandler struct {
	Config *configs.Config
}

func ProvideImportTemplateHandler(config *configs.Config) ImportTemplateHandler {
	return ImportTemplateHandler{Config: config}
}

func (h *ImportTemplateHandler) Router(r chi.Router, mw *middleware.JWT) {
	r.Route("/import/template", func(r chi.Router) {
		r.Use(mw.VerifyToken)
		r.Get("/{domain}", h.Download)
	})
}

// Download generates and serves Excel template for the specified domain.
// @Summary Download import template
// @Description Download Excel template for importing data into the specified domain.
// @Tags Import
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security BearerAuth
// @Param domain path string true "Domain name (academic_year, grade, organization, institution, class, course_subject, student, study_program, teacher)"
// @Success 200 {file} file "Excel template file"
// @Failure 404 {object} response.Base
// @Router /v1/import/template/{domain} [get]
func (h *ImportTemplateHandler) Download(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")

	template, ok := ImportTemplates[domain]
	if !ok {
		response.WithError(w, failure.NotFound("template not found for domain: "+domain))
		return
	}

	f := excelize.NewFile()
	sheet := "Sheet1"

	// Set headers (row 1)
	for i, header := range template.Headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	// Bold headers
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	f.SetRowStyle(sheet, 1, 1, style)

	// Set example data (row 2+)
	for rowIdx, example := range template.Examples {
		for colIdx, val := range example {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheet, cell, val)
		}
	}

	// Auto-width columns (approximate)
	for i := range template.Headers {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheet, col, col, 20)
	}

	// Serve file
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+domain+"_import_template.xlsx")
	f.Write(w)
}
