package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Energie-Burgenland/ausaestung-info/internal/entities"
	"github.com/Energie-Burgenland/ausaestung-info/internal/models"
	"github.com/Energie-Burgenland/ausaestung-info/utils/auth"
	"github.com/Energie-Burgenland/ausaestung-info/utils/database"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type BaseRepository[TEntity any, TModel any, TSaveModel any] struct {
	dbContext *database.DbContext
}

func NewBaseRepository[TEntity any, TModel any, TSaveModel any](dbContext *database.DbContext) *BaseRepository[TEntity, TModel, TSaveModel] {
	return &BaseRepository[TEntity, TModel, TSaveModel]{
		dbContext: dbContext,
	}
}

func (r *BaseRepository[TEntity, TModel, TSaveModel]) GetList(ctx context.Context, filter expression.ConditionBuilder, lastEvaluatedKey string) (models.ListResult[TModel], error) {
	var result models.ListResult[TModel]

	items, lastEvaluatedKey, err := r.dbContext.GetList(ctx, filter, 100, lastEvaluatedKey)
	if err != nil {
		return result, err
	}

	var list []TEntity

	for _, item := range items {
		var entity TEntity
		err = attributevalue.UnmarshalMap(item, &entity)
		if err != nil {
			return result, err
		}
		list = append(list, entity)
	}

	var models []TModel

	copier.Copy(&models, &list)

	result.Items = models
	result.LastEvaluatedKey = lastEvaluatedKey

	return result, nil
}

func (r *BaseRepository[TEntity, TModel, TSaveModel]) Get(ctx context.Context, id string) (*TModel, error) {
	filter := expression.Name("Id").Equal(expression.Value(id))

	result, err := r.dbContext.GetSingle(ctx, filter)
	if err != nil {
		return nil, err
	}

	var entity TEntity
	err = attributevalue.UnmarshalMap(result, &entity)
	if err != nil {
		return nil, err
	}

	var model TModel

	copier.Copy(&model, &entity)

	return &model, nil
}

func (r *BaseRepository[TEntity, TModel, TSaveModel]) Create(ctx context.Context, model TSaveModel) (*TModel, error) {
	entity, ok := any(new(TEntity)).(entities.IEntityBase)
	if !ok {
		return nil, errors.New("entity does not implement IEntityBase")
	}

	if err := copier.Copy(entity, model); err != nil {
		return nil, fmt.Errorf("failed to copy model to entity: %w", err)
	}

	entity.SetId(uuid.New().String())
	entity.SetCreator(auth.GetUserName())
	entity.SetModifier(auth.GetUserName())
	entity.SetModified(time.Now().UTC())
	entity.SetCreated(time.Now().UTC())

	if _, err := r.dbContext.Save(ctx, entity); err != nil {
		return nil, fmt.Errorf("failed to save entity: %w", err)
	}

	return r.Get(ctx, entity.GetId())
}

func (r *BaseRepository[TEntity, TModel, TSaveModel]) Update(ctx context.Context, id string, model TSaveModel) (*TModel, error) {
	filter := expression.Name("Id").Equal(expression.Value(id))

	result, err := r.dbContext.GetSingle(ctx, filter)
	if err != nil {
		return nil, err
	}

	entity, ok := any(new(TEntity)).(entities.IEntityBase)
	if !ok {
		return nil, errors.New("entity does not implement EntityBase")
	}

	err = attributevalue.UnmarshalMap(result, entity)
	if err != nil {
		return nil, err
	}

	copier.Copy(entity, model)

	entity.SetModifier(auth.GetUserName())
	entity.SetModified(time.Now().UTC())

	_, err = r.dbContext.Save(ctx, entity)
	if err != nil {
		return nil, err
	}

	return r.Get(ctx, id)
}

func (r *BaseRepository[TEntity, TModel, TSaveModel]) Delete(ctx context.Context, id string) error {
	condition := map[string]interface{}{
		"Id": id,
	}
	_, err := r.dbContext.Delete(ctx, condition)
	if err != nil {
		return err
	}

	return nil
}
