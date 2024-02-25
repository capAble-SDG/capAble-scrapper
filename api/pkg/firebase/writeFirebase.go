package firebase

import (
	"capAble-scrapper/pkg/crypto"
	"capAble-scrapper/pkg/objects"
	"context"
	"fmt"
)

func (fb *FireBaseSvcImpl) WriteToFirebase(data []objects.Opportunity) error {
	collectionName := "opportunities"

	batch := fb.FirestoreClient.BulkWriter(context.Background())
	opportunitiesCollection := fb.FirestoreClient.Collection(collectionName)

	fmt.Printf("Received %d\n", len(data))
	for count, obj := range data {
		if obj.Company == "" {
			continue
		}
		docID, err := crypto.HashObject(obj)
		if err != nil {
			return err
		}
		docRef := opportunitiesCollection.Doc(docID)
		_, err = batch.Create(docRef, obj)
		if err != nil {
			return err
		}

		fmt.Printf("Writing Object: %d\n", count)
	}

	batch.End()
	fmt.Println("Objects written to Firestore.")

	return nil
}
