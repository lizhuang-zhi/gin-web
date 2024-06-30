package manager

import (
	pb "booking-app/micro-service/protobuf/gen-pb"
	"context"
	"sync"
)

type UserServcie struct {
	pb.UnimplementedUserServiceServer
	users map[string]*pb.User
	mu    sync.RWMutex
}

func NewUserService() *UserServcie {
	return &UserServcie{
		users: make(map[string]*pb.User, 0),
	}
}

func (s *UserServcie) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exist := s.users[req.User.Id]; exist {
		return &pb.CreateUserResponse{
			Success: false,
			Message: "user already exist",
		}, nil
	}

	user := &pb.User{
		Id:    req.User.Id,
		Name:  req.User.Name,
		Email: req.User.Email,
	}

	s.users[req.User.Id] = user

	return &pb.CreateUserResponse{
		Success: true,
		Message: "user created",
		User:    user,
	}, nil
}

func (s *UserServcie) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if req.Id == "" {
		return &pb.GetUserResponse{
			Success: false,
			Message: "user id is required",
		}, nil
	}
	user, exist := s.users[req.Id]
	if !exist {
		return &pb.GetUserResponse{
			Success: false,
			Message: "user not found",
		}, nil
	}

	return &pb.GetUserResponse{
		Success: true,
		Message: "user found",
		User:    user,
	}, nil
}
