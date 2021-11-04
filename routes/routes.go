package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gregoflash05/trove/controllers"
)

type Handler struct {
	Router *mux.Router
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) SetupRoutes() {
	h.Router = mux.NewRouter().StrictSlash(true)

	// frontend
	h.Router.HandleFunc("/", HomePageHandler)
	h.Router.HandleFunc("/signup", SignUpPageHandler)
	h.Router.HandleFunc("/login", LoginPageHandler)
	h.Router.HandleFunc("/account", AccountPageHandler)
	h.Router.HandleFunc("/profile", ProfilePageHandler)
	h.Router.HandleFunc("/loan", LoanPageHandler)

	// Backend
	h.Router.HandleFunc("/v1/register", controllers.UserCreate).Methods("POST")
	h.Router.HandleFunc("/v1/login", controllers.UserLogin).Methods("POST")
	h.Router.HandleFunc("/v1/verify-token", controllers.VerifyTokenHandler).Methods("POST")
	h.Router.HandleFunc("/v1/user", controllers.GetUser).Methods("GET")
	h.Router.HandleFunc("/v1/profile", controllers.UserUpdate).Methods("PUT")
	h.Router.HandleFunc("/v1/password", controllers.PasswordUpdate).Methods("PUT")
	h.Router.HandleFunc("/v1/loan", controllers.TakeLoan).Methods("POST")
	h.Router.HandleFunc("/v1/loan", controllers.GetLoan).Methods("GET")
	h.Router.HandleFunc("/v1/loan", controllers.PayBackLoan).Methods("PUT")

	// statics
	h.Router.PathPrefix("/statics/").Handler(http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics/"))))

}
