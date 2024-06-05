package handlers

import (
	"context"
	"encoding/json"
	"http/models"
	"http/pb"
	"http/repository"
	"log"
	"net/http"
)

type Handler struct {
	authClient  pb.AuthServiceClient
	emailClient pb.EmailServiceClient
	repo        repository.Repository
}

func NewHandler(authClient pb.AuthServiceClient, emailClient pb.EmailServiceClient, repo repository.Repository) *Handler {
	return &Handler{authClient: authClient, emailClient: emailClient, repo: repo}
}

func (s *Handler) Login(writer http.ResponseWriter, request *http.Request) {

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

}

func (s *Handler) CreateBook(writer http.ResponseWriter, request *http.Request) {

}

func (s *Handler) GetBook(writer http.ResponseWriter, request *http.Request) {

}

func (s *Handler) UpdateBook(writer http.ResponseWriter, request *http.Request) {

}

func (s *Handler) DeleteBook(writer http.ResponseWriter, request *http.Request) {

}

func (s *Handler) Profile(writer http.ResponseWriter, request *http.Request) {

}

func (s *Handler) VerifyCode(writer http.ResponseWriter, request *http.Request) {

}
