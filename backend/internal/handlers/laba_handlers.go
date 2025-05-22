package handlers

import (
	"laba_service/internal/models"
	service "laba_service/internal/services"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type LabaHandler interface {
	GetAll(ctx *gin.Context)
	ExportExcel(ctx *gin.Context)
	ImportExcel(ctx *gin.Context)
}

type labaHandlerImpl struct {
	service service.LabaService
}

func NewproductHandler(service service.LabaService) *labaHandlerImpl {
	return &labaHandlerImpl{service}
}

func (h *labaHandlerImpl) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	laba, err := h.service.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get Laba",
			Error:   true,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Successfully fetched Laba",
		Error:   false,
		Data:    laba,
	})
}

func (h *labaHandlerImpl) ExportExcel(c *gin.Context) {
	ctx := c.Request.Context()

	data, err := h.service.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch data",
			Error:   true,
		})
		return
	}

	sortedKeys := func(m map[string]map[string]float64) []string {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return keys
	}

	f := excelize.NewFile()
	sheet := "Laba"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	f.SetCellValue(sheet, "A1", "NO")
	f.SetCellValue(sheet, "B1", "Label Rekonsiliasi")

	periodMap := make(map[string]struct{})
	for _, periodData := range data {
		for p := range periodData {
			periodMap[p] = struct{}{}
		}
	}

	var periods []string
	for p := range periodMap {
		periods = append(periods, p)
	}
	sort.Strings(periods)

	for i, p := range periods {
		col, _ := excelize.ColumnNumberToName(i + 3)
		f.SetCellValue(sheet, col+"1", p)
	}

	row := 2
	for i, label := range sortedKeys(data) {
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), i+1)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), label)

		for j, p := range periods {
			col, _ := excelize.ColumnNumberToName(j + 3)
			val := data[label][p]
			f.SetCellValue(sheet, col+strconv.Itoa(row), val)
		}
		row++
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=laba_export.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to write excel file",
			Error:   true,
		})
	}
}

func (h *labaHandlerImpl) ImportExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "File is required",
			Error:   true,
		})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to open file",
			Error:   true,
		})
		return
	}
	defer f.Close()

	excelFile, err := excelize.OpenReader(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to parse excel file",
			Error:   true,
		})
		return
	}

	sheetName := excelFile.GetSheetName(0)
	if sheetName == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Excel file has no sheets",
			Error:   true,
		})
		return
	}

	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to read rows",
			Error:   true,
		})
		return
	}

	if len(rows) < 2 {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Excel file must have header and at least one row",
			Error:   true,
		})
		return
	}

	header := rows[0]

	if len(header) < 3 {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid header columns",
			Error:   true,
		})
		return
	}

	periods := header[2:]

	var labas []models.Laba

	for i, row := range rows[1:] {
		if len(row) < 2 {
			continue
		}

		label := row[1]

		for j, periode := range periods {
			if len(row) <= j+2 {
				continue
			}

			nilaiStr := row[j+2]
			if nilaiStr == "" {
				continue
			}

			val, err := strconv.ParseFloat(nilaiStr, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, models.Response{
					Status:  http.StatusBadRequest,
					Message: "Invalid data at row " + strconv.Itoa(i+2) + " column " + strconv.Itoa(j+3),
					Error:   true,
				})
				return
			}

			t, err := time.Parse("2006-01-02", periode)
			if err != nil {
				c.JSON(http.StatusBadRequest, models.Response{
					Status:  http.StatusBadRequest,
					Message: "Invalid periode format in header: " + periode,
					Error:   true,
				})
				return
			}

			laba := models.Laba{
				LabelRekonsiliasiFiskal: label,
				Periode:                 t,
				Nilai:                   val,
			}

			labas = append(labas, laba)
		}
	}

	err = h.service.Create(c.Request.Context(), labas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to save data",
			Error:   true,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Data imported successfully",
		Error:   false,
	})
}
