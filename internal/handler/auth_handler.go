package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/golang/internal/usecase"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

type CreateUserRequest struct {
	Email    string   `json:"email" example:"admin@example.com"`
	Password string   `json:"password" example:"securePassword123"`
	Roles    []string `json:"roles" example:"admin,user"`
}

type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secret123"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secret123"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// CreateUser godoc
// @Summary      Create a new user (Admin)
// @Description  Create a new user with specified roles
// @Tags         admin, users
// @Accept       json
// @Produce      json
// @Param        request body CreateUserRequest true "User details (email, password, roles)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /admin/users [post]
// @Security     BearerAuth
func (h *AuthHandler) CreateUser(c fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.authUsecase.CreateUser(req.Email, req.Password, req.Roles)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User created successfully"})
}

// GetUserByID godoc
// @Summary      Get user by ID (Admin)
// @Description  Retrieve user information by their ID
// @Tags         admin, users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /admin/users/{id} [get]
// @Security     BearerAuth
func (h *AuthHandler) GetUserByID(c fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.authUsecase.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

// GetAllUsers godoc
// @Summary      Get all users (Admin)
// @Description  Retrieve a list of all users
// @Tags         admin, users
// @Accept       json
// @Produce      json
// @Success      200 {array} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /admin/users [get]
// @Security     BearerAuth
func (h *AuthHandler) GetAllUsers(c fiber.Ctx) error {
	users, err := h.authUsecase.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// DeleteUser godoc
// @Summary      Delete user (Admin)
// @Description  Delete a user by their ID
// @Tags         admin, users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /admin/users/{id} [delete]
// @Security     BearerAuth
func (h *AuthHandler) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")
	err := h.authUsecase.DeleteUser(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Registration details (email, password)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /register [post]
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.authUsecase.Register(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and receive access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login details (email, password)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /login [post]
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	accessToken, refreshToken, err := h.authUsecase.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"access_token": accessToken, "refresh_token": refreshToken})
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Get a new access token using a refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RefreshRequest true "Refresh token (refresh_token)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /refresh [post]
func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	var req RefreshRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	accessToken, refreshToken, err := h.authUsecase.Refresh(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"access_token": accessToken, "refresh_token": refreshToken})
}