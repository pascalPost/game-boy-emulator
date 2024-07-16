package internal

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestByteUnmarshal(t *testing.T) {
//	u := ByteKey{}
//	err := json.Unmarshal([]byte(`{"value":"0x00"}`), &u)
//	assert.NoError(t, err)
//}

func TestByteUnmarshal2(t *testing.T) {
	u := map[ByteKey]int{}
	err := json.Unmarshal([]byte(`{"0x00":0,"0x01":1}`), &u)
	assert.NoError(t, err)
}
