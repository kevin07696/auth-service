package domain

import "context"

type AuthPubSub struct {
	registerTopicID string
	loginTopicID    string
	pubsub          PubSub
}

func NewAuthPubSub(registerTopicID, loginTopicID string, pubsub PubSub) AuthPubSub {
	return AuthPubSub{
		registerTopicID: registerTopicID,
		loginTopicID:    loginTopicID,
		pubsub:          pubsub,
	}
}

func SendRegisterEvent(ctx context.Context, event RegisterEvent) StatusCode {
	return StatusOK
}
