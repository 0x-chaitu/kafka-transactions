package suspicioustransaction

import (
	"context"
	"errors"
	"transactions/mongodb"
	"transactions/mongodb/suspicioustransaction/models"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	collection = func(mongoCLient *mongo.Client, databaseName, collectionName string) *mongo.Collection {
		return mongoCLient.Database(databaseName).Collection(collectionName)
	}
	insertCollection = func(
		ctx context.Context,
		collection *mongo.Collection,
		document interface{}) (*mongo.InsertOneResult, error) {
		return collection.InsertOne(ctx, document)
	}
)

func Insert(ctx context.Context, db *mongodb.Mongodb, newSuspiciousTransaction *models.SuspiciousTransaction) error {
	const collectionName = "suspicious_transactions"
	coll := collection(db.Client, db.DatabaseName, collectionName)
	_, err := insertCollection(ctx, coll, newSuspiciousTransaction)
	if err != nil {
		return errors.Unwrap(err)
	}
	return nil
}
