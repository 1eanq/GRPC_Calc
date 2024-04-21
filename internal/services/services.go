package services

import (
	"GRPC_Calc/internal/domain/models"
	"GRPC_Calc/internal/lib/calculator"
	"GRPC_Calc/internal/lib/jwt"
	"GRPC_Calc/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
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
	User(ctx context.Context, email string) (models.User, error)
}

type Calc struct {
	log          *slog.Logger
	exprSaver    ExpressionSaver
	exprProvider ExpressionProvider
}

type ExpressionSaver interface {
	SaveExpression(ctx context.Context, expr string, uid int64) (int64, error)
}

type ExpressionProvider interface {
	Expression(ctx context.Context, id int64) (models.Expression, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

// NewAuth returns a new interface of the Auth service
func NewAuth(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		tokenTTL:    tokenTTL,
	}
}

func NewCalc(
	log *slog.Logger,
	expressionSaver ExpressionSaver,
	expressionProvider ExpressionProvider,
) *Calc {
	return &Calc{
		log:          log,
		exprSaver:    expressionSaver,
		exprProvider: expressionProvider,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attempting to login user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", err)

			return "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		a.log.Error("failed to get user", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credenticals", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logined successfilly")

	token, err := jwt.NewToken(user, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (a *Auth) RegisterNewUser(ctx context.Context, email string, pass string) (int64, error) {
	const op = "Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// Calculate calculates given expression and returns answer.
// If expression written in wrong spot, returns error
func (c *Calc) CalculateExpression(ctx context.Context, expr string) (string, error) {
	const op = "Auth.Calculate"

	log := c.log.With(
		slog.String("op", op),
		slog.String("expression", expr),
	)

	log.Info("calculating expression")

	res, err := calculator.CalculateExpr(expr)
	if err != nil {
		log.Error("failed to calculate expression", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return fmt.Sprintf("%f", res), nil
}
