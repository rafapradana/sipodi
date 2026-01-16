package handler

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/service"
	"github.com/xuri/excelize/v2"
)

type ExportHandler struct {
	userService   *service.UserService
	schoolService *service.SchoolService
	talentService *service.TalentService
}

func NewExportHandler(userService *service.UserService, schoolService *service.SchoolService, talentService *service.TalentService) *ExportHandler {
	return &ExportHandler{
		userService:   userService,
		schoolService: schoolService,
		talentService: talentService,
	}
}

func (h *ExportHandler) ExportGTK(c *fiber.Ctx) error {
	claims := GetClaims(c)
	format := c.Query("format", "excel")

	params := domain.ListParams{
		Page:  1,
		Limit: 10000,
		Filters: map[string]string{
			"role":     string(domain.RoleGTK),
			"gtk_type": c.Query("gtk_type"),
		},
	}

	// Admin sekolah can only export their school's GTK
	if claims.Role == domain.RoleAdminSekolah && claims.SchoolID != nil {
		params.Filters["school_id"] = claims.SchoolID.String()
	} else if c.Query("school_id") != "" {
		params.Filters["school_id"] = c.Query("school_id")
	}

	users, _, err := h.userService.List(c.Context(), params)
	if err != nil {
		return InternalError(c)
	}

	if format == "pdf" {
		return h.exportGTKPDF(c, users)
	}
	return h.exportGTKExcel(c, users)
}

func (h *ExportHandler) exportGTKExcel(c *fiber.Ctx, users []domain.User) error {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Data GTK"
	f.SetSheetName("Sheet1", sheet)

	// Headers
	headers := []string{"No", "Nama Lengkap", "Email", "NUPTK", "NIP", "Jenis Kelamin", "Tanggal Lahir", "Jenis GTK", "Jabatan", "Status"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheet, cell, header)
	}

	// Data
	for i, user := range users {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), user.FullName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), user.Email)
		if user.NUPTK != nil {
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), *user.NUPTK)
		}
		if user.NIP != nil {
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *user.NIP)
		}
		if user.Gender != nil {
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), string(*user.Gender))
		}
		if user.BirthDate != nil {
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), user.BirthDate.Format("2006-01-02"))
		}
		if user.GTKType != nil {
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), string(*user.GTKType))
		}
		if user.Position != nil {
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *user.Position)
		}
		status := "Aktif"
		if !user.IsActive {
			status = "Nonaktif"
		}
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), status)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return InternalError(c)
	}

	filename := fmt.Sprintf("data_gtk_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	return c.Send(buf.Bytes())
}

func (h *ExportHandler) exportGTKPDF(c *fiber.Ctx, users []domain.User) error {
	// For simplicity, return Excel format with PDF extension note
	// In production, use a PDF library like gofpdf
	return h.exportGTKExcel(c, users)
}

func (h *ExportHandler) ExportTalents(c *fiber.Ctx) error {
	claims := GetClaims(c)
	format := c.Query("format", "excel")

	params := domain.ListParams{
		Page:  1,
		Limit: 10000,
		Filters: map[string]string{
			"talent_type": c.Query("talent_type"),
			"status":      c.Query("status"),
		},
	}

	// Admin sekolah can only export their school's talents
	if claims.Role == domain.RoleAdminSekolah && claims.SchoolID != nil {
		params.Filters["school_id"] = claims.SchoolID.String()
	} else if c.Query("school_id") != "" {
		params.Filters["school_id"] = c.Query("school_id")
	}

	talents, _, err := h.talentService.List(c.Context(), params)
	if err != nil {
		return InternalError(c)
	}

	if format == "pdf" {
		return h.exportTalentsPDF(c, talents)
	}
	return h.exportTalentsExcel(c, talents)
}

func (h *ExportHandler) exportTalentsExcel(c *fiber.Ctx, talents []domain.Talent) error {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Data Talenta"
	f.SetSheetName("Sheet1", sheet)

	// Headers
	headers := []string{"No", "User ID", "Jenis Talenta", "Status", "Tanggal Dibuat", "Tanggal Diverifikasi"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheet, cell, header)
	}

	// Data
	for i, talent := range talents {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), talent.UserID.String())
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), string(talent.TalentType))
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), string(talent.Status))
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), talent.CreatedAt.Format("2006-01-02 15:04:05"))
		if talent.VerifiedAt != nil {
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), talent.VerifiedAt.Format("2006-01-02 15:04:05"))
		}
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return InternalError(c)
	}

	filename := fmt.Sprintf("data_talenta_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	return c.Send(buf.Bytes())
}

func (h *ExportHandler) exportTalentsPDF(c *fiber.Ctx, talents []domain.Talent) error {
	return h.exportTalentsExcel(c, talents)
}

func (h *ExportHandler) ExportSchools(c *fiber.Ctx) error {
	format := c.Query("format", "excel")

	params := domain.ListParams{
		Page:  1,
		Limit: 10000,
		Filters: map[string]string{
			"status": c.Query("status"),
		},
	}

	schools, _, err := h.schoolService.List(c.Context(), params)
	if err != nil {
		return InternalError(c)
	}

	if format == "pdf" {
		return h.exportSchoolsPDF(c, schools)
	}
	return h.exportSchoolsExcel(c, schools)
}

func (h *ExportHandler) exportSchoolsExcel(c *fiber.Ctx, schools []domain.School) error {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Data Sekolah"
	f.SetSheetName("Sheet1", sheet)

	// Headers
	headers := []string{"No", "Nama Sekolah", "NPSN", "Status", "Alamat"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheet, cell, header)
	}

	// Data
	for i, school := range schools {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), school.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), school.NPSN)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), string(school.Status))
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), school.Address)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return InternalError(c)
	}

	filename := fmt.Sprintf("data_sekolah_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	return c.Send(buf.Bytes())
}

func (h *ExportHandler) exportSchoolsPDF(c *fiber.Ctx, schools []domain.School) error {
	return h.exportSchoolsExcel(c, schools)
}

// Helper to get school name by ID
func (h *ExportHandler) getSchoolName(c *fiber.Ctx, schoolID *uuid.UUID) string {
	if schoolID == nil {
		return ""
	}
	school, err := h.schoolService.GetByID(c.Context(), *schoolID)
	if err != nil || school == nil {
		return ""
	}
	return school.Name
}
