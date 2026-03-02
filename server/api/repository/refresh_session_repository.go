package repository

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type RefreshSessionRepository interface {
	EnsureIndexes(context.Context) error
	Insert(context.Context, models.RefreshSession) error
	FindByTokenHash(context.Context, string) (*models.RefreshSession, error)
	RevokeByTokenHash(context.Context, string) error
	RevokeAndReplace(context.Context, bson.ObjectID, bson.ObjectID) error
	RevokeActiveByUserID(context.Context, bson.ObjectID) error
	RotateSession(context.Context, models.RefreshSession, models.RefreshSession) error
}

type refreshSessionRepository struct {
	collection *mongo.Collection
}

func NewRefreshSessionRepository(collection *mongo.Collection) RefreshSessionRepository {
	return refreshSessionRepository{
		collection: collection,
	}
}

func (rsr refreshSessionRepository) EnsureIndexes(ctx context.Context) error {
	_, err := rsr.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "token_hash", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	})
	return err
}

func (rsr refreshSessionRepository) Insert(ctx context.Context, session models.RefreshSession) error {
	_, err := rsr.collection.InsertOne(ctx, session)
	return err
}

func (rsr refreshSessionRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshSession, error) {
	var session models.RefreshSession
	err := rsr.collection.FindOne(ctx, bson.M{"token_hash": tokenHash}).Decode(&session)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("[err] finding refresh session: %w", err)
	}
	return &session, nil
}

func (rsr refreshSessionRepository) RevokeByTokenHash(ctx context.Context, tokenHash string) error {
	now := time.Now()

	result, err := rsr.collection.UpdateOne(ctx, bson.M{"token_hash": tokenHash}, bson.M{
		"$set": bson.M{"revoked_at": now},
	})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (rsr refreshSessionRepository) RevokeAndReplace(ctx context.Context, oldID, newID bson.ObjectID) error {
	now := time.Now()

	result, err := rsr.collection.UpdateOne(ctx, bson.M{"_id": oldID}, bson.M{
		"$set": bson.M{
			"revoked_at":  now,
			"replaced_by": newID,
		},
	})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (rsr refreshSessionRepository) RevokeActiveByUserID(ctx context.Context, userID bson.ObjectID) error {
	now := time.Now()

	_, err := rsr.collection.UpdateMany(ctx, bson.M{
		"user_id":    userID,
		"revoked_at": bson.M{"$exists": false},
	}, bson.M{
		"$set": bson.M{
			"revoked_at": now,
			"expires_at": now,
		},
	})

	return err
}

func (rsr refreshSessionRepository) RotateSession(
	ctx context.Context,
	oldSession models.RefreshSession,
	newSession models.RefreshSession,
) error {
	// Insert new session first. If old session revoke fails due a race,
	// we immediately revoke the new session as compensation.
	if _, err := rsr.collection.InsertOne(ctx, newSession); err != nil {
		return err
	}

	now := time.Now()
	updateResult, err := rsr.collection.UpdateOne(ctx, bson.M{
		"_id":        oldSession.ID,
		"token_hash": oldSession.TokenHash,
		"revoked_at": bson.M{"$exists": false},
	}, bson.M{
		"$set": bson.M{
			"revoked_at":  now,
			"replaced_by": newSession.ID,
		},
	})
	if err != nil {
		_, _ = rsr.collection.UpdateOne(ctx, bson.M{"_id": newSession.ID}, bson.M{
			"$set": bson.M{
				"revoked_at": now,
				"expires_at": now,
			},
		})
		return err
	}
	if updateResult.MatchedCount == 0 {
		_, _ = rsr.collection.UpdateOne(ctx, bson.M{"_id": newSession.ID}, bson.M{
			"$set": bson.M{
				"revoked_at": now,
				"expires_at": now,
			},
		})
		return ErrConflict
	}

	return nil
}
