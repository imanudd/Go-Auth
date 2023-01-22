package user

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterRoute(e *echo.Echo, svc *API) {
	e.POST("/register", svc.Register)
	e.POST("/login", svc.Login)
	e.GET("/token/claim", svc.TokenClaim,
		middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:     &MyClaims{},
			SigningKey: []byte("secret"),
		}))
}

type API struct {
	svc UserService
}

func NewUserAPIImpl(svc UserService) *API {
	return &API{
		svc: svc,
	}
}

func (api *API) Register(c echo.Context) error {
	var user = new(User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	res, err := api.svc.Register(c.Request().Context(), *user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	if res.Status == http.StatusOK {
		return c.JSON(http.StatusOK, res)
	}

	return c.JSON(http.StatusCreated, res)
}

func (api *API) Login(c echo.Context) error {
	var user = new(User)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
	}
	if user.Phone_number == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  http.StatusBadRequest,
			Message: "Phone Number or Password can't be empty",
		})
	}
	token, err := api.svc.Login(c, *user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  http.StatusBadRequest,
			Message: "Phone Number or Password is not match",
		})
	}
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusAccepted, Response{
		Status:  http.StatusAccepted,
		Message: "Login Success",
		Data: map[string]string{
			"token": t,
		},
	})
}
func (api *API) TokenClaim(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims, err := api.svc.TokenClaim(c, user)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, Response{
			Status:  http.StatusUnauthorized,
			Message: "Error Info",
		})
	}
	res := Response{
		Status:  http.StatusOK,
		Message: "Success Get Data From JWT",
		Data:    claims,
	}

	return c.JSON(http.StatusOK, res)
}
