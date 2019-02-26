package nativemessaging

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
)

type TestStruct struct {
	Name string `json:"name"`
}

func TestReceive(t *testing.T) {
	// test json: {"name":"Receive"}
	bs := []byte{18, 0, 0, 0, 123, 34, 110, 97, 109, 101, 34, 58, 34, 82, 101, 99, 101, 105, 118, 101, 34, 125}
	buffer := bytes.NewBuffer(bs)

	receiveTarget := &TestStruct{}

	if err := Receive(receiveTarget, buffer); err != nil {
		t.Fatalf("Fail Receive. err: %s", err)
	}

	if receiveTarget.Name != "Receive" {
		t.Fatalf("Expected: Receive, got: %s", receiveTarget.Name)
	}
}

func TestSend(t *testing.T) {
	testSendStruct := &TestStruct{Name: "Send"}
	exceptJsonStr := `{"name":"Send"}`
	exceptJsonStrLength := uint32(len(exceptJsonStr))

	buffer := new(bytes.Buffer)
	if err := Send(testSendStruct, buffer); err != nil {
		t.Fatalf("Fail Send. err: %s", err)
	}

	length := uint32(0)
	if err := binary.Read(buffer, binary.LittleEndian, &length); err != nil {
		t.Fatalf("Read header error. header length: %d", length)
	}
	// header length check
	if length != exceptJsonStrLength {
		t.Fatalf("header length is incorrect. header length: %d, json length: %d", length, exceptJsonStrLength)
	}

	// contents check
	if strings.Index(buffer.String(), exceptJsonStr) < 0 {
		t.Fatalf("Unexpected return")
	}
}
