package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/sethvargo/go-password/password"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type UserRepositoryImpl struct {
	db *bun.DB
}

func NewUserRepositoryImpl(connectionStr string) UserRepository {
	config, err := pgx.ParseConfig(connectionStr)
	if err != nil {
		panic(err)
	}
	config.PreferSimpleProtocol = true
	sqlDb := stdlib.OpenDB(*config)
	db := bun.NewDB(sqlDb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))
	return &UserRepositoryImpl{
		db: db,
	}
}

func (repo *UserRepositoryImpl) Register(ctx context.Context, req User) (res Response, err error) {
	user := new(User)
	err = repo.db.NewSelect().Model(user).Where("phone_number = ?", req.Phone_number).Scan(ctx)
	if err != sql.ErrNoRows {
		res.Data = map[string]string{
			"generate_password ": user.Password,
		}
		res.Message = "userExist"
		return res, err
	}

	pass, _ := password.Generate(4, 0, 0, false, false)
	req.Password = pass
	result, err := repo.NewUser(ctx, req)
	if err != nil {
		return res, err
	}
	res.Data = result.Data
	return res, nil
}

func (repo *UserRepositoryImpl) Login(ctx context.Context, req User) (User, bool, error) {
	var user User
	err := repo.db.NewSelect().Where("phone_number = ? AND password = ?", req.Phone_number, req.Password).Model(&user).Scan(ctx)
	if err == sql.ErrNoRows {
		fmt.Println("User Not Found")
		return user, false, err
	}

	if err != nil {
		//Handle and Show to API
		fmt.Println("Query ERROR")
		return user, false, err
	}
	return user, true, nil
}

func (repo *UserRepositoryImpl) NewUser(ctx context.Context, req User) (res Response, err error) {
	user := &User{Name: req.Name, Password: req.Password, Phone_number: req.Phone_number, Role: req.Role, Created_at: time.Now()}
	_, err = repo.db.NewInsert().Model(user).ExcludeColumn("id").Returning("id").Exec(ctx)
	if err != nil {
		return res, err
	}
	password := user.Password
	res.Data = map[string]string{
		"generate_password ": password,
	}
	return res, nil
}
