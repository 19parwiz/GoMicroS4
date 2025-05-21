package usecase

import (
	"context"

	"crypto/rand"
	"encoding/hex"
	"github.com/19parwiz/user-service/internal/adapter/mail"
	"github.com/19parwiz/user-service/internal/domain"
)

type UserUsecase struct {
	aiRepo   AutoIncRepo
	userRepo UserRepo
	pHasher  PasswordHasher
	mailer   *mail.Mailer
}

func generateToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func NewUserUsecase(ai AutoIncRepo, userRepo UserRepo, pHasher PasswordHasher, mailer *mail.Mailer) UserUsecase {
	return UserUsecase{
		aiRepo:   ai,
		userRepo: userRepo,
		pHasher:  pHasher,
		mailer:   mailer,
	}
}

func (uc UserUsecase) Register(ctx context.Context, req domain.User) (domain.User, error) {
	emailFilter := domain.UserFilter{
		Email: &req.Email,
	}
	if exists, _ := uc.userRepo.GetWithFilter(ctx, emailFilter); exists != (domain.User{}) {
		return domain.User{}, domain.ErrUserExists
	}

	id, err := uc.aiRepo.Next(ctx, domain.UserDB)
	if err != nil {
		return domain.User{}, err
	}
	req.ID = id

	req.HashedPassword, err = uc.pHasher.Hash(req.HashedPassword)
	if err != nil {
		return domain.User{}, err
	}

	// === ADD THIS: generate and set email confirmation token here ===
	req.EmailConfirmToken = generateToken() // Implement generateToken to create a random token string

	err = uc.userRepo.Create(ctx, req)
	if err != nil {
		return domain.User{}, err
	}

	confirmationLink := "http://localhost:3000/confirm?email=" + req.Email + "&token=" + req.EmailConfirmToken

	emailBody := "New user registered with email: " + req.Email + "\n\nPlease click the following link to confirm the email:\n" + confirmationLink

	err = uc.mailer.SendEmail([]string{"aliparwizbaktash19@gmail.com"}, "Email Confirmation (Test)", emailBody)

	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:   id,
		Name: req.Name,
	}, nil
}

func (uc UserUsecase) Authenticate(ctx context.Context, req domain.User) (domain.User, error) {
	emailFilter := domain.UserFilter{
		Email: &req.Email,
	}
	existingUser, err := uc.userRepo.GetWithFilter(ctx, emailFilter)
	if err != nil {
		return domain.User{}, err
	}

	if existingUser == (domain.User{}) {
		return domain.User{}, domain.ErrUserNotFound
	}

	isValid := uc.pHasher.Verify(existingUser.HashedPassword, req.HashedPassword)
	if !isValid {
		return domain.User{}, domain.ErrInvalidPassword
	}

	return domain.User{
		ID:   existingUser.ID,
		Name: existingUser.Name,
	}, nil
}

func (uc UserUsecase) Get(ctx context.Context, filter domain.UserFilter) (domain.User, error) {
	user, err := uc.userRepo.GetWithFilter(ctx, filter)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
