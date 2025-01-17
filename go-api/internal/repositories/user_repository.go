package repositories

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/Energie-Burgenland/ausaestung-info/internal/entities"
	"github.com/Energie-Burgenland/ausaestung-info/internal/models"
	"github.com/Energie-Burgenland/ausaestung-info/utils/auth"
	"github.com/Energie-Burgenland/ausaestung-info/utils/database"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type UserRepository struct {
	*BaseRepository[entities.User, models.UserModel, models.SaveUserModel]
}

func NewUserRepository(dbContext *database.DbContext) UserRepository {
	return UserRepository{
		BaseRepository: NewBaseRepository[entities.User, models.UserModel, models.SaveUserModel](dbContext),
	}
}

func (r *UserRepository) GetUsers(ctx context.Context, lastEvaluatedKey string, textQuery string) (models.ListResult[models.UserModel], error) {
	var filter expression.ConditionBuilder

	if textQuery != "" {
		orConditions := expression.Name("UserName").Contains(textQuery).
			Or(expression.Name("FirstName").Contains(textQuery)).
			Or(expression.Name("LastName").Contains(textQuery)).
			Or(expression.Name("Street").Contains(textQuery)).
			Or(expression.Name("City").Contains(textQuery)).
			Or(expression.Name("Zip").Contains(textQuery)).
			Or(expression.Name("Country").Contains(textQuery)).
			Or(expression.Name("Email").Contains(textQuery)).
			Or(expression.Name("PhoneNumber").Contains(textQuery))

		filter = filter.And(orConditions)
	}

	return r.BaseRepository.GetList(ctx, filter, lastEvaluatedKey)
}

func (r *UserRepository) GetUser(ctx context.Context, id string) (*models.UserModel, error) {
	return r.BaseRepository.Get(ctx, id)
}

func (r *UserRepository) CreateUser(ctx context.Context, saveModel *models.SaveUserModel) (*models.UserModel, error) {
	return r.BaseRepository.Create(ctx, *saveModel)
}

func (r *UserRepository) UpdateUser(ctx context.Context, id string, saveModel *models.SaveUserModel) (*models.UserModel, error) {
	model, err := r.BaseRepository.Update(ctx, id, *saveModel)
	if err != nil {
		return nil, err
	}

	return model, err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	err := r.BaseRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) ImportUsers(ctx context.Context, fileContent multipart.File, password string) (*models.ImportResponseModel, error) {
	response := &models.ImportResponseModel{}
	var existingUsers []entities.User

	var filter expression.ConditionBuilder

	items, _, err := r.dbContext.GetList(ctx, filter, 0, "")
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		var user entities.User
		err = attributevalue.UnmarshalMap(item, &user)
		if err != nil {
			return nil, err
		}
		existingUsers = append(existingUsers, user)
	}

	defer fileContent.Close()

	excelFile, err := excelize.OpenReader(fileContent, excelize.Options{Password: password})
	if err != nil {
		return response, fmt.Errorf("failed to load the file content: %w", err)
	}

	sheetName := excelFile.GetSheetName(0)
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return response, fmt.Errorf("failed to read rows from excel sheet: %w", err)
	}

	var newUsers []interface{}

	for rowIndex, row := range rows {

		if rowIndex == 0 {
			continue
		}

		user := entities.User{
			UserName:    strings.TrimSpace(row[0]),
			FirstName:   strings.TrimSpace(row[1]),
			LastName:    strings.TrimSpace(row[2]),
			Street:      strings.TrimSpace(row[3]),
			City:        strings.TrimSpace(row[4]),
			Zip:         strings.TrimSpace(row[5]),
			Country:     strings.TrimSpace(row[6]),
			PhoneNumber: strings.TrimSpace(row[7]),
			Email:       strings.TrimSpace(row[8]),
		}

		user.SetId(uuid.New().String())
		user.SetCreator(auth.GetUserName())
		user.SetModifier(auth.GetUserName())
		user.SetModified(time.Now().UTC())
		user.SetCreated(time.Now().UTC())

		userExists := false
		for _, existingUser := range existingUsers {
			if existingUser.UserName == user.UserName &&
				existingUser.FirstName == user.FirstName &&
				existingUser.LastName == user.LastName &&
				existingUser.Street == user.Street &&
				existingUser.City == user.City &&
				existingUser.Zip == user.Zip &&
				existingUser.Country == user.Country &&
				existingUser.PhoneNumber == user.PhoneNumber &&
				existingUser.Email == user.Email {
				userExists = true
				break
			}
		}

		if !userExists {
			newUsers = append(newUsers, user)
		}
	}

	err = r.dbContext.SaveBatch(ctx, newUsers)
	if err != nil {
		return response, fmt.Errorf("failed to add new users: %w", err)
	}

	response.ImportedUsers = len(newUsers)

	return response, nil
}
