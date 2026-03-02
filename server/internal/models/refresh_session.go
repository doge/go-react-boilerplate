package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RefreshSession struct {
	ID         bson.ObjectID  `bson:"_id"`
	UserID     bson.ObjectID  `bson:"user_id"`
	TokenHash  string         `bson:"token_hash"`
	ExpiresAt  time.Time      `bson:"expires_at"`
	CreatedAt  time.Time      `bson:"created_at"`
	RevokedAt  *time.Time     `bson:"revoked_at,omitempty"`
	ReplacedBy *bson.ObjectID `bson:"replaced_by,omitempty"`
	IP         string         `bson:"ip,omitempty"`
	UserAgent  string         `bson:"user_agent,omitempty"`
}
