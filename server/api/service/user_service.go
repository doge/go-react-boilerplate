package service

import (
	"context"
	"errors"
	"fmt"
	"server/api/dto"
	"server/api/repository"
	"server/internal/models"
	"server/internal/security"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserService struct {
	repo         repository.UserRepository
	refreshRepo  repository.RefreshSessionRepository
	crypto       *security.AESGCM
	tokenManager *security.TokenManager
	authSettings AuthSettings
}

func NewUserService(
	ur repository.UserRepository,
	rsr repository.RefreshSessionRepository,
	crypto *security.AESGCM,
	tokenManager *security.TokenManager,
	authSettings AuthSettings,
) *UserService {
	return &UserService{
		repo:         ur,
		refreshRepo:  rsr,
		crypto:       crypto,
		tokenManager: tokenManager,
		authSettings: authSettings,
	}
}

func (us *UserService) Login(ctx context.Context, p dto.UserLoginPayload, meta AuthClientMeta) (*LoginResult, error) {
	if err := p.Validate(); err != nil {
		return nil, ErrInvalidPayload
	}

	foundUser, err := us.repo.FindUserByUsername(ctx, p.Username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
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

	accessToken, err := us.tokenManager.GenerateAccessToken(foundUser.ID.Hex())
	if err != nil {
		return nil, err
	}

	refreshToken, err := security.GenerateRecoveryToken()
	if err != nil {
		return nil, err
	}

	err = us.refreshRepo.Insert(ctx, models.RefreshSession{
		ID:        bson.NewObjectID(),
		UserID:    foundUser.ID,
		TokenHash: security.HashToken(refreshToken),
		ExpiresAt: time.Now().Add(us.authSettings.RefreshTTL),
		CreatedAt: time.Now(),
		IP:        meta.IP,
		UserAgent: meta.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Email:        string(decryptedEmail),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (us *UserService) Register(p dto.UserRegistrationPayload) (bson.ObjectID, error) {
	if err := p.Validate(); err != nil {
		return bson.NilObjectID, ErrInvalidPayload
	}

	hashedPassword, err := security.GenerateHash(p.Password)
	if err != nil {
		return bson.NilObjectID, errors.New(`error hashing password`)
	}

	encryptedEmail, err := us.crypto.Encrypt([]byte(p.Email))
	if err != nil {
		return bson.NilObjectID, errors.New(`failed to encrypt`)
	}

	createdID, err := us.repo.Insert(models.User{
		ID:        bson.NewObjectID(),
		Username:  p.Username,
		Email:     encryptedEmail,
		EmailHash: security.HashNormalizedEmail(p.Email),
		Password:  hashedPassword,
	})
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return bson.NilObjectID, ErrUserExists
		}
		return bson.NilObjectID, err
	}

	return createdID, nil
}

func (us *UserService) Refresh(ctx context.Context, rawRefreshToken string, meta AuthClientMeta) (*RefreshResult, error) {
	if rawRefreshToken == "" {
		return nil, ErrMissingToken
	}

	tokenHash := security.HashToken(rawRefreshToken)
	existingSession, err := us.refreshRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	if existingSession.RevokedAt != nil {
		// Rotated tokens are expected to be revoked. If a stale rotated token is seen,
		// return unauthorized without invalidating active sessions.
		if existingSession.ReplacedBy != nil {
			return nil, ErrInvalidToken
		}

		_ = us.refreshRepo.RevokeActiveByUserID(ctx, existingSession.UserID)
		return nil, ErrTokenReuseDetected
	}
	if time.Now().After(existingSession.ExpiresAt) {
		return nil, ErrExpiredSession
	}

	accessToken, err := us.tokenManager.GenerateAccessToken(existingSession.UserID.Hex())
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := security.GenerateRecoveryToken()
	if err != nil {
		return nil, err
	}

	newSession := models.RefreshSession{
		ID:        bson.NewObjectID(),
		UserID:    existingSession.UserID,
		TokenHash: security.HashToken(newRefreshToken),
		ExpiresAt: time.Now().Add(us.authSettings.RefreshTTL),
		CreatedAt: time.Now(),
		IP:        meta.IP,
		UserAgent: meta.UserAgent,
	}

	err = us.refreshRepo.RotateSession(ctx, *existingSession, newSession)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	return &RefreshResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (us *UserService) Logout(ctx context.Context, rawRefreshToken string) error {
	if rawRefreshToken == "" {
		return ErrMissingToken
	}

	err := us.refreshRepo.RevokeByTokenHash(ctx, security.HashToken(rawRefreshToken))
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return err
	}

	return nil
}

func (us *UserService) Session(ctx context.Context, userID string) (*SessionResult, error) {
	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := us.repo.FindUserByID(ctx, objectID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	decryptedEmail, err := us.crypto.Decrypt(user.Email)
	if err != nil {
		return nil, err
	}

	return &SessionResult{
		UID:   user.ID.Hex(),
		Email: string(decryptedEmail),
	}, nil
}
