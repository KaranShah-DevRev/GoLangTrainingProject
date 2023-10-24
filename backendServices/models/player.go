package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Player struct {
	Id            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" validate:"required"`
	Age           int                `json:"age,omitempty" validate:"required"`
	Role          []string           `json:"role,omitempty" validate:"required"`
	DominantHand  string             `json:"dominantHand,omitempty" validate:"required"`
	BattingAvg    float32            `json:"battingAvg"`
	StrikeRate    float32            `json:"strikeRate"`
	Economy       float32            `json:"economy"`
	MatchesPlayed int                `json:"matchesPlayed"`
	TotalRuns     int                `json:"totalRuns"`
	TotalWickets  int                `json:"totalWickets"`
}
