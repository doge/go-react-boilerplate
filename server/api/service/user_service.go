package service

import (
	"context"
	"errors"
	"fmt"
	"server/api/dto"
	"server/api/repository"
	"server/internal/models"
	"server/internal/security"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserService struct {
	repo   repository.UserRepository
	crypto *security.AESGCM
}

func NewUserService(ur repository.UserRepository, crypto *security.AESGCM) *UserService {
	return &UserService{
		repo:   ur,
		crypto: crypto,
	}
}

func (us *UserService) Login(ctx context.Context, p dto.UserLoginPayload) (*LoginResult, error) {

	foundUser, err := us.repo.FindUserByUsername(ctx, p.Username)
	if err != nil {
		return nil, err
	}

	passwordMatch, err := security.VerifyPassword(p.Password, foundUser.Password)
	if err != nil {
		return nil, err
	}
	if !passwordMatch {
		return nil, ErrInvalidCredentials
	}

	decryptedEmail, err := us.crypto.Decrypt(foundUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt email for %s: %v", foundUser.ID.Hex(), err)
	}

	return &LoginResult{
		Email: string(decryptedEmail),
	}, nil
}

func (us *UserService) Register(p dto.UserRegistrationPayload) (bson.ObjectID, error) {

	hashedPassword, err := security.GenerateHash(p.Password)
	if err != nil {
		return bson.NilObjectID, errors.New(`error hashing password`)
	}

	encryptedEmail, err := us.crypto.Encrypt([]byte(p.Email))
	if err != nil {
		return bson.NilObjectID, errors.New(`failed to encrypt`)
	}

	return us.repo.Insert(models.User{
		ID:       bson.NewObjectID(),
		Username: p.Username,
		Email:    encryptedEmail,
		Password: hashedPassword,
	})
}
