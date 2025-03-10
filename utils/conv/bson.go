package conv

import "go.mongodb.org/mongo-driver/v2/bson"

func BsonObjectIDPtrToStringPtr(objectID *bson.ObjectID) *string {
	if objectID == nil {
		return nil
	}

	str := objectID.Hex()
	return &str
}
