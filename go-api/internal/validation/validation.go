package validation

import (
	"context"
	"errors"

	"github.com/Energie-Burgenland/ausaestung-info/internal/models"
	"github.com/Energie-Burgenland/ausaestung-info/utils/database"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
)

type Validator struct {
	validate   *validator.Validate
	dbContext  *database.DbContext
	Translator *ut.Translator
}

func InitValidation(dbContext *database.DbContext) (*Validator, error) {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	validate := validator.New(validator.WithRequiredStructEnabled())

	enTranslations.RegisterDefaultTranslations(validate, trans)

	validator := &Validator{
		validate:   validate,
		dbContext:  dbContext,
		Translator: &trans,
	}

	if err := validator.RegisterUserValidation(dbContext, trans); err != nil {
		return nil, err
	}

	return validator, nil
}

func (v *Validator) ValidateSave(c *fiber.Ctx, s interface{}) *models.ValidationErrorResponseModel {
	return v.ValidateSaveWithId(c, s, "")
}

type contextKey string

const idKey contextKey = "Id"

func (v *Validator) ValidateSaveWithId(c *fiber.Ctx, s interface{}, id string) *models.ValidationErrorResponseModel {
	var ctx context.Context

	if id != "" {
		ctx = context.WithValue(c.UserContext(), idKey, id)
	} else {
		ctx = c.UserContext()
	}

	err := v.validate.StructCtx(ctx, s)

	var errs validator.ValidationErrors
	validationErrors := []models.ValidationErrorModel{}
	errors.As(err, &errs)

	for _, err := range errs {
		var elem models.ValidationErrorModel

		elem.Field = err.Field() // Export struct field name
		elem.ErrorMessage = err.Translate(*v.Translator)
		elem.Value = err.Value().(string)

		validationErrors = append(validationErrors, elem)
	}

	return &models.ValidationErrorResponseModel{
		ErrorCode:    "400",
		ErrorMessage: "One ore more validation errors occurred",
		Errors:       validationErrors,
	}
}
