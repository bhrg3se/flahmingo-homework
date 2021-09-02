package server

import (
	"context"
	pb "github.com/bhrg3se/flahmingo-homework/services/auth/pb/proto"
	"github.com/bhrg3se/flahmingo-homework/services/auth/store"
	"github.com/bhrg3se/flahmingo-homework/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

var empty = &emptypb.Empty{}

// SignupWithPhoneNumber creates a user profile and begins the phone verification process by sending the otp
func (s Server) SignupWithPhoneNumber(ctx context.Context, request *pb.User) (*emptypb.Empty, error) {
	if request.PhoneNumber == "" {
		err := status.Error(codes.InvalidArgument, "phone number is empty")
		logrus.Error(err)
		return empty, err
	}

	user := store.User{
		ID:          uuid.New().String(),
		Name:        request.Name,
		IsVerified:  false,
		PhoneNumber: request.PhoneNumber,
	}
	err := s.store.CreateUser(&user)
	if err != nil {
		logrus.Error(err)
		return empty, status.Error(codes.Internal, "could not create user")
	}

	otp := utils.GetRandomOTP()
	err = s.store.SaveOTP(otp, request.PhoneNumber)
	if err != nil {
		logrus.Error(err)
		return nil, status.Error(codes.Internal, "could not save otp")
	}

	// publish the otp on pubsub
	s.store.PublishOTP(ctx, otp, request.PhoneNumber)
	return empty, nil
}

// VerifyPhoneNumber takes otp entered by client and checks in database to verify it.
// If everything is good, user is marked as verified
func (s Server) VerifyPhoneNumber(ctx context.Context, request *pb.VerifyPhoneNumberRequest) (*emptypb.Empty, error) {
	otp, err := s.store.GetOTP(request.PhoneNumber)
	if err != nil {
		logrus.Error(err)
		return empty, status.Error(codes.Internal, "could not get otp")
	}

	if request.Otp != otp {
		logrus.Trace(request.Otp, otp)
		return empty, status.Error(codes.Unauthenticated, "invalid otp")
	}

	err = s.store.VerifyUser(request.PhoneNumber)
	if err != nil {
		logrus.Error(err)
		return empty, status.Error(codes.Internal, "could not verify user")
	}

	return empty, nil
}

func (s Server) LoginWithPhoneNumber(ctx context.Context, request *pb.User) (*emptypb.Empty, error) {

	_, err := s.store.GetUser(request.PhoneNumber)
	if err != nil {
		logrus.Debug(err)
		return empty, status.Error(codes.InvalidArgument, "phone number nor registered")
	}

	otp := utils.GetRandomOTP()
	err = s.store.SaveOTP(otp, request.PhoneNumber)
	if err != nil {
		logrus.Error(err)
		return nil, status.Error(codes.Internal, "could not save otp")
	}

	// publish the otp on pubsub
	s.store.PublishOTP(ctx, otp, request.PhoneNumber)
	return empty, nil
}

// ValidatePhoneNumberLogin takes token from client,verifies it and then creates a jwt auth token and returns it.
func (s Server) ValidatePhoneNumberLogin(ctx context.Context, request *pb.VerifyPhoneNumberRequest) (*pb.Token, error) {
	otp, err := s.store.GetOTP(request.PhoneNumber)
	if err != nil {
		logrus.Error(err)
		return nil, status.Error(codes.Internal, "could not get otp")
	}

	if request.Otp != otp {
		return nil, status.Error(codes.Unauthenticated, "invalid otp")
	}

	token, err := generateAuthToken(request.PhoneNumber, s.store.GetJWTPrivateKey())
	if err != nil {
		logrus.Error(err)
		return nil, status.Error(codes.Internal, "could not generate token")
	}
	return &pb.Token{Token: token}, nil
}

// GetProfile return profile of user based on auth token if the given token is valid
func (s Server) GetProfile(ctx context.Context, e *emptypb.Empty) (*pb.User, error) {
	// get auth token from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "could not find auth token")
	}
	authToken := md.Get("token")
	if len(authToken) < 1 {
		return nil, status.Error(codes.InvalidArgument, "could not find auth token")
	}

	// parse auth token to get phone number
	token, err := parseAuthToken(authToken[0], s.store.GetJWTPublicKey())
	if err != nil {
		logrus.Error(err)
		return nil, status.Error(codes.InvalidArgument, "could not parse auth token")
	}

	// check if the token is expired
	if !token.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, status.Error(codes.Unauthenticated, "auth token expired")
	}

	// get user profile from database
	user, err := s.store.GetUser(token.PhoneNumber)
	if err != nil {
		logrus.Error(err)
		return nil, status.Error(codes.Internal, "could not fetch user")
	}

	userPb := pb.User{
		Id:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
	}

	if !user.IsVerified {
		//TODO notify user that profile is not verified
		//return &userPb, status.Error(codes.PermissionDenied,"user not verified")
		return &userPb, nil
	}

	return &userPb, nil
}
