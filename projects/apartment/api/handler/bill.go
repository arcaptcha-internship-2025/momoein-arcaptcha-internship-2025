package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/dto"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/domain"
	billPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
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

		if status := domain.PaymentStatus(r.FormValue("status")); status.IsValid() {
			b.Status = status
		} else {
			Error(w, r, http.StatusBadRequest, "invalid payment status")
			return
		}

		if paidAtStr := r.FormValue("paidAt"); paidAtStr != "" {
			if b.PaidAt, _ = parseDateField(r, "paidAt", dateLayout, log, w); b.PaidAt.IsZero() {
				return
			}
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
			if errors.Is(err, bill.ErrBillOnValidate) {
				BadRequestError(w, r, err.Error())
				return
			}
			InternalServerError(w, r)
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

	b.HasImage = true
	b.Image = domain.Image{
		Name:    header.Filename,
		Type:    contentType,
		Size:    header.Size,
		Content: content,
	}
	return nil
}

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

func GetBillImage(svcGetter ServiceGetter[billPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())
		log.Debug("GetBillImage")

		var req dto.GetBillImageRequest
		if err := BodyParse(r, &req); err != nil {
			log.Error("GetBillImage", zap.Error(err))
			BadRequestError(w, r, "invalid image id")
		}

		log.Debug("GetBillImage", zap.String("imageId", req.ImageID.String()))
	})
}
