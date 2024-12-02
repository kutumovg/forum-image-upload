// handlers/handlers_test.go
package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Define MockUserRepository struct with mock fields
type MockUserRepository struct {
	CheckEmailExistsMock    func(email string) (bool, error)
	CheckUsernameExistsMock func(username string) (bool, error)
	RegisterUserMock        func(email, username, password string) (string, error)
	AuthenticateUserMock    func(email, password string) (string, error)
	CreatePostMock          func(userID, content string, categories []string) error
}

// Define the UserRepository interface with required methods
type UserRepository interface {
	CheckEmailExists(email string) (bool, error)
	CheckUsernameExists(username string) (bool, error)
	RegisterUser(email, username, password string) (string, error)
	AuthenticateUser(email, password string) (string, error)
	CreatePost(userID, content string, categories []string) error
}

// Declare userRepo as a global variable implementing UserRepository
var userRepo UserRepository

// Implement UserRepository methods on MockUserRepository

func (m *MockUserRepository) CheckEmailExists(email string) (bool, error) {
	return m.CheckEmailExistsMock(email)
}

func (m *MockUserRepository) CheckUsernameExists(username string) (bool, error) {
	return m.CheckUsernameExistsMock(username)
}

func (m *MockUserRepository) RegisterUser(email, username, password string) (string, error) {
	return m.RegisterUserMock(email, username, password)
}

func (m *MockUserRepository) AuthenticateUser(email, password string) (string, error) {
	return m.AuthenticateUserMock(email, password)
}

func (m *MockUserRepository) CreatePost(userID, content string, categories []string) error {
	return m.CreatePostMock(userID, content, categories)
}

func TestRegisterHandler(t *testing.T) {
	// Setup mock repository with specific behavior
	mockRepo := &MockUserRepository{
		CheckEmailExistsMock: func(email string) (bool, error) {
			if email == "existing@example.com" {
				return true, nil
			}
			return false, nil
		},
		CheckUsernameExistsMock: func(username string) (bool, error) {
			return false, nil
		},
		RegisterUserMock: func(email, username, password string) (string, error) {
			return "mock_user_id", nil
		},
	}

	// Assign the mock to the global userRepo
	userRepo = mockRepo

	req := httptest.NewRequest("POST", "/register", strings.NewReader("email=new@example.com&username=newuser&password=secret"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler_test)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %v; got %v", http.StatusSeeOther, rr.Code)
	}
}

func TestLoginHandler(t *testing.T) {
	mockRepo := &MockUserRepository{
		AuthenticateUserMock: func(email, password string) (string, error) {
			if email == "valid@example.com" && password == "correctpassword" {
				return "mock_user_id", nil
			}
			return "", errors.New("invalid credentials")
		},
	}
	userRepo = mockRepo

	// Case 1: Successful login
	req := httptest.NewRequest("POST", "/login", strings.NewReader("email=valid@example.com&password=correctpassword"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler_test)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %v; got %v", http.StatusSeeOther, rr.Code)
	}

	// Case 2: Invalid credentials
	req = httptest.NewRequest("POST", "/login", strings.NewReader("email=invalid@example.com&password=wrongpassword"))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest || !strings.Contains(rr.Body.String(), "Invalid email or password") {
		t.Errorf("expected 'Invalid email or password' message; got %v", rr.Body.String())
	}
}

func TestCreatePostHandler(t *testing.T) {
	mockRepo := &MockUserRepository{
		CreatePostMock: func(userID, content string, categories []string) error {
			return nil // Simulate successful post creation
		},
	}
	userRepo = mockRepo

	// Case 1: Successful post creation
	req := httptest.NewRequest("POST", "/create_post", strings.NewReader("content=This is a new post&categories=1"))
	req.AddCookie(&http.Cookie{Name: "session_token", Value: "mock_user_id"})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreatePostHandler_test)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %v; got %v", http.StatusSeeOther, rr.Code)
	}

	// Case 2: Unauthorized (no session token)
	req = httptest.NewRequest("POST", "/create_post", strings.NewReader("content=This is a new post&categories=1"))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized || !strings.Contains(rr.Body.String(), "Unauthorized: Please log in to create a post") {
		t.Errorf("expected 'Unauthorized' error; got %v", rr.Body.String())
	}

	// Case 3: Missing content or categories
	req = httptest.NewRequest("POST", "/create_post", strings.NewReader("content=&categories="))
	req.AddCookie(&http.Cookie{Name: "session_token", Value: "mock_user_id"})
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest || !strings.Contains(rr.Body.String(), "Content and at least one category are required to create a post") {
		t.Errorf("expected 'Content and at least one category are required' error; got %v", rr.Body.String())
	}
}

// Example RegisterHandler function using userRepo
func RegisterHandler_test(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Check if email already exists
	exists, err := userRepo.CheckEmailExists(email)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Email is already registered", http.StatusBadRequest)
		return
	}

	// Check if username already exists
	exists, err = userRepo.CheckUsernameExists(username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Username is already taken", http.StatusBadRequest)
		return
	}

	// Register the user
	_, err = userRepo.RegisterUser(email, username, password)
	if err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginHandler_test(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Authenticate the user
	userID, err := userRepo.AuthenticateUser(email, password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	// Set a session or cookie upon successful authentication
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    userID, // In a real application, this should be a secure token
		HttpOnly: true,
		Path:     "/",
		// Optionally set other cookie attributes such as Secure, SameSite, etc.
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CreatePostHandler_test(w http.ResponseWriter, r *http.Request) {
	// Check for session token (simple authentication check)
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		http.Error(w, "Unauthorized: Please log in to create a post", http.StatusUnauthorized)
		return
	}

	// Retrieve form values
	content := r.FormValue("content")
	categories := r.Form["categories"]

	// Validate post content
	if content == "" || len(categories) == 0 {
		http.Error(w, "Content and at least one category are required to create a post", http.StatusBadRequest)
		return
	}

	// Simulate saving the post by calling a repository method (e.g., userRepo.CreatePost)
	err = userRepo.CreatePost(sessionCookie.Value, content, categories)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// Redirect to the main page or post page after successful creation
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
