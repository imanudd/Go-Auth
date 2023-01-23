package user

import (
	"context"

	"github.com/golang-jwt/jwt"
)

type UserRepository interface {
	Register(ctx context.Context, req User) (res Response, err error)
	Login(ctx context.Context, req User) (User, bool, error)
	NewUser(ctx context.Context, req User) (res Response, err error)
}
type UserService interface {
	Register(c context.Context, req User) (res Response, err error)
	Login(c context.Context, req User) (jwt *jwt.Token, err error)
	TokenClaim(c context.Context, user *jwt.Token) (MyClaims, error)
	NewUser(c context.Context, req User) (res Response, err error)
}
