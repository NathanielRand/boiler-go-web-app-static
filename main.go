package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NathanielRand/boiler-static-web-app/views"
	"github.com/gorilla/mux"
)

const port = ":8080"

// Global variables related to templates.
// ATTN: Relocate for production. Global variables are bad practice.
var (
	homeView    *views.View
	contactView *views.View
	faqView     *views.View
)

var notfound http.Handler = http.HandlerFunc(notFound404)

type Visitor struct {
	IPAddress string
}

func notFound404(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>404 Not Foundish</h1>")
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	visitor := Visitor{
		IPAddress: GetIP(r),
	}

	must(homeView.Render(w, visitor))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(faqView.Render(w, nil))
}

// must is a helper for errors
func must(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func main() {
	homeView = views.NewView("bulma", "views/pages/home.html")
	contactView = views.NewView("bulma", "views/pages/contact.html")
	faqView = views.NewView("bulma", "views/pages/faq.html")

	// Gorilla Mux router
	r := mux.NewRouter()

	// Handle 404s
	r.NotFoundHandler = notfound

	// Server static assests (css, js, images, etc..)
	// Assest Routes
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	// Routes
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/faq", faq)

	// Start web server.
	fmt.Println("Listening on port", port)
	http.ListenAndServe(port, r)
}
