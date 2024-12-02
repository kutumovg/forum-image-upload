package handlers

import (
	"html/template"
	"log"
	"net/http"

	"forum/models"
)

// MainPageHandler - Displays the main page with posts and user information if logged in
func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	var username string
	loggedIn := false
	var userID string

	// Check if the user is logged in
	cookie, err := r.Cookie("session_token")
	if err == nil {
		sessionToken := cookie.Value

		// Get the username of the logged-in user
		userID, username, err = models.GetIDBySessionToken(sessionToken)
		if err == nil {
			loggedIn = true // User is logged in
		}
	}

	// Get filters from query parameters
	categoryID := r.URL.Query().Get("category")

	// Retrieve all posts
	posts, err := models.GetFilteredPosts(loggedIn, userID, categoryID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching posts")
		return
	}

	// Retrieve all categories
	categories, err := models.GetAllCategories()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching categories")
		return
	}

	// Check if there is a notification query parameter
	notification := r.URL.Query().Get("notification")

	// Load the index.html template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error loading template")
		return
	}
	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	data := struct {
		Posts            []models.Post
		Categories       []models.Category
		LoggedIn         bool
		Username         string
		Notification     string
		SelectedCategory string
	}{
		Posts:            posts,
		Categories:       categories,
		LoggedIn:         loggedIn,
		Username:         username,
		Notification:     notification,
		SelectedCategory: categoryID,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}
