package repo

import (
	"context"
	"errors"

	"compConfigurator/internal/db"
	domainErr "compConfigurator/internal/domain/errors"
	"compConfigurator/internal/domain/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UsersRepo struct {
	pool *db.Pool
}

func NewUsersRepo(pool *db.Pool) *UsersRepo {
	return &UsersRepo{pool: pool}
}

const usersSelectBase = `
SELECT
  u.user_id, u.email, u.nickname, u.avatar_url, u.created_at,
  u.password_hash,
  r.code as role_code
FROM pc_configurator.users u
JOIN pc_configurator.roles r ON r.role_id = u.role_id
`

func (r *UsersRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var u models.User

	err := r.pool.QueryRow(ctx, usersSelectBase+` WHERE u.email = $1`, email).
		Scan(&u.UserID, &u.Email, &u.Nickname, &u.AvatarURL, &u.CreatedAt, &u.PasswordHash, &u.RoleCode)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, domainErr.ErrNotFound
		}
		return models.User{}, err
	}
	return u, nil
}

func (r *UsersRepo) GetByID(ctx context.Context, id int64) (models.User, error) {
	var u models.User

	err := r.pool.QueryRow(ctx, usersSelectBase+` WHERE u.user_id = $1`, id).
		Scan(&u.UserID, &u.Email, &u.Nickname, &u.AvatarURL, &u.CreatedAt, &u.PasswordHash, &u.RoleCode)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, domainErr.ErrNotFound
		}
		return models.User{}, err
	}
	return u, nil
}

func (r *UsersRepo) CreateUser(ctx context.Context, email, passwordHash string, nickname *string) (models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx, `
	INSERT INTO pc_configurator.users (email, password_hash, nickname, role_id)
	VALUES ($1, $2, $3, (SELECT role_id FROM pc_configurator.roles WHERE code = 'USER'))
	RETURNING user_id
`, email, passwordHash, nickname).Scan(&u.UserID)

	if err != nil {
		// уникальность email → conflict
		// pgx вернёт PgError, но чтобы не тащить типы: можно проверять SQLSTATE "23505"
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.User{}, domainErr.ErrConflict
		}
		return models.User{}, err
	}
	return r.GetByID(ctx, u.UserID)
}
