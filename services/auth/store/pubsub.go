package store

import (
	"cloud.google.com/go/pubsub"
	"context"
)

// SendOTP publishes a pubsub message with phone number and otp as attributes
func (s Store) SendOTP(ctx context.Context, otp, phoneNumber string) {
	msg := pubsub.Message{
		Data: nil,
		Attributes: map[string]string{
			"OTP":          otp,
			"PHONE_NUMBER": phoneNumber,
		},
	}
	s.pubsub.Topic("verification").Publish(ctx, &msg)
}
