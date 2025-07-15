package dto

import "time"

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

type SignUpResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SignInRequest struct{}

type SignInResponse struct{}

type Apartment struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	BillType    string    `json:"billType"`
	BillId      uint      `json:"billId"`
	Amount      uint      `json:"amount"`
	DueDate     time.Time `json:"dueDate"`
	ImageId     string    `json:"imageId"`
	ApartmentId string    `json:"apartmentId"`
}

type AddApartmentRequest struct {
	Apartment Apartment `json:"apartment"`
}
type AddApartmentResponse struct {
	Apartment Apartment `json:"apartment"`
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

type InviteUserToApartmentRequest struct{}
type InviteUserToApartmentResponse struct {
}

type RemoveUserFromApartmentRequest struct{}
type RemoveUserFromApartmentResponse struct {
}

type ListApartmentUsersRequest struct{}
type ListApartmentUsersResponse struct {
}

type AddBillRequest struct{}
type AddBillResponse struct {
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
