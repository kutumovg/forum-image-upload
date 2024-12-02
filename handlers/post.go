package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"forum/models"
)

// Handler for creating a post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	// Check if the user is logged in
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	content := r.FormValue("content")
	categories := r.Form["categories"]

	content = models.SanitizeInput(content)
	if !models.IsValidContent(content) || len(categories) == 0 {
		ErrorHandler(w, r, http.StatusBadRequest, "Content and at least one category are required to create a post")
		return
	}

	var imagePath string
	if file, header, err := r.FormFile("image"); err == nil {
		defer file.Close()

		// Validate the image
		if err := validateImage(file, header); err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, err.Error())
			return
		}

		// Save the image
		imagePath, err = saveImage(file, header)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}

	postID, err := models.CreatePost(userID, content, imagePath)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error creating post")
		return
	}

	for _, categoryID := range categories {
		err = models.AddCategoryToPost(postID, categoryID)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, "Error associating category")
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Handler for liking a post
func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	// Check if the user is logged in
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	postID := r.FormValue("post_id")

	// Like the post
	err = models.LikePost(userID, postID)
	if err != nil {
		// if err.Error() == "you have already liked this post" {
		// 	// Redirect back to the main page with a notification
		// 	http.Redirect(w, r, "/?notification=already_liked", http.StatusSeeOther)
		// 	return
		// }

		http.Error(w, "Error liking post: "+err.Error(), http.StatusInternalServerError)
		ErrorHandler(w, r, http.StatusInternalServerError, "Error liking post")
		return
	}

	// Update the post's like and dislike counts
	err = models.UpdatePostLikesDislikes(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error updating like count")
		return
	}

	// Redirect back to the main page
	// http.Redirect(w, r, "/", http.StatusSeeOther)
	referer := r.Header.Get("Referer")
	http.Redirect(w, r, referer, http.StatusSeeOther)
	// if strings.Contains(referer, "/post") {
	// 	// If the referer is a post page, redirect back to that page
	// 	http.Redirect(w, r, referer, http.StatusSeeOther)
	// } else {
	// 	// Otherwise, redirect to the homepage
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// }
}

// Handler for disliking a post
func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	// Check if the user is logged in
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	postID := r.FormValue("post_id")

	// Dislike the post
	err = models.DislikePost(userID, postID)
	if err != nil {
		// if err.Error() == "you have already disliked this post" {
		// 	// Redirect back to the main page with a notification
		// 	http.Redirect(w, r, "/?notification=already_disliked", http.StatusSeeOther)
		// 	return
		// }
		ErrorHandler(w, r, http.StatusInternalServerError, "Error disliking post")
		return
	}

	// Update the post's like and dislike counts
	err = models.UpdatePostLikesDislikes(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error updating dislike count")
		return
	}

	// Redirect back to the main page
	// http.Redirect(w, r, "/", http.StatusSeeOther)
	referer := r.Header.Get("Referer")
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func PostPageHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure we're handling a GET request
	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	// Get the post ID from the query string
	postID := r.URL.Query().Get("id")
	if postID == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "Missing post ID")
		return
	}

	// Fetch the post by ID
	post, err := models.GetPostByID(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			ErrorHandler(w, r, http.StatusNotFound, "Post not found")
			return
		}
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching post")
		return
	}

	// Fetch comments for the post
	comments, err := models.GetCommentsForPost(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching comments")
		return
	}

	notification := r.URL.Query().Get("notification")

	// Load the comments.html template
	tmpl, err := template.ParseFiles("templates/comments.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	var loggedIn bool
	var username string
	cookie, err := r.Cookie("session_token")
	if err == nil {
		loggedIn = true
		_, username, _ = models.GetIDBySessionToken(cookie.Value)
	}

	data := struct {
		Post         models.Post
		Comments     []models.Comment
		LoggedIn     bool
		Username     string
		Notification string
	}{
		Post:         post,
		Comments:     comments,
		LoggedIn:     loggedIn,
		Username:     username,
		Notification: notification,
	}

	tmpl.Execute(w, data)
}

func MyPostsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	// Fetch posts created by the logged-in user
	posts, err := models.GetPostsByUser(userID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching posts")
		return
	}

	categories, err := models.GetAllCategories()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching categories")
		return
	}

	// Render the posts page with "My Posts"
	tmpl, err := template.ParseFiles("templates/posts.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := struct {
		Posts            []models.Post
		Categories       []models.Category
		LoggedIn         bool
		Username         string
		SelectedCategory string
		SelectedFilter   string
	}{
		Posts:            posts,
		Categories:       categories,
		LoggedIn:         true,
		SelectedCategory: "",
		SelectedFilter:   "",
	}

	tmpl.Execute(w, data)
}

func LikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	// Fetch posts liked by the logged-in user
	posts, err := models.GetLikedPostsByUser(userID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching liked posts")
		return
	}

	categories, err := models.GetAllCategories()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching categories")
		return
	}

	// Render the posts page with "Liked Posts"
	tmpl, err := template.ParseFiles("templates/posts.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := struct {
		Posts            []models.Post
		Categories       []models.Category
		LoggedIn         bool
		Username         string
		SelectedCategory string
		SelectedFilter   string
	}{
		Posts:            posts,
		Categories:       categories,
		LoggedIn:         true,
		SelectedCategory: "",
		SelectedFilter:   "",
	}

	tmpl.Execute(w, data)
}
