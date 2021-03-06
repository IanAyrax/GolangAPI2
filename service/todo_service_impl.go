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

type ToDoServiceImpl struct {
	ToDoRepository 	repository.ToDoRepository
	DB				*sql.DB
	Validate		*validator.Validate
}

func NewToDoService(todoRepository repository.ToDoRepository, DB *sql.DB, validate *validator.Validate) ToDoService {
	return &ToDoServiceImpl {
		ToDoRepository:	todoRepository,
		DB:				DB,
		Validate:		validate,
	}
}

func (service *ToDoServiceImpl) Create(ctx context.Context, request model.ToDoCreateRequest) (model.ToDoResponse) {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	todo := model.ToDo{
		UserId: request.UserId,
		Title:	request.Title,
	}

	todo = service.ToDoRepository.Save(ctx, tx, todo)

	return helper.ToToDoResponse(todo)
}

func (service *ToDoServiceImpl) Update(ctx context.Context, get model.ToDoResponse, request model.ToDoUpdateRequest, roleId string, userId string) model.ToDoResponse{
	if fmt.Sprintf("%v", get.UserId) != userId && helper.IsAdmin(roleId) != nil{
		helper.PanicIfError(errors.New("Action Not Allowed : Not the Owner!!!!"))
	}
	
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	todo, err := service.ToDoRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	todo.UserId = request.UserId
	todo.Title = request.Title

	todo = service.ToDoRepository.Update(ctx, tx, todo)

	return helper.ToToDoResponse(todo)
}

func (service *ToDoServiceImpl) Delete(ctx context.Context, get model.ToDoResponse, roleId string, userId string, todoId int) {
	if fmt.Sprintf("%v", get.UserId) != userId && helper.IsAdmin(roleId) != nil{
		helper.PanicIfError(errors.New("Action Not Allowed : Not the Owner!!!!"))
	}

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	todo, err := service.ToDoRepository.FindById(ctx, tx, todoId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.ToDoRepository.Delete(ctx, tx, todo)
}

func (service *ToDoServiceImpl) FindById(ctx context.Context, userId string, roleId string, todoId int) model.ToDoResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	todo, err := service.ToDoRepository.FindById(ctx, tx, todoId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	
	toDoResponse := helper.ToToDoResponse(todo)
	if fmt.Sprintf("%v", toDoResponse.UserId) != userId && helper.IsAdmin(roleId) != nil{
		helper.PanicIfError(errors.New("Action Not Allowed : Not the Owner!!!!"))
	}

	return toDoResponse
}

func (service *ToDoServiceImpl) GetAll(ctx context.Context, roleId string) []model.ToDoResponse {
	if helper.IsAdmin(roleId) != nil {
		helper.PanicIfError(errors.New("Action Not Allowed : Not the Owner!!!!"))
	}

	fmt.Println("Service OK")
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	todos := service.ToDoRepository.GetAll(ctx, tx)

	return helper.ToToDoResponses(todos)
}