package authservice

import (
	"context"
	"log/slog"
	"sso/internal/domain/models"
	"time"
)

type AuthService struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvide  AppProvide
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (user models.User, err error)
	IsAmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvide interface {
	App(ctx context.Context, appId int) (models.App, error)
}

// New return new instance of the AuthService
func New(
	log *slog.Logger,
	usrSaver UserSaver,
	usrProvider UserProvider,
	appProvide AppProvide,
	tokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		log:         log,
		usrSaver:    usrSaver,
		usrProvider: usrProvider,
		appProvide:  appProvide,
		tokenTTL:    tokenTTL,
	}
}

func (a *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (string, error) {
	panic("not implemented")
}

func RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	panic("not implemented")
}

func IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	panic("not implemented")
}
