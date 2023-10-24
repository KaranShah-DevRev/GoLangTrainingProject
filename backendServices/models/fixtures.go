package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Fixture struct {
	Id           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`         // Unique identifier for the match
	MatchNumber  int                `json:"matchNumber,omitempty" validate:"required"` // Unique identifier for the match
	TeamA        string             `json:"teamA,omitempty" validate:"required"`       // Name of the first team
	TeamB        string             `json:"teamB,omitempty" validate:"required"`       // Name of the second team
	MatchDate    time.Time          `json:"matchDate,omitempty" validate:"required"`   // Date and time of the match
	Venue        string             `json:"venue,omitempty" validate:"required"`       // Location where the match is played
	IsFinalMatch bool               `json:"finalMatch,omitempty"`                      // Indicates if it's a final match (true/false)
}
