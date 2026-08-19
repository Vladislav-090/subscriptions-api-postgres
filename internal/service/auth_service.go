package service

import (
	"context"
	"database/sql"
	"errors"
	"subscriptions-api-postgres/internal/auth"
	"subscriptions-api-postgres/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, email string, passwordHash string, role models.Role) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type RevokedTokensRepository interface {
	RevokeToken(ctx context.Context, jti string, expiresAt time.Time) error
	IsRevoked(ctx context.Context, jti string) (bool, error)
}

type AuthService struct {
	usersRepo         UsersRepository
	revokedTokensRepo RevokedTokensRepository
	jwtManager        *auth.JWTManager
}

func NewAuthService(repo UsersRepository, revokedTokensRepo RevokedTokensRepository, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{
		usersRepo:         repo,
		revokedTokensRepo: revokedTokensRepo,
		jwtManager:        jwtManager,
	}
}

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailRequired      = errors.New("email cannot be empty")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters")

	pgUniqueViolationCode = "23505"
)

func (s *AuthService) Register(ctx context.Context, input models.UserRegisterInput) (models.User, error) {
	if input.Email == "" {
		return models.User{}, ErrEmailRequired
	}

	if len(input.Password) < 8 {
		return models.User{}, ErrPasswordTooShort
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user, err := s.usersRepo.CreateUser(ctx, input.Email, string(passwordHash), models.RoleUser)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolationCode {
			return models.User{}, ErrEmailTaken
		}
		return models.User{}, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, input models.UserLoginInput) (string, error) {
	user, err := s.usersRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.jwtManager.Generate(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Logout(ctx context.Context, claims *auth.Claims) error {
	return s.revokedTokensRepo.RevokeToken(ctx, claims.ID, claims.ExpiresAt.Time)
}

func (s *AuthService) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
	return s.revokedTokensRepo.IsRevoked(ctx, jti)
}
