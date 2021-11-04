package routes

import "net/http"

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	p := "./statics/indexpage.html"
	http.ServeFile(w, r, p)
}
func SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	p := "./statics/signuppage.html"
	http.ServeFile(w, r, p)
}
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	p := "./statics/loginpage.html"
	http.ServeFile(w, r, p)
}

func AccountPageHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	p := "./statics/accountpage.html"
	http.ServeFile(w, r, p)
}
func ProfilePageHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	p := "./statics/profilepage.html"
	http.ServeFile(w, r, p)
}
func LoanPageHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	p := "./statics/loanpage.html"
	http.ServeFile(w, r, p)
}
