package manager

import (
	"time"

	"github.com/capucinoxx/forlorn/model"
	"github.com/capucinoxx/forlorn/network"
)

const (
	// writeWait est le délai maximum pour écrire un message sur la connexion.
	writeWait = 1 * time.Second

	// pongWait est le délai maximum pour lire le prochain pong du client.
	pongWait = 5 * time.Second

	// pingPeriod envoie des pings au client avec ce délai. Doit être inférieur à pongWait.
	pingPeriod = (pongWait * 9) / 10

	// maxMessageSize est la taille maximale d'un message en octets.
	maxMessageSize = 1024
)

// Protocol est une interface pour encoder et décoder des messages réseau.
type Protocol interface {
	Encode(id int, currentGameTime uint32, message []byte) []byte
	Decode(data []byte) []byte
}

// Network est une interface pour les transports réseau. Il doit être capable
// d'initialiser le réseau, de retourner l'adresse du serveur et de démarrer
// le serveur.
type Network interface {
	Init()
	Address() string
	Run() error
}

// NetworkManager maintient une liste de clients et gère les messages entrants et sortants.
type NetworkManager struct {
	Ready bool

	transport network.Network

	protocol Protocol

	clients map[*model.Client]bool

	broadcast chan []byte

	register chan *model.Client

	unregister chan *model.Client
}

// NewNetworkManager crée un nouveau NetworkManager avec le transport réseau et le
// protocole spécifiés.
func NewNetworkManager(transport network.Network, protocol Protocol) *NetworkManager {
	return &NetworkManager{
		Ready:      false,
		transport:  transport,
		protocol:   protocol,
		clients:    make(map[*model.Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *model.Client),
		unregister: make(chan *model.Client),
	}
}

// Address retourne l'adresse du serveur.
func (nm *NetworkManager) Address() string {
	return nm.transport.Address()
}

// Start initialise le NetworkManager et démarre le serveur.
func (nm *NetworkManager) Start() error {
	nm.transport.Init()
	go nm.run()
	return nm.transport.Run()
}

// run est la boucle principale du NetworkManager. Il gère les messages entrants
// et sortants, ainsi que l'inscription et la désinscription des clients. La
// boucle de jeu appelle cette méthode dans une goroutine séparée. Cela permet
// à la boucle de jeu de continuer à s'exécuter même si le réseau est occupé.
func (nm *NetworkManager) run() {
	for {
		select {
		case c := <-nm.register:
			nm.clients[c] = true
			go nm.writer(c)
			go nm.reader(c)

		case c := <-nm.unregister:
			if _, ok := nm.clients[c]; ok {
				// TODO: gracefully disconnect client
				delete(nm.clients, c)
				nm.transport.Register(c.Connection)
			}

		case message := <-nm.broadcast:
			for client := range nm.clients {
				select {
				case client.Out <- message:
				default:
					nm.unregister <- client
				}
			}
		}
	}
}

// Register ajoute un joueur à la partie. Cela envoie également l'état actuel
// de la partie au joueur. Cette méthode est appelée par la boucle de jeu lorsqu'un
// client se connecte.
func (nm *NetworkManager) Register(player *model.Player) {
	// TODO: aussi envoyer l'état de la map au joueur

	nm.register <- player.Client
}

// ForceDisconnect déconnecte un joueur de la partie.
func (nm *NetworkManager) ForceDisconnect(player *model.Player) {
	client := player.Client
	client.Connection.Close(writeWait)
	nm.unregister <- player.Client
}

// Send envoie un message à un client.
func (nm *NetworkManager) Send(client *model.Client, message []byte) {
	client.Out <- message
}

// BroadcastGameState envoie l'état actuel de la partie à tous les joueurs.
func (nm *NetworkManager) BroadcastGameState() {
	// TODO: envoyer l'état de la map à tous les joueurs
}

// writer écrit les messages sortants dans le réseau WebSocket. Si un message
// ne peut pas être écrit, la connexion est fermée. La boucle de jeu ferme la
// connexion en cas d'erreur. Cela empêche les goroutines de lecture et
// d'écriture de fuir. La boucle de jeu ferme également la connexion si le
// client ne répond pas aux pings. Cela empêche les connexions inactives de
// consommer des ressources.
func (nm *NetworkManager) writer(client *model.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Connection.Close(writeWait)
	}()

	for {
		select {
		case msg, ok := <-client.Out:
			client.Connection.PrepareWrite(writeWait)
			if !ok {
				client.Connection.Close(writeWait)
				return
			}

			client.Connection.Write(msg)

		case <-ticker.C:
			client.Connection.Ping(writeWait)
		}
	}
}

// reader lit les messages entrants dans le réseau WebSocket et les envoie à la
// boucle de jeu. L'application lit les messages entrants dans une goroutine
// séparée pour éviter de bloquer la boucle de jeu.
func (nm *NetworkManager) reader(client *model.Client) {
	defer func() {
		client.Connection.Close(writeWait)
		nm.unregister <- client
	}()
	client.Connection.PrepareRead(maxMessageSize, pongWait)

	for {
		msg, err := client.Connection.Read()
		if err != nil {
			break
		}
		client.In <- nm.protocol.Decode(msg)
	}
}
