package nativemessaging

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strings"
	"testing"
)

type TestStruct struct {
	Name string
}

func TestReceive(t *testing.T) {
	testJsonStr := `
{ "name": "Receive" }
`

	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, uint32(len(testJsonStr)))
	if err != nil {
		t.Fatalf("Fail write header length")
	}

	receiveTarget := &TestStruct{}
	buffer.Write([]byte(testJsonStr))
	err = Receive(receiveTarget, buffer)

	if receiveTarget.Name != "Receive" {
		t.Fatalf("Expected: Receive, got: %s", receiveTarget.Name)
	}
}

func TestSend(t *testing.T) {
	testSendStruct := &TestStruct{Name: "Send"}

	buffer := new(bytes.Buffer)
	if err := Send(testSendStruct, buffer); err != nil {
		t.Fatalf("Fail send")
	}

	jsonData, _ := json.Marshal(testSendStruct)

	length := uint32(0)
	if err := binary.Read(buffer, binary.LittleEndian, &length); err != nil {
		t.Fatal("Read header error")
	}
	// header length check
	if length != uint32(len(jsonData)) {
		t.Fatalf("header length is incorrect. header length: %d, json length: %d", length, len(jsonData))
	}

	// contents check
	if strings.Index(buffer.String(), string(jsonData)) < 0 {
		t.Fatalf("Unexpected return")
	}
}
