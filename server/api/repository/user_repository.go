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

type UserRepository interface {
	EnsureIndexes(context.Context) error
	Insert(models.User) (bson.ObjectID, error)
	FindUserByUsername(context.Context, string) (*models.User, error)
	FindUserByID(context.Context, bson.ObjectID) (*models.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return userRepository{
		collection: collection,
	}
}

func (ur userRepository) EnsureIndexes(ctx context.Context) error {
	_, err := ur.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "email_hash", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetPartialFilterExpression(bson.M{"email_hash": bson.M{"$exists": true}}),
		},
	})
	return err
}

func (ur userRepository) Insert(u models.User) (bson.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := ur.collection.InsertOne(ctx, u)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return bson.NilObjectID, ErrConflict
		}
		return bson.NilObjectID, err
	}

	return result.InsertedID.(bson.ObjectID), err
}

func (ur userRepository) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {

	var user models.User

	err := ur.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("[error] finding user: %w", err)
	}

	return &user, nil
}

func (ur userRepository) FindUserByID(ctx context.Context, id bson.ObjectID) (*models.User, error) {
	var user models.User

	err := ur.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("[error] finding user by id: %w", err)
	}

	return &user, nil
}
