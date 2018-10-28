package meow

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"io/ioutil"
	"strconv"
	"testing"
)

var testVectorsFilename = flag.String("testvectors", "testdata/testvectors.json", "test vectors filename")

type TestVector struct {
	Seed   uint64
	Input  []byte
	Hash   []byte
	Hash64 uint64
	Hash32 uint32
}

func (v *TestVector) UnmarshalJSON(b []byte) error {
	var data struct {
		Seed   string `json:"seed"`
		Input  string `json:"input"`
		Hash   string `json:"hash"`
		Hash64 string `json:"hash64"`
		Hash32 string `json:"hash32"`
	}

	var err error

	if err = json.Unmarshal(b, &data); err != nil {
		return err
	}

	if v.Seed, err = strconv.ParseUint(data.Seed, 16, 64); err != nil {
		return err
	}

	if v.Input, err = hex.DecodeString(data.Input); err != nil {
		return err
	}

	if v.Hash, err = hex.DecodeString(data.Hash); err != nil {
		return err
	}

	if v.Hash64, err = strconv.ParseUint(data.Hash64, 16, 64); err != nil {
		return err
	}

	h32, err := strconv.ParseUint(data.Hash32, 16, 32)
	if err != nil {
		return err
	}
	v.Hash32 = uint32(h32)

	return nil
}

type TestData struct {
	Version     int          `json:"version_number"`
	VersionName string       `json:"version_name"`
	TestVectors []TestVector `json:"test_vectors"`
}

var testdata *TestData

func LoadTestData(t *testing.T) *TestData {
	t.Helper()

	if testdata != nil {
		return testdata
	}

	b, err := ioutil.ReadFile(*testVectorsFilename)
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(b, &testdata); err != nil {
		t.Fatal(err)
	}

	return testdata
}

func AssertBytesEqual(t *testing.T, expect, got []byte) {
	if len(expect) != len(got) {
		t.Fatalf("length mismatch got=%d expect=%d", len(expect), len(got))
	}

	if bytes.Equal(expect, got) {
		return
	}

	n := len(expect)
	delta := make([]byte, n)
	for i := 0; i < len(expect); i++ {
		delta[i] = expect[i] ^ got[i]
	}

	t.Fatalf("expected equal\n   got=%x\nexpect=%x\n delta=%x\n", got, expect, delta)
}
