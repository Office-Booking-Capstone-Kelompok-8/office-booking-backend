package firebase

import (
	"context"
	"log"

	fb "firebase.google.com/go"
	"firebase.google.com/go/messaging"

	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

func NewApp(ctx context.Context, conf *viper.Viper) *fb.App {
	opt := option.WithCredentialsFile(conf.GetString("service.firebase.keyPath"))
	app, err := fb.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase: %v", err)
	}

	return app
}

func NewMessagingClient(ctx context.Context, app *fb.App) (*messaging.Client, error) {
	messagingClient, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v", err)
	}

	return messagingClient, nil
}
