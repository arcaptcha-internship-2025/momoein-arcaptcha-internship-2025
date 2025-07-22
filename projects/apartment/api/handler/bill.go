package handler

import (
	"errors"
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

func AddBill(svcGetter ServiceGetter[billPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())

		if err := r.ParseMultipartForm(1 * MiB); err != nil {
			log.Error("add bill", zap.Error(err))
			Error(w, r, http.StatusBadRequest, err.Error())
			return
		}

		log.Info("add bill", zap.Any("form values", r.MultipartForm.Value))

		var b domain.Bill

		// --- Name ---
		b.Name = r.FormValue("name")

		// --- Bill Type ---
		billTypeStr := r.FormValue("type")
		billType := domain.BillType(billTypeStr)
		if !billType.IsValid() {
			log.Error("AddBill", zap.String("error", "invalid bill type"))
			Error(w, r, http.StatusBadRequest, "invalid bill type")
			return
		}
		b.Type = billType

		// --- Bill Number ---
		billNumber := r.FormValue("billNumber")
		billNum, err := strconv.ParseInt(billNumber, 10, 64)
		if err != nil {
			log.Error("AddBill", zap.Error(err))
			Error(w, r, http.StatusBadRequest, "invalid bill number")
			return
		}
		b.BillNumber = billNum

		// --- Due Date ---
		dueDateStr := r.FormValue("dueDate")
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil {
			log.Error("AddBill", zap.String("dueDate", dueDateStr), zap.Error(err))
			Error(w, r, http.StatusBadRequest, "invalid due date format (expected YYYY-MM-DD)")
			return
		}
		b.DueDate = dueDate

		// --- Amount ---
		amountStr := r.FormValue("amount")
		amount, err := strconv.ParseInt(amountStr, 10, 64)
		if err != nil {
			log.Error("AddBill", zap.String("amount", amountStr), zap.Error(err))
			Error(w, r, http.StatusBadRequest, "invalid amount")
			return
		}
		b.Amount = amount

		// --- Status ---
		statusStr := r.FormValue("status")
		status := domain.PaymentStatus(statusStr)
		if !status.IsValid() {
			log.Error("AddBill", zap.String("status", statusStr))
			Error(w, r, http.StatusBadRequest, "invalid payment status")
			return
		}
		b.Status = status

		// --- Paid At (optional) ---
		paidAtStr := r.FormValue("paidAt")
		if paidAtStr != "" {
			paidAt, err := time.Parse("2006-01-02", paidAtStr)
			if err != nil {
				log.Error("AddBill", zap.String("paidAt", paidAtStr), zap.Error(err))
				Error(w, r, http.StatusBadRequest, "invalid paidAt date format (expected YYYY-MM-DD)")
				return
			}
			b.PaidAt = paidAt
		}

		// --- Apartment ID ---
		apartmentID := r.FormValue("apartmentID")
		if apartmentID == "" {
			Error(w, r, http.StatusBadRequest, "apartmentID is required")
			return
		}
		aptID := common.NilID
		if err := aptID.UnmarshalText([]byte(apartmentID)); err != nil {
			log.Error("AddBil", zap.Error(err))
			Error(w, r, http.StatusBadRequest, "invalid apartment id")
			return
		}
		b.ApartmentID = aptID

		// --- Optional Image ---
		file, header, err := r.FormFile("image")
		if err != nil && !errors.Is(err, http.ErrMissingFile) {
			log.Error("failed to get image", zap.Error(err))
			InternalServerError(w, r)
			return
		}

		if file != nil {
			defer file.Close()

			if header.Size > maxFileSize {
				log.Error("uploaded file too large", zap.Int64("size", header.Size))
				Error(w, r, http.StatusRequestEntityTooLarge, "uploaded file is too large")
				return
			}

			content, err := io.ReadAll(file)
			if err != nil {
				log.Error("failed to read file", zap.Error(err))
				InternalServerError(w, r)
				return
			}

			contentType := http.DetectContentType(content)
			if strings.HasPrefix(contentType, "image/") {
				b.HasImage = true
				b.Image.Name = header.Filename
				b.Image.Type = contentType
				b.Image.Size = header.Size
				b.Image.Content = content
			}
		}

		// --- Service Call ---
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

		// --- Success Response ---
		WriteJson(w, http.StatusCreated, map[string]any{
			"message": "bill created successfully",
			"id":      newBill.ID,
		})
	})
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
