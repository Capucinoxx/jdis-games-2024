package manager

// Network manager provides the necessary tools to manage network communications
// for a multiplayer game environment using WebSocket. This package contains the
// NetworkManager which orchestrates the handling of client connections, message
// transmission, and the enforcement of network protocols.
//
// The manager is responsible for:
// - Initializing and running the network server.
// - Registering and unregistering clients as they connect and disconnect.
// - Broadcasting game state updates and other messages to all connected clients.
// - Handling incoming messages from clients and routing them appropriately.
// - Performing periodic health checks on connections via pings.
// - Gracefully handling client disconnections and ensuring resource cleanup.
//
// This package utilizes interfaces to abstract the details of network protocols and
// transport mechanisms, allowing for easier adaptation and testing. The manager
// integrates tightly with model and protocol packages to manage the game's state
// and the serialization/deserialization of network messages.
//
// Usage of this package involves creating an instance of NetworkManager with a specific
// network transport and protocol implementation. The NetworkManager then manages the
// lifecycle of client connections and data flow throughout the game session.

import (
	"time"

	"github.com/capucinoxx/forlorn/pkg/model"
	"github.com/capucinoxx/forlorn/pkg/network"
)

const (
	// writeWait is the maximum duration to wait before timing out writes of the message.
	writeWait = 1 * time.Second

	// pongWait is the maximum time to wait for the next pong message from the client.
	pongWait = 5 * time.Second

	// pingPeriod is the duration to send pings to the client. This must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// maxMessageSize is the maximum size of a message in bytes.
	maxMessageSize = 1024
)

// Protocol is an interface to encode and decode network messages.
type Protocol interface {
	Encode(message *model.ClientMessage) []byte
	Decode(data []byte) model.ClientMessage
}

// Network is an interface for network transports. It should be capable of initializing the network,
// returning the server's address, and starting the server.
type Network interface {
	Init()
	Address() string
	Run() error
}

// NetworkManager maintains a list of clients and manages incoming and outgoing messages.
type NetworkManager struct {
	// transport holds a reference to a network.Network instance used to manage
	// low-level network operations such as opening server sockets and handling WebSocket connections.
	transport *network.Network

	// protocol defines the interface for encoding and decoding messages that are sent over the network.
	// This allows for abstraction and easy swapping of different communication protocols if needed.
	protocol Protocol

	// clients is a map that tracks the active clients connected to the server.
	// The keys are pointers to model.Client objects, and the values are booleans indicating
	// whether the client is currently active (true) or inactive (false).
	clients map[model.Connection]*model.Client

	// broadcast is a channel used to send messages to all connected clients.
	// Messages sent here are broadcasted in the network manager's main loop.
	broadcast chan []byte

	// register is a channel used for registering new clients to the server.
	// Clients are added to the network manager's client map via this channel.
	register chan *model.Client

	// unregister is a channel used for unregistering clients from the server.
	// It allows for clean removal of clients from the network manager's client map
	// and proper resource cleanup.
	unregister chan model.Connection
}

// NewNetworkManager creates a new NetworkManager with the specified network transport and
// protocol.
func NewNetworkManager(transport *network.Network, protocol Protocol) *NetworkManager {
	return &NetworkManager{
		transport:  transport,
		protocol:   protocol,
		clients:    make(map[model.Connection]*model.Client),
		broadcast:  make(chan []byte),
		register:   make(chan *model.Client),
		unregister: make(chan model.Connection),
	}
}

// Address returns the server's address.
func (nm *NetworkManager) Address() string {
	return nm.transport.Address()
}

// Start initializes the NetworkManager and starts the server.
func (nm *NetworkManager) Start() error {
	nm.transport.Init()
	go nm.run()
	return nm.transport.Run()
}

// run is the main loop of the NetworkManager. It handles incoming and outgoing messages,
// as well as registration and deregistration of clients. The game loop calls this method in
// a separate goroutine, allowing the game loop to continue running even if the network is busy.
func (nm *NetworkManager) run() {
	for {
		select {
		case c := <-nm.register:
			conn := c.GetConnection()
			nm.clients[conn] = c
			go nm.writer(c)
			if conn.Identifier() != "" {
				go nm.reader(c)
			}

		case c := <-nm.unregister:
			if client, ok := nm.clients[c]; ok {
        client.Disconnect()
				delete(nm.clients, c)

				nm.transport.Unregister(c)
			}

		case message := <-nm.broadcast:
			for conn, client := range nm.clients {
				if client.IsBlind() {
					continue
				}

				select {
				case client.Out <- message:
				default:
					nm.unregister <- conn
				}
			}
		}
	}
}

// Register adds a player to the game and also sends the current state of the game to the player.
// This method is called by the game loop when a client connects.
func (nm *NetworkManager) Register(client *model.Client) {
	nm.register <- client
}

// ForceDisconnect forcibly disconnects a client from the game.
func (nm *NetworkManager) ForceDisconnect(conn model.Connection) {
	conn.Close(writeWait, false)
	nm.unregister <- conn
}

// Send sends a message to a client.
func (nm *NetworkManager) Send(client *model.Client, message []byte) {
	client.Out <- message
}

// BroadcastGameState sends the current state of the game to all players.
// This involves sending the positions of all players and coins in the game.
func (nm *NetworkManager) BroadcastGameState(state *model.GameState, tick int32, round int8) {
	nm.broadcast <- nm.protocol.Encode(&model.ClientMessage{
		MessageType: model.MessageGameState,
		Body: model.MessageGameStateToEncode{
			CurrentTick:  tick,
			CurrentRound: round,
			Players:      state.Players(),
			Coins:        state.Coins().List(),
		},
	})
}

// BroadcastGameEnd sends a game end message to all players.
func (nm *NetworkManager) BroadcastGameEnd() {
	nm.broadcast <- nm.protocol.Encode(&model.ClientMessage{
		MessageType: model.MessageGameEnd,
	})
}

// BroadcastGameStart sends a game start message to all players.
func (nm *NetworkManager) BroadcastGameStart(state *model.GameState) {
	encodeMessage := func(isAdmin bool) []byte {
		return nm.protocol.Encode(&model.ClientMessage{
			MessageType: model.MessageMapState,
			Body: model.MessageMapStateToEncode{
				Map:     state.Map,
				IsAdmin: isAdmin,
				Storage: [100]byte{},
			},
		})
	}

	msgAdmin := encodeMessage(true)
	msg := encodeMessage(false)

	for conn, client := range nm.clients {
		var msgToSend []byte
		if conn.IsAdmin() {
			msgToSend = msgAdmin
		} else {
			msgToSend = msg
		}
		
		select {
		case client.Out <- msgToSend:
		default:
			nm.unregister <- conn
		}
	}
}

// writer writes outgoing messages to the WebSocket network. If a message cannot be written,
// the connection is closed. The game loop closes the connection in case of an error to prevent
// read and write goroutines from leaking. The game loop also closes the connection if the
// client does not respond to pings to prevent inactive connections from consuming resources.
func (nm *NetworkManager) writer(client *model.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.GetConnection().Close(writeWait, false)
	}()

	for {
		select {
		case msg, ok := <-client.Out:
			conn := client.GetConnection()
			conn.PrepareWrite(writeWait)
			if !ok {
				conn.Close(writeWait, true)
				return
			}

			if err := conn.Write(msg); err != nil {
				nm.unregister <- conn
			}

		case <-ticker.C:
			client.GetConnection().Ping(writeWait)
		}
	}
}

// reader reads incoming messages from the WebSocket network and sends them to the
// game loop. The application reads incoming messages in a separate goroutine to avoid
// blocking the game loop.
func (nm *NetworkManager) reader(client *model.Client) {
	defer func() {
		client.GetConnection().Close(writeWait, false)
		nm.unregister <- client.GetConnection()
	}()
	client.GetConnection().PrepareRead(maxMessageSize, pongWait)

	for {
		msg, err := client.GetConnection().Read()
		if err != nil {
			break
		}

		decoded := nm.protocol.Decode(msg)
		client.In <- decoded
	}
}
