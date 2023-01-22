package app

import (
	"Go-Auth/infra/user"

	"github.com/labstack/echo/v4"
)

func Init() {
	e := echo.New()
	registerProductAPI(e, "postgres")
	e.Logger.Fatal(e.Start(":1323"))
}

func registerProductAPI(e *echo.Echo, db string) {
	var userRepository user.UserRepository
	switch db {
	case "postgres":
		userRepository = user.NewUserRepositoryImpl("postgres://postgres:12345678@localhost:5432/goauth?sslmode=disable")
	default:
		panic(`unknown orm selections!`)
	}

	userService := user.NewUserServiceImpl(userRepository)
	userApi := user.NewUserAPIImpl(userService)
	user.RegisterRoute(e, userApi)
}
