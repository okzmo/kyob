package services

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/livekit/protocol/auth"
)

type LivekitResponse struct {
	Token string `json:"token"`
}

func GenerateCallToken(roomName string, userId string) (string, error) {
	apiKey := os.Getenv("LIVEKIT_API_KEY")
	apiSecret := os.Getenv("LIVEKIT_API_SECRET")

	at := auth.NewAccessToken(apiKey, apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         roomName,
		CanPublish:   aws.Bool(true),
		CanSubscribe: aws.Bool(true),
	}

	at.SetVideoGrant(grant).SetIdentity(userId).SetValidFor(time.Hour)

	return at.ToJWT()
}
