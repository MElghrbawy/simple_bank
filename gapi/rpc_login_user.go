package gapi

import (
	"context"
	"database/sql"

	db "github.com/MElghrbawy/simple_bank/db/sqlc"
	"github.com/MElghrbawy/simple_bank/pb"
	"github.com/MElghrbawy/simple_bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(c context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.store.GetUser(c, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "user not found")

		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	if err := util.CheckPasswordHash(req.Password, user.HashedPassword); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid username or password")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create access token")

	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create refresh token")
	}

	session, err := server.store.CreateSession(c, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "c.ClientIP()",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create session")
	}

	res := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiresAt),
	}

	return res, nil
}
