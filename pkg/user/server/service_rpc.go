package server

import (
	"UserGrpcProj/models"
	pb "UserGrpcProj/pkg/user/service"
	"UserGrpcProj/repository/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"regexp"
)

type ServiceUser struct {
	pb.UserServer

	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewUserService(db *pgxpool.Pool, logger *zap.Logger) *ServiceUser {
	return &ServiceUser{
		db:     db,
		logger: logger,
	}
}
func (s *ServiceUser) AddUser(ctx context.Context, request *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	login := request.GetLogin()

	templ := "^[a-zA-Z0-9]{5,}$"
	match, err := regexp.MatchString(templ, login)
	if err != nil {
		return nil, err
	}
	if !match {
		err := errors.New("wrong login")
		return nil, err
	}
	userRepository := postgres.NewUserRepo(ctx, s.db, s.logger)

	if userRepository.IsSetUser(login, 0) {
		return nil, errors.New("login already exists")
	}

	password := request.GetPassword()
	if len([]rune(password)) < 5 {
		err := errors.New("the password is too short")
		return nil, err
	}
	name := request.GetName()
	if len([]rune(name)) < 3 {
		err := errors.New("the name is too short")
		return nil, err
	}
	phone := request.GetPhone()
	if !IsPhoneValid(phone) {
		err := errors.New("wrong phone")
		return nil, err
	}
	user := models.UserInfo{
		Login:    login,
		Password: password,
		Name:     name,
		Phone:    phone,
	}

	ok, err := userRepository.AddUser(&user)
	status := &pb.AddUserResponse{
		Status: ok,
	}
	if err != nil {
		fmt.Println(err)
		return status, err
	}
	return status, nil
}
func (s *ServiceUser) RemoveUser(ctx context.Context, request *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	removeId := request.GetId()
	if removeId <= 0 {
		return nil, errors.New("wrong request")
	}
	userRepository := postgres.NewUserRepo(ctx, s.db, s.logger)
	if userRepository.IsSetUser("", int(removeId)) == false {
		return nil, errors.New("id not found")
	}
	if removeId < 0 {
		return nil, errors.New("wrong id")
	}

	result, err := userRepository.RemoveUser(int(removeId))
	if err != nil {
		return nil, err
	}
	status := &pb.RemoveUserResponse{
		Status: result,
	}
	return status, nil
}

func (s *ServiceUser) UserList(ctx context.Context, request *pb.UserListRequest) (*pb.UserListResponse, error) {
	filter := request.GetFilter()
	phone := filter.Phone
	name := filter.Name
	login := filter.Login

	userRepository := postgres.NewUserRepo(ctx, s.db, s.logger)
	resp, err := userRepository.UserList(login, name, phone)
	if err != nil {
		return nil, err
	}
	result := &pb.UserListResponse{}

	for _, respUser := range resp.User {
		result.UserList = append(result.UserList, &pb.UserInfo{
			Id:    int32(respUser.Id),
			Login: respUser.Login,
			Name:  respUser.Name,
			Phone: respUser.Phone,
		})
	}

	return result, nil
}
func IsPhoneValid(p string) bool {
	phoneRegex := regexp.MustCompile("^((8|\\+7)[\\- ]?)?(\\(?\\d{3}\\)?[\\- ]?)?[\\d\\- ]{7,10}$")
	if !phoneRegex.MatchString(p) || p == "" {
		return false
	}

	return true
}
