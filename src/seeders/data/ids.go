package data

import (
	"crypto/sha1"
	"encoding/hex"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func stableID(parts ...string) bson.ObjectID {
	hash := sha1.New()
	for _, part := range parts {
		_, _ = hash.Write([]byte(part))
		_, _ = hash.Write([]byte{0})
	}
	id, err := bson.ObjectIDFromHex(hex.EncodeToString(hash.Sum(nil))[:24])
	if err != nil {
		panic(err)
	}
	return id
}
