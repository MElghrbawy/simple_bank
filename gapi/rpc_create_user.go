package gapi

import (
	"context"

	db "github.com/MElghrbawy/simple_bank/db/sqlc"
	"github.com/MElghrbawy/simple_bank/pb"
	"github.com/MElghrbawy/simple_bank/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(c context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(c, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Error(codes.AlreadyExists, "username or email already exists")
			}
		}
		return nil, status.Error(codes.Internal, "could not create user")
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil

}
