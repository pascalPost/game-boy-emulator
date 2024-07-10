package internal

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Operands struct {
	Name      string `json:"name"`
	Bytes     int    `json:"bytes,omitempty"`
	Immediate bool   `json:"immediate"`
}

type Flags struct {
	Z string `json:"Z"`
	N string `json:"N"`
	H string `json:"H"`
	C string `json:"C"`
}

type Opcode struct {
	Mnemonic  string     `json:"mnemonic"`
	Bytes     int        `json:"bytes"`
	Cycles    []int      `json:"cycles"`
	Operands  []Operands `json:"operands"`
	Immediate bool       `json:"immediate"`
	Flags     Flags      `json:"flags"`
}

type ByteKey struct {
	Value byte
}

func (k *ByteKey) UnmarshalJSON(bytes []byte) error {
	text := string(bytes)

	text = RemoveQuotes(text)
	text = RemoveHexPrefix(text)

	decodeString, err := hex.DecodeString(text)
	if err != nil {
		return err
	}

	if len(decodeString) > 1 {
		return errors.New("invalid byte key")
	}

	*k = ByteKey{decodeString[0]}
	return nil
}

func RemoveQuotes(text string) string {
	if len(text) > 0 && text[0] == '"' {
		text = text[1:]
	}
	if len(text) > 0 && text[len(text)-1] == '"' {
		text = text[:len(text)-1]
	}
	return text
}

func RemoveHexPrefix(text string) string {
	if strings.HasPrefix(text, "0x") {
		return text[2:]
	}

	return text
}

func (v *ByteKey) UnmarshalText(text []byte) error {
	return json.Unmarshal(text, v)
}

type OpcodeList struct {
	UnPrefixed map[ByteKey]Opcode `json:"unprefixed"`
	CbPrefixed map[ByteKey]Opcode `json:"cbprefixed"`
}

func ParseOpcodes() (*OpcodeList, error) {
	list := &OpcodeList{}

	_, filename, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(filename), "../", "Opcodes.json")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(file).Decode(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
