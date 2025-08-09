package dto

import (
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type SignUpRequest struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
type Apartment struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	UnitNumber int64  `json:"unitNumber"`
	AdminID    string `json:"adminID"`
}

type ListUserApartmentsRequest struct{}
type ListUserApartmentsResponse struct {
	Apartments []Apartment `json:"apartments"`
	Error      Error       `json:"error"`
}

type RemoveApartmentRequest struct {
	Apartment Apartment
}
type RemoveApartmentResponse struct {
}

type InviteUserToApartmentRequest struct {
	UserEmail   common.Email `json:"userEmail"`
	ApartmentID common.ID    `json:"apartmentID"`
}

type UserTotalDebt struct {
	TotalDebt int `json:"totalDebt"`
}

type InviteUserToApartmentResponse struct {
}

type RemoveUserFromApartmentRequest struct{}
type RemoveUserFromApartmentResponse struct {
}

type ListApartmentUsersRequest struct{}
type ListApartmentUsersResponse struct {
}

type AddBillRequest struct {
	Name          string    `form:"name"`
	BillType      string    `form:"billType"`
	BillNumber    int64     `form:"billNumber"`
	DueDate       time.Time `form:"dueDate"`
	Amount        int64     `form:"amount"`
	PaymentStatus string    `form:"paymentStatus"`
	PaidAt        time.Time `form:"paidAt"`
}
type AddBillResponse struct {
}

type GetBillRequest struct {
	ID common.ID `json:"id" form:"id"`
}

type GetBillImageRequest struct {
	ImageID common.ID `json:"imageID"`
}

type BillSharesResponse struct {
	BillShares []domain.UserBillShare `json:"billShares"`
}

type RemoveBillRequest struct{}
type RemoveBillResponse struct {
}

type SendBillNotificationToUsersRequest struct{}
type SendBillNotificationToUsersResponse struct {
}

type CalculateUserBillRequest struct{}
type CalculateUserBillResponse struct {
}

type PayRequest struct {
	Amount      int64  `json:"amount"`
	CallbackURL string `json:"returnUrl"`
}

type PayResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type VerifyRequest struct {
	Token string `json:"token"`
}

type VerifyResponse struct {
	Code    int    `json:"code"`    //  ==0 success, !=0 failed
	Message string `json:"message"` // descriptive message
}

type PayBillRequest struct {
	BillID  string `json:"billID"`
	Gateway string `json:"gateway"`
}

type PayTotalDebtRequest struct {
	Gateway string `json:"gateway"`
}

type RedirectGateway struct {
	Method string         `json:"method"`
	URL    string         `json:"url"`
	Body   map[string]any `json:"body"`
}

type SupportedGatewaysResponse struct {
	SupportedGateways []string `json:"supportedGateways"`
}
