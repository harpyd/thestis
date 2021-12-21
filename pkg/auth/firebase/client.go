package firebase

import (
	"context"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

const connectTimeout = 10 * time.Second

func NewClient(serviceAccountFilePath string) (*auth.Client, error) {
	opt := option.WithCredentialsFile(serviceAccountFilePath)

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	firebaseApp, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	return firebaseApp.Auth(ctx)
}
