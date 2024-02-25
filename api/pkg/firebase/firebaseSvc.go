package firebase

import (
	"capAble-scrapper/pkg/objects"
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

type FireBaseSvcImpl struct {
	FirestoreClient *firestore.Client
}

type FireBaseSvc interface {
	WriteToFirebase(data []objects.Opportunity) error
	ReadFromFirebase() error
}

func NewService() FireBaseSvc {
	opt := option.WithCredentialsFile(".firebase-cred.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
		return nil
	}
	firestoreClient, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error initializing firestore: %v", err)
		return nil
	}

	return &FireBaseSvcImpl{
		FirestoreClient: firestoreClient,
	}
}
