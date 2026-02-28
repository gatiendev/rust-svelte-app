package repository

import (
	"database/sql"
	"errors"

	"myproject/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(email, password string) (*models.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hashed),
	}
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	err = r.db.QueryRowx(query, email, string(hashed)).StructScan(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.Get(&user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) CheckPassword(user *models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}
