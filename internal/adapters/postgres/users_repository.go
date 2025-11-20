package postgres

import (
	"context"
	"database/sql"

	"drones/internal/core/domain"
	"drones/internal/ports"
)

type UserRepositoryImpl struct {
	db                   *sql.DB
	logger               ports.Logger
	getByIDStmt          *sql.Stmt
	getByNameAndTypeStmt *sql.Stmt
}

// NewUserRepository creates a new PostgreSQL User repository
func NewUserRepository(db *sql.DB, logger ports.Logger) ports.UserRepository {
	repo := &UserRepositoryImpl{
		db:     db,
		logger: logger,
	}
	// Prepare statements
	if err := repo.prepareStatements(); err != nil {
		logger.Error("Users: Failed to prepare statements", "error", err)
		return repo // Return anyway, fallback to non-prepared queries
	}

	return repo
}

func (r *UserRepositoryImpl) Close() error {
	if r.getByIDStmt != nil {
		if err := r.getByIDStmt.Close(); err != nil {
			r.logger.Error("Failed to close getByIDStmt", "error", err)
			return err
		}
	}
	if r.getByNameAndTypeStmt != nil {
		if err := r.getByNameAndTypeStmt.Close(); err != nil {
			r.logger.Error("Failed to close getByNameAndTypeStmt", "error", err)
			return err
		}
	}
	return nil
}

func (r *UserRepositoryImpl) GetDB() *sql.DB {
	return r.db
}

// prepareStatements prepares the SQL statements for the repository
// Optimizes performance by pre-compiling frequently used queries
func (r *UserRepositoryImpl) prepareStatements() error {
	var err error

	// Get by ID statement
	r.getByIDStmt, err = r.db.Prepare(`
		SELECT u.id, u.name, u.email, u.phone, u.type, u.active, u.country, u.locale, u.device_id, u.avatar_url, u.bio, u.created_at, u.updated_at, d.id as drone_id
		FROM users u
		LEFT JOIN drones d ON d.user_id = u.id
		WHERE u.id = $1`)
	if err != nil {
		return err
	}

	// Get by email and type statement
	r.getByNameAndTypeStmt, err = r.db.Prepare(`
		SELECT u.id, u.name, u.email, u.phone, u.type, u.active, u.country, u.locale, u.device_id, u.notification_token, u.avatar_url, u.bio, u.created_at, u.updated_at, d.id as drone_id
		FROM users u
		LEFT JOIN drones d ON d.user_id = u.id
		WHERE u.name = $1 AND u.type = $2`)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByID implements ports.UserRepository.
func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	var user domain.User
	var err error

	// Use prepared statement if available, otherwise fall back to regular query
	if r.getByIDStmt != nil {
		err = r.getByIDStmt.QueryRowContext(ctx, userID).Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Type,
			&user.Active,
			&user.Country,
			&user.Locale,
			&user.DeviceID,
			&user.AvatarUrl,
			&user.Bio,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DroneId,
		)
	} else {
		err = r.db.QueryRowContext(ctx, `
			SELECT u.id, u.name, u.email, u.phone, u.type, u.active, u.country, u.locale, u.device_id, u.avatar_url, u.bio, u.created_at, u.updated_at, d.id as drone_id
			FROM users u
			LEFT JOIN drones d ON d.user_id = u.id
			WHERE u.id = $1`, userID).Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Type,
			&user.Active,
			&user.Country,
			&user.Locale,
			&user.DeviceID,
			&user.AvatarUrl,
			&user.Bio,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DroneId,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn("User not found", "userID", userID)
			return nil, domain.NewDomainError(domain.UserNotFoundError, "User not found", nil)
		}
		r.logger.Error("Failed to get user by ID", "userID", userID, "error", err)
		return nil, domain.NewDomainError(domain.UnableToProcessError, "Failed to get user by ID", err)
	}

	return &user, nil
}

// GetByEmailAndType implements ports.UserRepository.
func (r *UserRepositoryImpl) GetUserByNameAndType(ctx context.Context, name string, userType string) (*domain.User, error) {
	var user domain.User
	var err error

	// Use prepared statement if available, otherwise fall back to regular query
	if r.getByNameAndTypeStmt != nil {
		err = r.getByNameAndTypeStmt.QueryRowContext(ctx, name, userType).Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Type,
			&user.Active,
			&user.Country,
			&user.Locale,
			&user.DeviceID,
			&user.NotificationToken,
			&user.AvatarUrl,
			&user.Bio,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DroneId,
		)
	} else {
		err = r.db.QueryRowContext(ctx, `
			SELECT u.id, u.name, u.email, u.phone, u.type, u.active, u.country, u.locale, u.device_id, u.notification_token, u.avatar_url, u.bio, u.created_at, u.updated_at, d.id as drone_id
			FROM users u
			LEFT JOIN drones d ON d.user_id = u.id
			WHERE u.name = $1 AND u.type = $2`, name, userType).Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Type,
			&user.Active,
			&user.Country,
			&user.Locale,
			&user.DeviceID,
			&user.NotificationToken,
			&user.AvatarUrl,
			&user.Bio,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DroneId,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn("User not found by name and type", "name", name, "userType", userType)
			return nil, domain.NewDomainError(domain.UserNotFoundError, "User not found by name and type", nil)
		}
		r.logger.Error("Failed to get user by name and type", "name", name, "userType", userType, "error", err)
		return nil, domain.NewDomainError(domain.UnableToProcessError, "Failed to get user by name and type", err)
	}

	return &user, nil
}
