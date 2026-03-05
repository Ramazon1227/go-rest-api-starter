package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Ramazon1227/go-rest-api-starter/models"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/email"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/utils"
	"github.com/Ramazon1227/go-rest-api-starter/storage"
)

type userRepoImpl struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) storage.UserRepoImpl {
	return &userRepoImpl{
		db: db,
	}
}

func (r *userRepoImpl) Add(ctx context.Context, entity *models.UserCreateModel) (pKey *models.PrimaryKey, err error) {
	// Check if user exists and is soft deleted
	checkQuery := `
		SELECT id, deleted_at
		FROM "user"
		WHERE email = $1
	`
	var existingId string
	var deletedAt *time.Time
	err = r.db.QueryRow(ctx, checkQuery, entity.Email).Scan(&existingId, &deletedAt)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	// If user exists and is soft deleted, reactivate them
	if err != pgx.ErrNoRows && deletedAt != nil {
		updateQuery := `
			UPDATE "user"
			SET name = $2,
				role = $3,
				phone = $4,
				deleted_at = NULL,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`
		_, err = r.db.Exec(ctx, updateQuery,
			existingId,
			entity.Name,
			entity.Role,
			entity.Phone,
		)
		if err != nil {
			return nil, err
		}
		return &models.PrimaryKey{Id: existingId}, nil
	}

	// Generate random password
	plainPassword, err := utils.GenerateRandomPassword(8)
	if err != nil {
		return nil, err
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(plainPassword)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO "user" (
		id,
		name,
		email,
		password,
		role,
		phone,
		expires_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, '2072-05-01T11:21:59.001+0000'
	)`

	uuid, err := uuid.NewRandom()
	if err != nil {
		return pKey, err
	}

	_, err = r.db.Exec(ctx, query,
		uuid.String(),
		entity.Name,
		entity.Email,
		hashedPassword,
		entity.Role,
		entity.Phone,
	)
	if err != nil {
		return nil, err
	}

	// Send welcome email with the plain password
	err = email.SendWelcomeEmail(entity.Email, entity.Name, plainPassword)
	if err != nil {
		log.Printf("Failed to send welcome email: %v", err)
	}

	pKey = &models.PrimaryKey{
		Id: uuid.String(),
	}

	return pKey, nil
}

func (r *userRepoImpl) UpdateProfile(ctx context.Context, entity *models.UpdateUserProfileModel) error {
	query := `
		UPDATE "user"
		SET name=$2, email=$3, updated_at=CURRENT_TIMESTAMP
		WHERE id=$1 AND deleted_at IS NULL`

	result, err := r.db.Exec(ctx, query,
		entity.Id,
		entity.Name,
		entity.Email,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return storage.ErrorNotFound
	}

	return nil
}

func (r *userRepoImpl) GetById(ctx context.Context, pKey *models.PrimaryKey) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, name, email, phone, password, role, created_at, updated_at
		FROM "user"
		WHERE id=$1 AND deleted_at IS NULL`

	err := r.db.QueryRow(ctx, query, pKey.Id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepoImpl) GetList(ctx context.Context, queryParam *models.QueryParam) (*models.GetUserListModel, error) {
	var (
		users []*models.User
		count int
	)

	countQuery := `SELECT count(1) FROM "user" WHERE deleted_at IS NULL`
	err := r.db.QueryRow(ctx, countQuery).Scan(&count)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, name, email, role, created_at, updated_at
		FROM "user"
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, queryParam.Limit, queryParam.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return &models.GetUserListModel{
		Users: users,
		Count: count,
	}, nil
}

func (r *userRepoImpl) Delete(ctx context.Context, pKey *models.PrimaryKey) error {
	query := `
		UPDATE "user"
		SET deleted_at=CURRENT_TIMESTAMP
		WHERE id=$1 AND deleted_at IS NULL`

	result, err := r.db.Exec(ctx, query, pKey.Id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return storage.ErrorNotFound
	}

	return nil
}

func (r *userRepoImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	query := `
		SELECT
			id,
			name,
			role,
			phone,
			email,
			password,
			expires_at,
			created_at,
			updated_at
		FROM "user"
		WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Role,
		&user.Phone,
		&user.Email,
		&user.Password,
		&user.ExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, storage.ErrorNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepoImpl) UpdateUserProfile(ctx context.Context, userId string, req *models.UpdateProfileRequest) error {
	var setValues []string
	var args []interface{}
	argId := 1

	if req.Name != "" {
		setValues = append(setValues, fmt.Sprintf("name = $%d", argId))
		args = append(args, req.Name)
		argId++
	}

	if req.Phone != "" {
		setValues = append(setValues, fmt.Sprintf("phone = $%d", argId))
		args = append(args, req.Phone)
		argId++
	}

	if req.Email != "" {
		setValues = append(setValues, fmt.Sprintf("email = $%d", argId))
		args = append(args, req.Email)
		argId++
	}

	if len(setValues) == 0 {
		return nil
	}

	setValues = append(setValues, fmt.Sprintf("updated_at = $%d", argId))
	args = append(args, time.Now())
	argId++

	args = append(args, userId)
	query := fmt.Sprintf(`
		UPDATE "user"
		SET %s
		WHERE id = $%d`,
		strings.Join(setValues, ", "),
		argId,
	)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return storage.ErrorNotFound
	}

	return nil
}

func (r *userRepoImpl) UpdatePassword(ctx context.Context, userId string, currentPassword, newPassword string) error {
	// First, get the user to verify current password
	user, err := r.GetById(ctx, &models.PrimaryKey{Id: userId})
	if err != nil {
		return err
	}

	// Verify current password
	if !utils.CheckPassword(user.Password, currentPassword) {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	query := `
		UPDATE "user"
		SET password = $1, updated_at = $2
		WHERE id = $3`

	result, err := r.db.Exec(ctx, query, hashedPassword, time.Now(), userId)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return storage.ErrorNotFound
	}

	return nil
}
