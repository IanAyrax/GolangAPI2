package service

import(
	"context"
	"database/sql"
	"github.com/go-playground/validator/v10"
	"example.com/GolangAPI2/model"
	"example.com/GolangAPI2/repository"
	"example.com/GolangAPI2/helper"
	"example.com/GolangAPI2/exception"
	"fmt"
	"errors"
)

type UserServiceImpl struct {
	UserRepository 	repository.UserRepository
	DB				*sql.DB
	Validate		*validator.Validate
}

func NewUserService(userRepository repository.UserRepository, DB *sql.DB, validate *validator.Validate) UserService {
	return &UserServiceImpl {
		UserRepository:	userRepository,
		DB:				DB,
		Validate:		validate,
	}
}

func (service *UserServiceImpl) Create(ctx context.Context, request model.UserCreateRequest) model.UserResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user := model.User{
		FullName: request.FullName,
		Email:	request.Email,
		Password:	request.Password,
		RoleId:	request.RoleId,
	}

	user = service.UserRepository.Save(ctx, tx, user)

	return helper.ToUserResponse(user)
}

func (service *UserServiceImpl) Update(ctx context.Context, request model.UserUpdateRequest, roleId string, loggedId string) model.UserResponse{
	if helper.IsAdmin(roleId) != nil && loggedId == fmt.Sprintf("%v", request.Id){
		helper.PanicIfError(errors.New("Access Not Allowed !!!!"))
	}

	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	user.FullName = request.FullName
	user.Email = request.Email
	user.RoleId = request.RoleId
	user.Password = request.Password

	user = service.UserRepository.Update(ctx, tx, user)

	return helper.ToUserResponse(user)
}

func (service *UserServiceImpl) Delete(ctx context.Context, roleId string, loggedId string, userId int) {
	if helper.IsAdmin(roleId) != nil && loggedId == fmt.Sprintf("%v", userId){
		helper.PanicIfError(errors.New("Access Not Allowed !!!!"))
	}

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, userId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.UserRepository.Delete(ctx, tx, user)
}

func (service *UserServiceImpl) FindById(ctx context.Context, roleId string, loggedId string, userId int) model.UserResponse {
	if helper.IsAdmin(roleId) != nil && loggedId == fmt.Sprintf("%v", userId){
		helper.PanicIfError(errors.New("Access Not Allowed !!!!"))
	}

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, userId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	return helper.ToUserResponse(user)
}

func (service *UserServiceImpl) GetAll(ctx context.Context, roleId string) []model.UserResponse {
	err := helper.IsAdmin(roleId)
	helper.PanicIfError(err)

	fmt.Println("Service OK")
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	users := service.UserRepository.GetAll(ctx, tx)

	return helper.ToUserResponses(users)
}