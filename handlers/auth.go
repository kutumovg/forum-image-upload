package handlers

import (
	"html/template"
	"net/http"
	"regexp"
	"time"

	"forum/models"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// authorization
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Validate email format
		if !isValidEmail(email) {
			ErrorHandler(w, r, http.StatusBadRequest, "Invalid email format")
			return
		}

		// Check if the email is already in use
		emailExists, err := models.CheckEmailExists(email)

		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		//if in use
		if emailExists {
			tmpl, _ := template.ParseFiles("templates/register.html")
			tmpl.Execute(w, struct{ Error string }{Error: "Email is already registered"})
			return
		}

		// Check if the username is already in use
		usernameExists, err := models.CheckUsernameExists(username)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		if usernameExists {
			tmpl, _ := template.ParseFiles("templates/register.html")
			tmpl.Execute(w, struct{ Error string }{Error: "Username is already taken"})
			return
		}

		// Register the user (create user in the database)
		sessionToken, err := models.RegisterUser(email, username, password)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		// Automatically log in the user after registration
		cookie := http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Now().Add(24 * time.Hour),
		}
		http.SetCookie(w, &cookie)

		// Redirect to the main page after successful registration and login
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, _ := template.ParseFiles("templates/register.html")
	tmpl.Execute(w, nil)
}

// LoginHandler - Handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Authenticate the user
		sessionToken, err := models.AuthenticateUser(email, password)
		if err != nil {
			tmpl, _ := template.ParseFiles("templates/login.html")
			tmpl.Execute(w, struct{ Error string }{Error: "Invalid email or password"})
			return
		}

		// Set session cookie
		cookie := http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Now().Add(24 * time.Hour),
		}
		http.SetCookie(w, &cookie)

		// Redirect to the main page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Render the login page
	tmpl, _ := template.ParseFiles("templates/login.html")
	tmpl.Execute(w, nil)
}

// LogoutHandler - Logs the user out by clearing the session cookie
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the session cookie
	cookie := http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour), // Expire the cookie immediately
	}
	http.SetCookie(w, &cookie)

	// Redirect to the main page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
