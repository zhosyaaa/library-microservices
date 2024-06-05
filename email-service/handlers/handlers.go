package handlers

import (
	"context"
	"email/pb"
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"
)

type Server struct {
	pb.UnimplementedEmailServiceServer
}

func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(100000 + rand.Intn(900000))
}

func sendEmail(to, subject, body string) error {
	from := "musabecova05@gmail.com"
	password := ""

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	return smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
}

func (s *Server) SendVerificationCode(ctx context.Context, in *pb.EmailRequest) (*pb.EmailResponse, error) {
	log.Printf("[Email Service] Sending verification code to: %s", in.To)
	code := generateVerificationCode()
	subject := "Your Verification Code"
	body := "Your verification code is: " + code
	err := sendEmail(in.To, subject, body)
	if err != nil {
		log.Printf("[Email Service] Failed to send verification code: %v", err)
		return &pb.EmailResponse{Success: false}, err
	}
	log.Println("[Email Service] Verification code sent successfully")
	return &pb.EmailResponse{Success: true}, nil
}

func (s *Server) SendConfirmationEmail(ctx context.Context, in *pb.EmailRequest) (*pb.EmailResponse, error) {
	log.Printf("[Email Service] Sending confirmation email to: %s", in.To)
	subject := "Welcome!"
	body := "Thank you for verifying your email."
	err := sendEmail(in.To, subject, body)
	if err != nil {
		log.Printf("[Email Service] Failed to send confirmation email: %v", err)
		return &pb.EmailResponse{Success: false}, err
	}
	log.Println("[Email Service] Confirmation email sent successfully")
	return &pb.EmailResponse{Success: true}, nil
}
