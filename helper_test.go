package ddxf

import (
	"encoding/json"
	"math/rand"
	"testing"

	"reflect"

	assert2 "gotest.tools/assert"
)

func TestHelper(t *testing.T) {
	origin := []byte{0xff, 0x00}
	bytes, err := json.Marshal(origin)

	assert2.Assert(t, err == nil)

	var decoded []byte
	err = json.Unmarshal(bytes, &decoded)
	assert2.Assert(t, err == nil && reflect.DeepEqual(decoded, origin))

	bytes, err = Object2Bytes(0)
	bytes2, err2 := json.Marshal(0)
	assert2.Assert(t, err == nil && err2 == nil && reflect.DeepEqual(bytes, bytes2))

	{
		rand.Seed(42)

		n1 := rand.Intn(100000)

		rand.Seed(42)

		n2 := rand.Intn(100000)

		assert2.Assert(t, n1 == n2)
	}

}
