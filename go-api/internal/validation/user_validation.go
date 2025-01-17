package validation

import (
	"context"

	"github.com/Energie-Burgenland/ausaestung-info/utils/database"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func (v *Validator) RegisterUserValidation(dbContext *database.DbContext, trans ut.Translator) error {
	if err := v.validate.RegisterValidationCtx("uniqueUserName", func(ctx context.Context, fl validator.FieldLevel) bool {
		return uniqueUser(ctx, fl, dbContext)
	}); err != nil {
		return err
	}

	if err := v.validate.RegisterTranslation("uniqueUserName", trans, func(ut ut.Translator) error {
		return ut.Add("uniqueUserName", "Another user with the username '{0}' already exists", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(" name", fe.Value().(string))
		return t
	}); err != nil {
		return err
	}

	return nil
}

func uniqueUser(ctx context.Context, fl validator.FieldLevel, dbContext *database.DbContext) bool {
	filter := expression.Name("UserName").Equal(expression.Value(fl.Field().String()))

	id, idExists := ctx.Value("Id").(string)
	if idExists {
		filter = filter.And(expression.Name("Id").NotEqual(expression.Value(id)))
	}

	matches, err := dbContext.Count(ctx, filter)
	if err != nil {
		return false
	}

	return matches == 0
}
