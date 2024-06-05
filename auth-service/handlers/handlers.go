package handlers

import (
	"auth/pb"
	"auth/pkg"
	"context"
	"log"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
}

func (s *Server) ValidateToken(ctx context.Context, in *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	log.Printf("[Auth Service] ValidateToken request received for token: %s", in.Token)
	email, err := pkg.VerifyToken(in.Token)
	if err != nil {
		log.Println("[Auth Service] Token validation failed")
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}
	log.Println("[Auth Service] Token validation succeeded")
	return &pb.ValidateTokenResponse{Valid: true, Email: email}, nil
}

func (s *Server) GetTokenAndHash(ctx context.Context, in *pb.GetTokenHashRequest) (*pb.GetTokenHashResponse, error) {
	log.Printf("[Auth Service] Get token and hash for: %s", in.Email)
	token, err := pkg.CreateToken(in.Email)
	if err != nil {
		log.Println("[Auth Service] Token getting failed")
		return &pb.GetTokenHashResponse{Pass: "", Token: ""}, nil
	}
	newPass, err := pkg.HashPassword(in.Password)
	if err != nil {
		log.Println("[Auth Service] hashing password failed")
		return &pb.GetTokenHashResponse{Pass: "", Token: ""}, nil
	}
	log.Println("[Auth Service] Token validation succeeded")
	return &pb.GetTokenHashResponse{Token: token, Pass: newPass}, nil
}
