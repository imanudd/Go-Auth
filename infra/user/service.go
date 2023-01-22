package user

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type UserServiceImpl struct {
	repo UserRepository
}

func NewUserServiceImpl(repo UserRepository) UserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (svc *UserServiceImpl) Register(c context.Context, req User) (res Response, err error) {
	res, err = svc.repo.Register(c, req)

	if err != nil {
		return res, err
	}

	if res.Message == "userExist" {
		res.Status = http.StatusOK
		res.Message = "Successfully Get Data"
		return res, nil
	}

	res.Status = http.StatusCreated
	res.Message = "Success Adding Data"
	return res, nil

}

func (svc *UserServiceImpl) Login(c echo.Context, req User) (token *jwt.Token, err error) {
	user, res, err := svc.repo.Login(c, req)
	if err != nil {
		return nil, err
	}
	if !res {
		return nil, echo.ErrUnauthorized
	}
	token = jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["phone_number"] = user.Phone_number
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Minute * 2).Unix()
	return token, nil
}

func (svc *UserServiceImpl) TokenClaim(c echo.Context, user *jwt.Token) (MyClaims, error) {
	token, err := jwt.ParseWithClaims(user.Raw, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return MyClaims{}, err
	}

	claims, _ := token.Claims.(*MyClaims)
	myClaims := MyClaims{
		Name:         claims.Name,
		Phone_number: claims.Phone_number,
		Role:         claims.Role,
	}
	return myClaims, nil
}

func (svc *UserServiceImpl) NewUser(c echo.Context, req User) (res Response, err error) {
	res, err = svc.repo.NewUser(c.Request().Context(), req)
	if err != nil {
		return res, err
	}
	res.Status = http.StatusCreated
	res.Message = "Succsefully Insert Data"
	return res, nil
}
