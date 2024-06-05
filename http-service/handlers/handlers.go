package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"http/models"
	"http/pb"
	"http/repository"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	authClient  pb.AuthServiceClient
	emailClient pb.EmailServiceClient
	repo        repository.Repository
}

func NewHandler(authClient pb.AuthServiceClient, emailClient pb.EmailServiceClient, repo repository.Repository) *Handler {
	return &Handler{authClient: authClient, emailClient: emailClient, repo: repo}
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *Handler) Login(writer http.ResponseWriter, request *http.Request) {
	log.Println("[Main Service] Login request received")

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, "Invalid request payload", http.StatusBadRequest)
		log.Println("[Main Service] Error decoding request payload:", err)
		return
	}
	authReq := &pb.GetTokenHashRequest{Email: req.Email, Password: req.Password}
	authRes, err := s.authClient.GetTokenAndHash(context.Background(), authReq)
	if err != nil || authRes.Token == "" {
		http.Error(writer, "Failed to authenticate user", http.StatusUnauthorized)
		log.Println("[Main Service] Failed to authenticate user:", err)
		return
	}
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(writer, "Failed to authenticate user", http.StatusUnauthorized)
		log.Println("[Main Service] Error getting user from database:", err)
		return
	}
	if !user.Is_verified {
		http.Error(writer, "User email not verified", http.StatusUnauthorized)
		log.Println("[Main Service] User email not verified")
		return
	}
	if CheckPasswordHash(user.Password, req.Password) {
		http.Error(writer, "Invalid credentials", http.StatusUnauthorized)
		log.Println("[Main Service] Invalid credentials")
		return
	}
	log.Println("[Main Service] User authenticated successfully")
	response := struct {
		Token string `json:"token"`
	}{
		Token: authRes.Token,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(writer, "Failed to marshal response", http.StatusInternalServerError)
		log.Println("[Main Service] Error marshaling response:", err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonResponse)
	log.Println("[Main Service] Login successful")
}

func (s *Handler) Register(writer http.ResponseWriter, request *http.Request) {
	log.Println("[Main Service] Register request received")
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, "Invalid request payload", http.StatusBadRequest)
		return
	}
	authReq := &pb.GetTokenHashRequest{Email: req.Email, Password: req.Password}
	authRes, err := s.authClient.GetTokenAndHash(context.Background(), authReq)
	if err != nil || authRes.Token == "" {
		log.Println("[Main Service] Failed to get token and hash from auth-service")
		http.Error(writer, "Failed to register user", http.StatusInternalServerError)
		return
	}
	emailReq := &pb.EmailRequest{
		To: req.Email,
	}
	emailRes, err := s.emailClient.SendVerificationCode(context.Background(), emailReq)
	if err != nil {
		log.Println("[Main Service] Failed to send email:", err)
		http.Error(writer, "Failed to send email", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Email:            req.Email,
		Password:         authRes.Pass,
		Is_verified:      false,
		VerificationCode: emailRes.Code,
	}
	if err := s.repo.CreateUser(&user); err != nil {
		log.Println("[Main Service] Failed to save user to database:", err)
		http.Error(writer, "Failed to register user", http.StatusInternalServerError)
		return
	}

	log.Println("[Main Service] User registered successfully")
	response := struct {
		Token string `json:"token"`
	}{
		Token: authRes.Token,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(writer, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writer.Write(jsonResponse)
}

func (s *Handler) GetBooks(writer http.ResponseWriter, request *http.Request) {
	filter := request.URL.Query().Get("filter")
	sort := request.URL.Query().Get("sort")
	limitStr := request.URL.Query().Get("limit")
	offsetStr := request.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // default limit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // default offset
	}

	books, err := s.repo.GetBooks(filter, sort, limit, offset)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(books)
}

func (s *Handler) CreateBook(writer http.ResponseWriter, request *http.Request) {
	var book models.Book
	if err := json.NewDecoder(request.Body).Decode(&book); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.repo.CreateBook(&book); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusCreated)
}

func (s *Handler) GetBook(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(writer, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := s.repo.GetBookByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(writer, "Book not found", http.StatusNotFound)
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(book)
}

func (s *Handler) UpdateBook(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(writer, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book models.Book
	if err := json.NewDecoder(request.Body).Decode(&book); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	book.ID = id

	if err := s.repo.UpdateBook(&book); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

func (s *Handler) DeleteBook(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(writer, "Invalid book ID", http.StatusBadRequest)
		return
	}

	if err := s.repo.DeleteBook(id); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

func (s *Handler) VerifyCode(writer http.ResponseWriter, request *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, "Invalid request payload", http.StatusBadRequest)
		log.Println("[VerifyCode] Error decoding request payload:", err)
		return
	}
	log.Printf("[VerifyCode] Verifying code for email: %s, code: %s\n", req.Email, req.Code)

	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(writer, "Failed to get user", http.StatusInternalServerError)
		log.Println("[VerifyCode] Error getting user:", err)
		return
	}
	log.Printf("[VerifyCode] User found: %+v\n", user)

	if user.VerificationCode != req.Code {
		http.Error(writer, "Invalid verification code", http.StatusBadRequest)
		log.Println("[VerifyCode] Invalid verification code")
		return
	}

	log.Printf("[VerifyCode] Verification code matched for user: %+v\n", user)

	if err := s.repo.UpdateUser(user); err != nil {
		http.Error(writer, "Failed to update user", http.StatusInternalServerError)
		log.Println("[VerifyCode] Error updating user:", err)
		return
	}
	log.Printf("[VerifyCode] User updated successfully")

	emailReq := &pb.EmailRequest{
		To: req.Email,
	}
	if _, err := s.emailClient.SendConfirmationEmail(context.Background(), emailReq); err != nil {
		http.Error(writer, "Failed to send email", http.StatusInternalServerError)
		log.Println("[VerifyCode] Error sending confirmation email:", err)
		return
	}
	log.Printf("[VerifyCode] Confirmation email sent successfully")

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Verification successful"))
	log.Println("[VerifyCode] Verification successful")
}
