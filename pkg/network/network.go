package network

// Package network fournit une abstraction pour la gestion des connexions réseau
// en utilisant le protocole WebSocket pour une communication bidirectionnelle en
// temps réel entre le client et le serveur. Cela permet l'enregistrement et la
// désinscription de connexions, la mise à niveau des requêtes HTTP en connexions
// WebSocket, et la gestion des lectures et écritures de messages.
//
// Ce package utilise le package "github.com/gorilla/websocket" pour la gestion
// des connexions WebSocket.

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

// Network représente un serveur réseau capable de gérer des connexions WebSocket.
type Network struct {
	// upgrader configure les paramètres pour la mise à niveau des requêtes HTTP en
	// connexions WebSocket.
	upgrader websocket.Upgrader

	// register est une fonction appelée lorsqu'une nouvelle connexion est établie.
	register func(conn model.Connection)

	// uregister est une fonction appelée lorsqu'une connexion est fermée.
	uregister func(conn model.Connection)

	// port est le port sur lequel le serveur écoute.
	port int

	// address est l'adresse IP sur laquelle le serveur écoute.
	address string

	connected sync.Map
}

// NewNetwork crée un nouveau serveur réseau avec l'adresse et le port spécifiés.
func NewNetwork(address string, port int) *Network {
	return &Network{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		port:      port,
		address:   address,
		connected: sync.Map{},
	}
}

// Address retourne l'adresse IP et le port du serveur.
func (n *Network) Address() string {
	return strings.Join([]string{n.address, strconv.Itoa(n.port)}, ":")
}

// SetRegisterFunc définit la fonction lorsqu'une nouvelle connexion est établie.
func (n *Network) SetRegisterFunc(f func(conn model.Connection)) {
	n.register = f
}

// SetUnregisterFunc définit la fonction lorsqu'une connexion est fermée.
func (n *Network) SetUnregisterFunc(f func(conn model.Connection)) {
	n.uregister = f
}

// Register enregistre une nouvelle connexion en appelant la fonction de registre
// spécifiée.
func (n *Network) Register(conn model.Connection) {
	if n.register != nil {
		n.register(conn)
	}
}

// Unregister désenregistre une connexion en appelant la fonction de désenregistrement
// spécifiée.
func (n *Network) Unregister(conn model.Connection) {
	if n.uregister != nil {
		n.uregister(conn)
	}

	n.connected.Delete(conn.Identifier())
}

// Init initialise le serveur réseau en écoutant les requêtes HTTP.
func (n *Network) Init() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "" {
			if _, ok := n.connected.Load(token); ok {
				http.Error(w, "error expected", http.StatusUnauthorized)
				return
			}

			n.connected.Store(token, true)
		}

		ws, _ := n.upgrader.Upgrade(w, r, nil)
		n.register(NewConnection(ws, token))
	})
}

// Run démarre le serveur réseau en écoutant les connexions entrantes sur l'adresse
// IP et le port spécifiés. La méthode retourne une erreur si le serveur ne peut
// pas démarrer.
func (n *Network) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", n.address, n.port))
	if err != nil {
		return err
	}
	n.port = listener.Addr().(*net.TCPAddr).Port

	return http.Serve(listener, nil)

}

// Connection est une implémentation de l'interface model.Connection pour les
// connexions WebSocket. Elle encapsule un pointeur vers une connexion WebSocket
// et fournit des méthodes pour la lecture, l'écriture, la fermeture et le ping
// de la connexion.
type Connection struct {
	// conn est la connexion WebSocket.
	conn *websocket.Conn

	token string
}

// NewConnection crée une nouvelle connexion à partir d'une connexion WebSocket.
func NewConnection(conn *websocket.Conn, token string) *Connection {
	return &Connection{conn: conn, token: token}
}

// Identifier retourne une chaîne de caractères représentant l'adresse IP et le
// port de la connexion.
func (c *Connection) Identifier() string {
	return c.token
}

// Close ferme la connexion en envoyant un message de fermeture et en fermant la
// connexion.
func (c *Connection) Close(writeWait time.Duration) {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})

}

// PrepareRead configure la connexion pour la lecture en définissant la taille
// maximale des messages, le délai d'attente pour le pong, et le gestionnaire de
// pong. Cette méthode doit être appelée avant chaque lecture de message.
func (c *Connection) PrepareRead(maxMessageSize int64, pongWait time.Duration) {
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}

// Read lit un message de la connexion et retourne le message lu ou une erreur.
// Cette méthode doit être appelée après l'appel à la méthode PrepareRead.
func (c *Connection) Read() ([]byte, error) {
	_, msg, err := c.conn.ReadMessage()
	return msg, err
}

// PrepareWrite configure la connexion pour l'écriture en définissant le délai
// d'attente pour l'écriture. Cette méthode doit être appelée avant chaque écriture
// de message.
func (c *Connection) PrepareWrite(writeWait time.Duration) {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
}

// Write écrit un message dans la connexion et retourne une erreur si l'écriture
// échoue. Cette méthode doit être appelée après l'appel à la méthode PrepareWrite.
func (c *Connection) Write(msg []byte) error {
	writer, err := c.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}

	writer.Write(msg)
	return writer.Close()
}

// Ping envoie un message de ping à la connexion pour vérifier si elle est toujours
// active. Cette méthode doit être appelée régulièrement pour maintenir la connexion
// active.
func (c *Connection) Ping(writeWait time.Duration) {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	_ = c.conn.WriteMessage(websocket.PingMessage, nil)
}
