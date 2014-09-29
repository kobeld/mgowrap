package mgowrap

import (
	"encoding/hex"
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

func ToObjectId(idHex string) (bid bson.ObjectId, err error) {
	var d []byte
	d, err = hex.DecodeString(idHex)
	if err != nil || len(d) != 12 {
		err = fmt.Errorf("Invalid input to ObjectIdHex: %q", idHex)
		return
	}
	bid = bson.ObjectId(d)
	return
}
