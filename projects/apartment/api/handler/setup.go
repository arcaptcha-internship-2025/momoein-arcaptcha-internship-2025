package handler

import (
	"fmt"
	"net/http"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/handler/middleware"
	"github.com/arcaptcha-internship-2025/momoein-apartment/api/handler/router"
	"github.com/arcaptcha-internship-2025/momoein-apartment/app"
	_ "github.com/arcaptcha-internship-2025/momoein-apartment/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Run(app app.App) error {
	r := router.NewRouter()
	RegisterAPI(r, app)

	addr := fmt.Sprintf(":%d", app.Config().HTTP.Port)
	app.Logger().Info("listen on " + addr)
	return http.ListenAndServe(addr, r)
}

// RegisterAPI
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func RegisterAPI(r *router.Router, app app.App) {
	jwtSecret := []byte(app.Config().Auth.JWTSecret)

	usrSvcGtr := UserServiceGetter(app)
	bilSvcGtr := BillServiceGetter(app)
	aptSvcGtr := ApartmentServiceGetter(app)
	paySvcGtr := PaymentServiceGetter(app)

	r.Use(
		middleware.SetRequestContext(app),
		middleware.LogRequest(),
	)
	r.Get("/", getRootHandler())

	r.Group("/api/v1", func(r *router.Router) {

		r.Group("/docs", func(r *router.Router) {
			r.Get("/swagger/", httpSwagger.Handler())
		})

		r.Group("/auth", func(r *router.Router) {
			r.Post("/sign-up", getSignUpHandler(usrSvcGtr, app.Config().Auth))
			r.Get("/sign-in", getSignInHandler(usrSvcGtr, app.Config().Auth))
			r.Get("/refresh-token", RefreshTokenHandler(usrSvcGtr, app.Config().Auth))
		})

		r.Group("/apartment", func(r *router.Router) {
			r.Use(middleware.NewAuth(jwtSecret))
			acceptURL := app.Config().BaseURL + "/api/v1/apartment/invite/accept"

			r.Post("/", AddApartment(aptSvcGtr))
			r.Post("/invite", InviteApartmentMember(aptSvcGtr, acceptURL))
			r.Get("/invite/accept", AcceptApartmentInvite(aptSvcGtr))
		})

		r.Group("/bill", func(r *router.Router) {
			r.Use(middleware.NewAuth(jwtSecret))

			r.Post("/", AddBill(bilSvcGtr))
			r.Get("/", GetBill(bilSvcGtr))
			r.Get("/image", GetBillImage(bilSvcGtr))
		})

		r.Group("/user", func(r *router.Router) {
			r.Use(middleware.NewAuth(jwtSecret))

			r.Get("/total-debt", GetUserTotalDept(bilSvcGtr))
			r.Get("/bill-shares", GetUserBillShares(bilSvcGtr))
		})

		r.Group("/payment", func(r *router.Router) {
			callbackURL := app.Config().BaseURL + "/api/v1/payment/callback"
			chain := router.Chain{middleware.NewAuth(jwtSecret)}

			r.Post("/pay-bill", chain.Then(PayUserBill(paySvcGtr, callbackURL)))
			r.Post("/pay-total-debt", chain.Then(PayTotalDebt(paySvcGtr, callbackURL)))
			r.Post("/callback", CallbackHandler(paySvcGtr))
			r.Get("/supported-gateways", SupportedGateways(paySvcGtr))

			r.Group("/mock-gateway", func(r *router.Router) {
				r.Post("/pay", MockGatewayPay())
				r.Get("/verify", MockGatewayVerify())
			})
		})
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
