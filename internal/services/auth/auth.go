package authservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/lib/logger/sl"
	"sso/internal/storage"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvide interface {
	App(ctx context.Context, appId int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

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
	const op = "authservice.Login"

	log := a.log.With(
		slog.String("op", op),
	)

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found", sl.Err(err))
			return "", ErrInvalidCredentials
		}

		log.Error("error geting user", sl.Err(err))
		return "", fmt.Errorf("%s %v", op, err)
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		log.Error("invalid credentials", sl.Err(err))
		return "", ErrInvalidCredentials
	}

	app, err := a.appProvide.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s %v", op, err)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s %v", op, err)
	}

	log.Info("login user")
	return token, nil
}

func (a *AuthService) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	const op = "authservice.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
	)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("passHash generation failed", sl.Err(err))
		return 0, fmt.Errorf("%s %v", op, err)
	}

	uid, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			log.Error("user already exists", sl.Err(err))
			return 0, ErrUserAlreadyExists
		}

		log.Error("save user failed", sl.Err(err))
		return 0, fmt.Errorf("%s %v", op, err)
	}

	log.Info("register user")
	return uid, nil
}

func (a *AuthService) IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	const op = "authservice.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
	)

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Error("app not found", sl.Err(err))
			return false, fmt.Errorf("%s %v", op, ErrInvalidAppId)
		}

		if errors.Is(err, ErrUserNotFound) {
			log.Error("user not found", sl.Err(err))
			return false, fmt.Errorf("%s %v", op, ErrUserNotFound)
		}

		log.Error("check is admin error", sl.Err(err))
		return false, fmt.Errorf("%s %v", op, err)
	}

	return isAdmin, nil
}
