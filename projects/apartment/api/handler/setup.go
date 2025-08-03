package handler

import (
	"fmt"
	"net/http"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/handler/middleware"
	"github.com/arcaptcha-internship-2025/momoein-apartment/api/handler/router"
	"github.com/arcaptcha-internship-2025/momoein-apartment/app"
)

func Run(app app.App) error {
	r := router.NewRouter()
	r.Use(
		middleware.SetRequestContext(app),
		middleware.LogRequest(),
	)
	r.Get("/", getRootHandler())

	api := r.Group("/api/v1", nil)
	RegisterAPI(api, app)

	addr := fmt.Sprintf(":%d", app.Config().HTTP.Port)
	app.Logger().Info("listen on " + addr)
	return http.ListenAndServe(addr, r)
}

func RegisterAPI(r *router.Router, app app.App) {
	secret := []byte(app.Config().Auth.JWTSecret)

	usrSvcGetter := UserServiceGetter(app)
	bilSvcGtr := BillServiceGetter(app)
	aptSvcGetter := ApartmentServiceGetter(app)

	r.Group("/auth", func(r *router.Router) {
		r.Post("/sign-up", getSignUpHandler(usrSvcGetter, app.Config().Auth))
		r.Get("/sign-in", getSignInHandler(usrSvcGetter, app.Config().Auth))
		r.Get("/refresh-token", RefreshTokenHandler(usrSvcGetter, app.Config().Auth))
	})

	r.Group("/apartment", func(r *router.Router) {
		r.Use(middleware.NewAuth(secret))

		acceptURL := app.Config().BaseURL + "/api/v1/apartment/invite/accept"
		r.Post("/", AddApartment(aptSvcGetter))
		r.Post("/invite", InviteApartmentMember(aptSvcGetter, acceptURL))
		r.Get("/invite/accept", AcceptApartmentInvite(aptSvcGetter))

		r.Post("/bill", AddBill(bilSvcGtr))
		r.Get("/bill", GetBill(bilSvcGtr))
		r.Get("/bill/image", GetBillImage(bilSvcGtr))
	})

	r.Group("/user", func(r *router.Router) {
		r.Use(middleware.NewAuth(secret))
		r.Get("/total-debt", GetUserTotalDept(bilSvcGtr))
		r.Get("/bill-shares", GetUserBillShares(bilSvcGtr))
	})
}

func getRootHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte("hello from arcaptcha apartment api\n"))
	})
}
