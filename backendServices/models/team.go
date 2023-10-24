package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Team struct {
	Id            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" validate:"required"`
	Members       []string           `json:"members,omitempty" validate:"required"`
	MatchesPlayed int                `json:"matchesPlayed"`
	MatchesWon    int                `json:"matchesWon"`
	MatchesLost   int                `json:"matchesLost"`
	MatchesTied   int                `json:"matchesTied"`
	MatchPoints   int                `json:"matchPoints"`
	Captain       string             `json:"captain,omitempty"`
	PlayingXI     []string           `json:"playingXI,omitempty"`
}
