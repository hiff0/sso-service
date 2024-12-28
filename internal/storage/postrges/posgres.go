package postrges

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"sso/internal/domain/models"
	"sso/internal/storage"
)

type Storage struct {
	db *sql.DB
}

// New create new instance of PostgreSQL storage
func New() (*Storage, error) {
	const op = "storage.postrges.New"
	// TODO: Вынести параметры подключения в конфиг
	connStr := "host=postgres-sso user=sso password=sso dbname=sso sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func Close(s *Storage) error {
	return s.db.Close()
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postrges.SaveUser"

	var id int64
	err := s.db.QueryRowContext(ctx, "INSERT INTO users(email, pass_hash) VALUES($1, $2) RETURNING id", email, passHash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postrges.User"

	var user models.User
	err := s.db.QueryRowContext(ctx, "SELECT id, email, pass_hash FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postrges.IsAdmin"

	var isAdmin bool
	err := s.db.QueryRowContext(ctx, "SELECT is_admin FROM users WHERE id = $1", userID).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.postrges.App"

	var app models.App
	err := s.db.QueryRowContext(ctx, "SELECT id, name, secret FROM apps WHERE id = $1", id).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
