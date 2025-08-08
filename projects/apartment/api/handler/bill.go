package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/dto"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/domain"
	billPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	appjwt "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/jwt"
	"go.uber.org/zap"
)

const (
	B = 1 << (10 * iota)
	KiB
	MiB
	GiB
)

const maxFileSize = 1 * MiB
const dateLayout = "2006-01-02"

// AddBill
//
// @Summary      Add a new bill
// @Description  Adds a new bill to an apartment. Accepts multipart/form-data for image upload.
// @Tags         Bill
// @Accept       multipart/form-data
// @Produce      json
// @Security 	 BearerAuth
// @Param        name         formData  string  true   "Bill Name"
// @Param        type         formData  string  true   "Bill Type"
// @Param        billNumber   formData  integer true   "Bill Number"
// @Param        amount       formData  integer true   "Amount"
// @Param        dueDate      formData  string  true   "Due Date (YYYY-MM-DD)"
// @Param        status       formData  string  true   "Payment Status"
// @Param        paidAt       formData  string  false  "Paid At (YYYY-MM-DD)"
// @Param        apartmentID  formData  string  true   "Apartment ID"
// @Param        image        formData  file    false  "Bill Image"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  dto.Error
// @Failure      500   {object}  dto.Error
// @Router       /api/v1/bill [post]
func AddBill(svcGetter ServiceGetter[billPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())

		if err := r.ParseMultipartForm(1 * MiB); err != nil {
			log.Error("failed to parse multipart form", zap.Error(err))
			Error(w, r, http.StatusBadRequest, err.Error())
			return
		}

		log.Info("add bill", zap.Any("form values", r.MultipartForm.Value))

		var b domain.Bill

		// Parse and validate form fields
		b.Name = r.FormValue("name")

		if bt := domain.BillType(r.FormValue("type")); bt.IsValid() {
			b.Type = bt
		} else {
			Error(w, r, http.StatusBadRequest, "invalid bill type")
			return
		}

		if b.BillNumber, _ = parseIntField(r, "billNumber", log, w); b.BillNumber <= 0 {
			return
		}

		if b.Amount, _ = parseIntField(r, "amount", log, w); b.Amount < 0 {
			return
		}

		if b.DueDate, _ = parseDateField(r, "dueDate", dateLayout, log, w); b.DueDate.IsZero() {
			return
		}

		apartmentID := r.FormValue("apartmentID")
		if apartmentID == "" {
			Error(w, r, http.StatusBadRequest, "apartmentID is required")
			return
		}
		var aptID common.ID
		if err := aptID.UnmarshalText([]byte(apartmentID)); err != nil {
			log.Error("invalid apartment id", zap.Error(err))
			Error(w, r, http.StatusBadRequest, "invalid apartment id")
			return
		}
		b.ApartmentID = aptID

		// Handle optional image
		if err := handleImageUpload(r, &b, log, w); err != nil {
			return
		}

		// Service Call
		svc := svcGetter(r.Context())
		newBill, err := svc.AddBill(r.Context(), &b)
		if err != nil {
			log.Error("failed to create bill", zap.Error(err))
			switch {
			case errors.Is(err, bill.ErrBillOnValidate):
				BadRequestError(w, r, err.Error())
			case errors.Is(err, bill.ErrAlreadyExists):
				BadRequestError(w, r, err.Error())
			default:
				InternalServerError(w, r)
			}
			return
		}

		WriteJson(w, http.StatusCreated, map[string]any{
			"message": "bill created successfully",
			"id":      newBill.ID,
		})
	})
}

func parseIntField(r *http.Request, field string, log *zap.Logger, w http.ResponseWriter) (int64, bool) {
	valStr := r.FormValue(field)
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		log.Error("invalid int field", zap.String("field", field), zap.String("value", valStr), zap.Error(err))
		Error(w, r, http.StatusBadRequest, fmt.Sprintf("invalid %s", field))
		return 0, false
	}
	return val, true
}

func parseDateField(r *http.Request, field, layout string, log *zap.Logger, w http.ResponseWriter) (time.Time, bool) {
	valStr := r.FormValue(field)
	val, err := time.Parse(layout, valStr)
	if err != nil {
		log.Error("invalid date field", zap.String("field", field), zap.String("value", valStr), zap.Error(err))
		Error(w, r, http.StatusBadRequest, fmt.Sprintf("invalid %s format (expected YYYY-MM-DD)", field))
		return time.Time{}, false
	}
	return val, true
}

func handleImageUpload(r *http.Request, b *domain.Bill, log *zap.Logger, w http.ResponseWriter) error {
	file, header, err := r.FormFile("image")
	if err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			log.Error("failed to get image", zap.Error(err))
			InternalServerError(w, r)
			return err
		}
		return nil // image is optional
	}
	defer file.Close()

	if header.Size > maxFileSize {
		log.Error("uploaded file too large", zap.Int64("size", header.Size))
		Error(w, r, http.StatusRequestEntityTooLarge, "uploaded file is too large")
		return errors.New("file too large")
	}

	content, err := io.ReadAll(file)
	if err != nil {
		log.Error("failed to read file", zap.Error(err))
		InternalServerError(w, r)
		return err
	}

	contentType := http.DetectContentType(content)
	if !strings.HasPrefix(contentType, "image/") {
		Error(w, r, http.StatusBadRequest, "uploaded file is not an image")
		return errors.New("invalid image type")
	}

	dir, err := os.MkdirTemp(os.TempDir(), "*")
	if err != nil {
		log.Error("failed to make temp dir", zap.Error(err))
		return err
	}

	path := filepath.Join(dir, header.Filename)
	err = os.WriteFile(path, content, 0644)
	if err != nil {
		log.Error("failed to save image", zap.Error(err))
		return err
	}

	b.HasImage = true
	b.Image = &domain.Image{
		Name:    header.Filename,
		Path:    path,
		Type:    contentType,
		Size:    header.Size,
		Content: content,
	}
	return nil
}

// GetBill
//
// @Summary      Get bill details
// @Description  Returns details of a bill by ID
// @Tags         Bill
// @Accept       json
// @Produce      json
// @Security 	 BearerAuth
// @Param        body  body      dto.GetBillRequest  true  "Bill ID"
// @Success      200   {object}  domain.Bill
// @Failure      400   {object}  dto.Error
// @Failure      404   {object}  dto.Error
// @Failure      500   {object}  dto.Error
// @Router       /api/v1/bill [get]
func GetBill(svcGetter ServiceGetter[billPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())

		var req dto.GetBillRequest
		if err := BodyParse(r, &req); err != nil {
			log.Error("GetBill", zap.Error(err))
			BadRequestError(w, r)
			return
		}

		svc := svcGetter(r.Context())
		b, err := svc.GetBill(r.Context(), &domain.BillFilter{ID: req.ID})
		if err != nil {
			log.Error("GetBill", zap.Error(err))
			if errors.Is(err, bill.ErrNotFound) {
				Error(w, r, http.StatusNotFound, "bill not found")
				return
			}
			InternalServerError(w, r)
			return
		}

		if err := WriteJson(w, http.StatusOK, b); err != nil {
			log.Error("failed to write response", zap.Error(err))
			InternalServerError(w, r)
		}
	})
}

// GetBillImage
//
// @Summary      Get bill image
// @Description  Returns the image file for a bill
// @Tags         Bill
// @Accept       json
// @Produce      image/png
// @Security 	 BearerAuth
// @Param        body  body      dto.GetBillImageRequest  true  "Image ID"
// @Success      200   {file}    file
// @Failure      400   {object}  dto.Error
// @Failure      500   {object}  dto.Error
// @Router       /api/v1/bill/image [get]
func GetBillImage(svcGetter ServiceGetter[billPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())
		log.Debug("GetBillImage")

		var req dto.GetBillImageRequest
		if err := BodyParse(r, &req); err != nil {
			log.Error("GetBillImage", zap.Error(err))
			BadRequestError(w, r, "invalid image id")
		}

		svc := svcGetter(r.Context())
		path, err := svc.GetBillImage(r.Context(), req.ImageID)
		if err != nil {
			log.Error("GetBillImage", zap.Error(err))
			InternalServerError(w, r)
			return
		}

		file, err := os.Open(path)
		if err != nil {
			log.Error("GetBillImage", zap.Error(err))
			InternalServerError(w, r)
			return
		}
		defer file.Close()

		stat, err := os.Stat(path)
		if err == nil {
			log.Debug("", zap.Any("file stat", stat))
		}

		// Optionally detect content-type from file extension or contents
		w.Header().Set("Content-Type", "image/png") // or image/jpeg, etc.
		w.Header().Set("Content-Disposition", "inline; filename=\""+filepath.Base(path)+"\"")

		// Serve file content
		if _, err := io.Copy(w, file); err != nil {
			log.Error("GetBillImage - writing response", zap.Error(err))
		}

	})
}

// GetUserTotalDept
//
// @Summary      Get user's total debt
// @Description  Returns the total debt for the authenticated user
// @Tags         Bill
// @Produce      json
// @Security 	 BearerAuth
// @Success      200   {object}  dto.UserTotalDebt
// @Failure      500   {object}  dto.Error
// @Router       /api/v1/user/total-debt [get]
func GetUserTotalDept(svcGetter ServiceGetter[billPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())
		svc := svcGetter(r.Context())

		userID, ok := r.Context().Value(appjwt.UserIDKey).(string)
		if !ok {
			log.Error("invalid userid type")
			InternalServerError(w, r)
			return
		}

		id := common.IDFromText(userID)
		dbt, err := svc.GetUserTotalDebt(r.Context(), id)
		if err != nil {
			log.Error("", zap.Error(err))
			InternalServerError(w, r)
			return
		}

		resp := &dto.UserTotalDebt{TotalDebt: dbt}
		if err = WriteJson(w, http.StatusOK, resp); err != nil {
			log.Error("", zap.Error(err))
			InternalServerError(w, r)
		}
	})
}

// GetUserBillShares
//
// @Summary      Get user's bill shares
// @Description  Returns the bill shares for the authenticated user
// @Tags         Bill
// @Produce      json
// @Security 	 BearerAuth
// @Success      200   {object}  dto.BillSharesResponse
// @Failure      500   {object}  dto.Error
// @Router       /api/v1/user/bill-shares [get]
func GetUserBillShares(svcGetter ServiceGetter[billPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())
		svc := svcGetter(r.Context())

		userID, ok := r.Context().Value(appjwt.UserIDKey).(string)
		if !ok {
			log.Error("invalid userid type")
			InternalServerError(w, r)
			return
		}

		id := common.IDFromText(userID)
		billShares, err := svc.GetUserBillShares(r.Context(), id)
		if err != nil {
			log.Error("", zap.Error(err))
			switch {
			case errors.Is(err, bill.ErrNotFound):
				Error(w, r, http.StatusNotFound)
			default:
				InternalServerError(w, r)
			}
			return
		}

		resp := dto.BillSharesResponse{BillShares: billShares}
		if err = WriteJson(w, http.StatusOK, resp); err != nil {
			log.Error("", zap.Error(err))
			InternalServerError(w, r)
		}
	})
}
