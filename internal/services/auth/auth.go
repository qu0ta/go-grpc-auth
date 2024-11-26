package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/qu0ta/go-grpc-auth/internal/domain/models"
	"github.com/qu0ta/go-grpc-auth/internal/storage"
	"github.com/qu0ta/go-grpc-auth/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Auth struct {
	log      *slog.Logger
	storage  Storage
	tokenTTL time.Duration
}

type Storage interface {
	SaveUser(ctx context.Context, email string, passwordHash []byte, appId int32) (uid int64, err error)
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
	App(ctx context.Context, id int32) (models.App, error)
}

// New creates a new Auth instance with the given logger, storage, and token TTL.

func New(log *slog.Logger, storage Storage, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:      log,
		storage:  storage,
		tokenTTL: tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("logging in")
	user, err := a.storage.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found: ", err.Error())

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get the user: ", err.Error())

		return "", fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials")

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.storage.App(ctx, user.AppID)
	if err != nil {
		a.log.Error("failed to get the app: ", err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("logged in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to create token: ", err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil

}
func (a *Auth) RegisterUser(ctx context.Context, email string, password string, appId int32) (int64, error) {
	const op = "auth.RegisterUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering new user")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password: ", err.Error())
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.storage.SaveUser(ctx, email, passwordHash, appId)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("user already exists: ", err.Error())
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to save user: ", err.Error())
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User registered")

	return id, nil

}

// IsAdmin checks if the user with the given userID is an admin.
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.RegisterUser"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userID", userID),
	)

	log.Info("check if user is admin")

	isAdmin, err := a.IsAdmin(ctx, userID)
	if err != nil {
		log.Error("failed to check if user is admin: ", err.Error())
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
