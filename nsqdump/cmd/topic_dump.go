package cmd

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

const (
	msglensize int = 4
)

const (
	timestampEnd int = 8
	attamptsEnd  int = 10
	messageidEnd int = 42
)

//OneMessage one line of message
type OneMessage struct {
	DataLen     uint32
	TimeStamp   uint64
	Attampts    uint16
	MessageID   string
	MessageBody string
}

//Messages all message in a backup file
type Messages struct {
	Messages []OneMessage
}

//DecodeAndDumpNsqBackFile the structure of one data record is like:
//[recordlen(4 byte)]|[timestamp[0:8]|attampts[8:10]|messageid[10:42]|messagebody[42:]]
func DecodeAndDumpNsqBackFile(origFile, destFile string) error {
	var oneMsg OneMessage
	var allMsgs Messages
	var err error

	infile, err := os.Open(origFile)
	if err != nil {
		return err
	}
	defer infile.Close()

	reader := bufio.NewReader(infile)
	msglenPart := make([]byte, msglensize)

	for {
		if _, err = reader.Read(msglenPart); err != nil {
			break
		}

		//get current message length
		currMsgLen := binary.BigEndian.Uint32(msglenPart)

		//peek the message
		currMsg, err := reader.Peek(int(currMsgLen))
		if err != nil {
			break
		}

		//discard the message
		_, err = reader.Discard(int(currMsgLen))
		if err != nil {
			break
		}

		//convert message byte slice into a message struct
		oneMsg.DataLen = currMsgLen
		ConvByteSliceToMessage(currMsg, &oneMsg)

		allMsgs.Messages = append(allMsgs.Messages, oneMsg)
	}

	if err != io.EOF {
		return err
	}

	//write out file
	outBytes, err := json.Marshal(allMsgs)
	if err != nil {
		return err
	}

	//if the dump file not exist, then create it
	_, err = os.Stat(destFile)
	if os.IsNotExist(err) {

		_, err = os.Create(destFile)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(destFile, outBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

//ConvByteSliceToMessage Convert a byte slice into a Message Struct
func ConvByteSliceToMessage(msgSlice []byte, msg *OneMessage) {
	msg.TimeStamp = binary.BigEndian.Uint64(msgSlice[:timestampEnd])
	msg.Attampts = binary.BigEndian.Uint16(msgSlice[timestampEnd:attamptsEnd])
	msg.MessageID = string(msgSlice[attamptsEnd:messageidEnd])
	msg.MessageBody = string(msgSlice[messageidEnd:])
}
