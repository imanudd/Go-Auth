package user

import (
	"context"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type UserRepository interface {
	Register(ctx context.Context, req User) (res Response, err error)
	Login(ctx echo.Context, req User) (User, bool, error)
	NewUser(ctx context.Context, req User) (res Response, err error)
}
type UserService interface {
	Register(c context.Context, req User) (res Response, err error)
	Login(c echo.Context, req User) (jwt *jwt.Token, err error)
	TokenClaim(c echo.Context, user *jwt.Token) (MyClaims, error)
	NewUser(c echo.Context, req User) (res Response, err error)
}
