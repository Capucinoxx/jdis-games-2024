package model

import "fmt"

// Client représente un client connecté au serveur.
type Client struct {
	Out        chan []byte
	In         chan ClientMessage
	Connection Connection
}

// ClientMessage représente un message envoyé par un client.
type ClientMessage struct {
	MessageType uint8
	Body        interface{}
}

// String retourne une représentation textuelle de ClientMessage.
func (msg ClientMessage) String() string {
	return fmt.Sprintf("ClientMessage{MessageType: %d, Body: %+v}", msg.MessageType, msg.Body)
}
