package repository

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	Insert(models.User) (bson.ObjectID, error)
	FindUserByUsername(context.Context, string) (*models.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return userRepository{
		collection: collection,
	}
}

func (ur userRepository) Insert(u models.User) (bson.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := ur.collection.InsertOne(ctx, u)
	if err != nil {
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
		return nil, fmt.Errorf("[err] finding user: %w", err)
	}

	return &user, nil
}
