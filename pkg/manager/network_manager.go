package manager

import (
	"time"

	"github.com/capucinoxx/forlorn/pkg/model"
	"github.com/capucinoxx/forlorn/pkg/network"
	"github.com/capucinoxx/forlorn/pkg/protocol"
	"github.com/capucinoxx/forlorn/pkg/utils"
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
	Encode(id uint8, currentGameTime uint32, message *model.ClientMessage) []byte
	Decode(data []byte) model.ClientMessage
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

	transport *network.Network

	protocol Protocol

	clients map[*model.Client]bool

	broadcast chan []byte

	register chan *model.Client

	unregister chan *model.Client
}

// NewNetworkManager crée un nouveau NetworkManager avec le transport réseau et le
// protocole spécifiés.
func NewNetworkManager(transport *network.Network, protocol Protocol) *NetworkManager {
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
			utils.Log("client", "register", "%s", c.Connection.Identifier())
			nm.clients[c] = true
			go nm.writer(c)
			go nm.reader(c)

		case c := <-nm.unregister:
			if _, ok := nm.clients[c]; ok {
				utils.Log("client", "unregister", "%s", c.Connection.Identifier())

				// TODO: gracefully disconnect client
				delete(nm.clients, c)
				nm.transport.Unregister(c.Connection)
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
// On envoi pour chaque joueur un paquet de données contenant l'identifiant du
// joueur ainsi que les données du joueur, soit:
// [0:1 id][1:2 messageType][2:6 currentTime][6:fin (position)]
func (nm *NetworkManager) BroadcastGameState(state *model.GameState) {
	players := state.Players()
	buf := make([]byte, 0, len(players)*protocol.PlayerPacketSize)

	for _, p := range players {
		buf = append(buf, nm.protocol.Encode(p.ID, 0, &model.ClientMessage{
			MessageType: model.Position,
			Body:        p,
		})...)
	}

	if len(buf) > 0 {
		nm.broadcast <- buf
	}
}

// BroadcastGameEnd envoie un message de fin de partie à tous les joueurs.
func (nm *NetworkManager) BroadcastGameEnd() {
	utils.Log("network", "broadcast", "game end")

	nm.broadcast <- nm.protocol.Encode(0, 0, &model.ClientMessage{
		MessageType: model.GameEnd,
	})
}

// BroadcastGameStart envoie un message de début de partie à tous les joueurs.
func (nm *NetworkManager) BroadcastGameStart(state *model.GameState) {
	utils.Log("network", "broadcast", "game start")

	// nm.broadcast <- nm.protocol.Encode(0, 0, &model.ClientMessage{
	// 	MessageType: model.GameStart,
	// 	Body:        state.Map.Colliders,
	// })
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

		decoded := nm.protocol.Decode(msg)
		client.In <- decoded
	}
}