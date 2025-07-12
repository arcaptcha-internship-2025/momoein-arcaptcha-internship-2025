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

type SingUpRequest struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type SingUpResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Error        Error  `json:"error"`
}

type SingInRequest struct{}

type SingInResponse struct{}

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
	Error     Error     `json:"error"`
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
	Error Error `json:"error"`
}

type InviteUserToApartmentRequest struct{}
type InviteUserToApartmentResponse struct {
	Error Error `json:"error"`
}

type RemoveUserFromApartmentRequest struct{}
type RemoveUserFromApartmentResponse struct {
	Error Error `json:"error"`
}

type ListApartmentUsersRequest struct{}
type ListApartmentUsersResponse struct {
	Error Error `json:"error"`
}

type AddBillRequest struct{}
type AddBillResponse struct {
	Error Error `json:"error"`
}

type RemoveBillRequest struct{}
type RemoveBillResponse struct {
	Error Error `json:"error"`
}

type SendBillNotificationToUsersRequest struct{}
type SendBillNotificationToUsersResponse struct {
	Error Error `json:"error"`
}

type CalculateUserBillRequest struct{}
type CalculateUserBillResponse struct {
	Error Error `json:"error"`
}
