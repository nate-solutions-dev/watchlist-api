package controller

import (
	"net/http"

	"github.com/Friel909/watchlist-api/internal/dto"
	"github.com/Friel909/watchlist-api/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// @Summary      Register user
// @Description  Create a new account using username, email, password, and region
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RegisterRequest  true  "Request body"
// @Success      201   {object}  dto.Response{result=map[string]string}
// @Failure      400   {object}  dto.Response
// @Router       /public/auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Message:  err.Error(),
			Response: http.StatusBadRequest,
			Result:   nil,
		})
		return
	}

	if err := ac.authService.Register(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Message:  err.Error(),
			Response: http.StatusBadRequest,
			Result:   nil,
		})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Message:  "register success",
		Response: http.StatusCreated,
		Result:   gin.H{"message": "register success"},
	})
}

// @Summary      Login user
// @Description  Authenticate user and return JWT token with user info
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginRequest  true  "Request body"
// @Success      200   {object}  dto.Response{result=dto.AuthResponse}
// @Failure      400   {object}  dto.Response
// @Failure      401   {object}  dto.Response
// @Router       /public/auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Message:  err.Error(),
			Response: http.StatusBadRequest,
			Result:   nil,
		})
		return
	}

	resp, err := ac.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Message:  err.Error(),
			Response: http.StatusUnauthorized,
			Result:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Message:  "login success",
		Response: http.StatusOK,
		Result:   resp,
	})
}

// @Summary      Get current user
// @Description  Return authenticated user identity extracted from JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  dto.Response{result=dto.MeResponse}
// @Failure      401            {object}  dto.Response
// @Router       /private/auth/me [get]
func (ac *AuthController) Me(c *gin.Context) {
	callerID, _ := c.Get("caller_id")
	callerUsername, _ := c.Get("caller_username")

	c.JSON(http.StatusOK, dto.Response{
		Message:  "success",
		Response: http.StatusOK,
		Result: gin.H{
			"user_data_id": callerID,
			"username":     callerUsername,
		},
	})
}
