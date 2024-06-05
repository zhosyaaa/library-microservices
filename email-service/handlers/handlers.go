package handlers

import (
	"context"
	"email/pb"
	"log"
	"math/rand"
	"net/smtp"
	"time"
)

type Server struct {
	pb.UnimplementedEmailServiceServer
}

const verificationCodeLength = 6

func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	chars := "0123456789"
	code := make([]byte, verificationCodeLength)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}

func sendEmail(to, subject, body string) error {
	from := "musabecova05@gmail.com"
	password := "mayf ayum loqn haqs"

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
	return &pb.EmailResponse{Success: true, Code: code}, nil
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
