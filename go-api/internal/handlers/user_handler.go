package handlers

import (
	"log"

	"github.com/Energie-Burgenland/ausaestung-info/internal/models"
	"github.com/Energie-Burgenland/ausaestung-info/internal/repositories"
	"github.com/Energie-Burgenland/ausaestung-info/internal/validation"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	repo      *repositories.UserRepository
	validator *validation.Validator
}

func NewUserHandler(repo *repositories.UserRepository, validator *validation.Validator) *UserHandler {
	return &UserHandler{
		repo:      repo,
		validator: validator,
	}
}

// @Summary      Get All Users
// @Description  Retrieve a list of all users
// @Tags         Users
// @Param        lastEvaluatedKey  query  string  false  "Last evaluated key"
// @Param        textQuery  query  string  false  "text search query"
// @Produce      json
// @Security 	 OAuth2Implicit
// @Success      200  {array} User
// @Router       /user [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	lastEvaluatedKey := c.Query("lastEvaluatedKey", "")
	textQuery := c.Query("textQuery", "")

	users, err := h.repo.GetUsers(c.UserContext(), lastEvaluatedKey, textQuery)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(users)
}

// @Summary      Get User by ID
// @Description  Retrieve a specific user by ID
// @Tags         Users
// @Produce      json
// @Security 	 OAuth2Implicit
// @Param        id   path      string  true  "Users ID"
// @Success      200  {object}  User
// @Router       /user/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.repo.GetUser(c.UserContext(), id)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(user)
}

// @Summary      Create User
// @Description  Create a new user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security 	 OAuth2Implicit
// @Param        request  body      SaveUser  true  "Create User Request"
// @Success      201      {object}  User            "Successfully created"
// @Router       /user [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	saveModel := new(models.SaveUserModel)

	if err := c.BodyParser(saveModel); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if err := h.validator.ValidateSave(c, saveModel); err != nil {
		return err
	}

	user, err := h.repo.CreateUser(c.UserContext(), saveModel)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(201).JSON(user)
}

// @Summary      Update User
// @Description  Update a new user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security 	 OAuth2Implicit
// @Param        id  path      string  true  "User ID"
// @Param        request  body      SaveUser  true  "Update User Request"
// @Success      200      {object}  User            "Successfully created"
// @Router       /user/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	saveModel := new(models.SaveUserModel)

	if err := c.BodyParser(saveModel); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if err := h.validator.ValidateSaveWithId(c, saveModel, id); err != nil {
		return err
	}

	user, err := h.repo.UpdateUser(c.UserContext(), id, saveModel)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(user)
}

// @Summary      Delete User
// @Description  Delete an existing user by ID
// @Tags         Users
// @Security 	 OAuth2Implicit
// @Param        id  path      string  true  "User ID"
// @Success      204      "No Content"
// @Router       /user/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.repo.DeleteUser(c.UserContext(), id)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendStatus(204)
}

// @Summary      Import users
// @Description  Imports users from a excel file
// @Accept       multipart/form-data
// @Tags         Users
// @Param        file  formData  file  true  "Excel file to upload"
// @Param        password  query  string  false  "Password used to protect the Excel file (option)"
// @Produce      json
// @Security 	   OAuth2Implicit
// @Success      200  {array} ImportResponse
// @Router       /user/import [post]
func (h *UserHandler) ImportUsers(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	password := c.Query("password", "")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to get uploaded file",
		})
	}

	// Log the filename and filesize
	log.Printf("Uploaded file: %s, size: %d bytes, header: %s", file.Filename, file.Size, file.Header)

	// Open the uploaded file in memory
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	response, err := h.repo.ImportUsers(c.UserContext(), fileContent, password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.JSON(response)
}
