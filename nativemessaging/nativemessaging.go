package nativemessaging

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
)

var (
	byteOrder binary.ByteOrder = binary.LittleEndian
)

// Receive json data.
func Receive(data interface{}, reader io.Reader) error {
	message, err := receiveMessage(reader)

	if err != nil {
		return err
	}

	return json.Unmarshal(message, data)
}

// Send json data.
func Send(data interface{}, writer io.Writer) error {
	jsonData, err := json.Marshal(data)

	if err != nil {
		return err
	}

	return sendMessage(jsonData, writer)
}

func receiveMessage(reader io.Reader) ([]byte, error) {
	// for 4byte
	var length uint32

	if err := binary.Read(reader, byteOrder, &length); err != nil {
		return nil, err
	}

	if length == 0 {
		return nil, nil
	}

	message := make([]byte, length)

	n, err := reader.Read(message)
	if err != nil {
		return nil, err
	}
	if n != len(message) {
		return nil, errors.New("message length is different.")
	}

	return message, nil
}

func sendMessage(message []byte, writer io.Writer) error {
	header := make([]byte, 4)

	byteOrder.PutUint32(header, (uint32)(len(message)))

	n, err := writer.Write(header)
	if err != nil {
		return err
	}

	if n != len(header) {
		return errors.New("header length is different.")
	}

	n, err = writer.Write(message)
	if err != nil {
		return err
	}

	if n != len(message) {
		return errors.New("message length is different.")
	}

	return nil
}
