package network

// Package network provides an abstraction for managing network connections
// using the WebSocket protocol for real-time bidirectional communication
// between the client and the server. This includes registering and
// unregistering connections, upgrading HTTP requests to WebSocket
// connections, and managing message reading and writing.
//
// This package utilizes the "github.com/gorilla/websocket" package for
// managing WebSocket connections.

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/capucinoxx/forlorn/pkg/model"
	"github.com/gorilla/websocket"
)

// Network represents a network server capable of managing WebSocket connections.
type Network struct {
	// upgrader configures the settings for upgrading HTTP requests to
	// WebSocket connections.
	upgrader websocket.Upgrader

	// register is a function called when a new connection is established.
	register func(conn model.Connection)

	// uregister is a function called when a connection is closed.
	uregister func(conn model.Connection)

	// port is the port on which the server listens.
	port int

	// address is the IP address on which the server listens.
	address string

	// connected is used to keep track of the currently used tokens.
	connected sync.Map
}

// NewNetwork creates a new network server with the specified address and port.
// It configures a WebSocket upgrader with default buffer sizes and a lenient origin check
// policy and initializes a concurrent map to track active tokens, preventing multiple uses
// of the same token.
func NewNetwork(address string, port int) *Network {
	return &Network{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow connections from any origin.
				// TODO: Generalize the logic to allow restricting origins.
				return true
			},
		},
		port:      port,
		address:   address,
		connected: sync.Map{},
	}
}

// Address returns the IP address and port of the server as a formatted string.
// Format as "address:port"
func (n *Network) Address() string {
	return strings.Join([]string{n.address, strconv.Itoa(n.port)}, ":")
}

// SetRegisterFunc sets the function to be called when a new connection is established.
func (n *Network) SetRegisterFunc(f func(conn model.Connection)) {
	n.register = f
}

// SetUnregisterFunc sets the function to be called when a connection is closed.
func (n *Network) SetUnregisterFunc(f func(conn model.Connection)) {
	n.uregister = f
}

// Register registers a new connection by invoking the specified register function.
func (n *Network) Register(conn model.Connection) {
	if n.register != nil {
		n.register(conn)
	}
}

// Unregister deregisters a connection by invoking the specified unregister function.
// It also removes the token associated with the connection from the active tokens map.
func (n *Network) Unregister(conn model.Connection) {
	if n.uregister != nil {
		n.uregister(conn)
	}

  n.connected.Delete(conn.Identifier())
}

// Init initializes the network server by listening for HTTP requests and upgrading
// them to WebSocket connections.
// It rejects connections with duplicate tokens.
func (n *Network) Init() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "" {
			if _, ok := n.connected.Load(token); ok {
				http.Error(w, "token already in use", http.StatusUnauthorized)
				return
			}

			n.connected.Store(token, true)
		}
  
		ws, err := n.upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "failed to upgrade", http.StatusInternalServerError)
			return
		}
		n.register(NewConnection(ws, token))
	})
}

// Run starts the network server listening for incoming connections on the specified IP
// address and port.
// Returns an error if the server cannot start.
func (n *Network) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", n.address, n.port))
	if err != nil {
		return err
	}
	n.port = listener.Addr().(*net.TCPAddr).Port

	return http.Serve(listener, nil)

}

// Connection is an implementation of the model.Connection interface for WebSocket connections.
// It encapsulates a WebSocket connection pointer and provides methods for reading, writing,
// closing, and pinging the connection.
type Connection struct {
	// conn is the websocket connection.
	conn *websocket.Conn

	// token associatedwith the connection for authentication.
	// If the token is empty, it represents a read-only connection.
	token string
}

// NewConnection creates a new WebSocket connection instance.
func NewConnection(conn *websocket.Conn, token string) *Connection {
	return &Connection{conn: conn, token: token}
}

// Identifier returns a string representing the token associated with the connection,
// which is used as a unique identifier for managing connections.
func (c *Connection) Identifier() string {
	return c.token
}

// Close closes the connection by sending a close message and then closing the underlying
// WebSocket connection.
func (c *Connection) Close(writeWait time.Duration) {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})

}

// PrepareRead configures the connection for reading by setting the maximum message size,
// the pong wait timeout, and the pong handler. This method should be called before reading
// each message.
func (c *Connection) PrepareRead(maxMessageSize int64, pongWait time.Duration) {
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}

// Read reads a message from the connection and returns the message read or an error.
// This method should be called after PrepareRead has been invoked.
func (c *Connection) Read() ([]byte, error) {
	_, msg, err := c.conn.ReadMessage()
	return msg, err
}

// PrepareWrite configures the connection for writing by setting the write deadline.
// This method should be called before each message is written.
func (c *Connection) PrepareWrite(writeWait time.Duration) {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
}

// Write writes a message into the connection and returns an error if the writing fails.
// This method should be called after PrepareWrite has been invoked.
func (c *Connection) Write(msg []byte) error {
	writer, err := c.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}

	writer.Write(msg)
	return writer.Close()
}

// Ping sends a ping message to the connection to check if it is still active.
// This method should be called regularly to maintain the connection active.
func (c *Connection) Ping(writeWait time.Duration) {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	_ = c.conn.WriteMessage(websocket.PingMessage, nil)
}
