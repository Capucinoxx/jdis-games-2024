package model

type MessageType uint8

const (
	// Spawn est un message envoyé pour informer le client
	// qu'il vient de rejoindre la partie.
	Spawn MessageType = iota

	// Position est un message utilisé pour notifié soit du déplacement d'un joueur,
	// soit de la rotation d'un joueur lorsque cela est une communication client -> serveur.
	// Quand la communication est serveur -> client, ce message est utilisé pour
	// informer le client de la position des joueurs.
	Position

	// Register est un message envoyé pour enregistrer un joueur à la partie.
	Register

	// Projectile est un
	Projectile

	// GameStart est un message envoyé pour informer les clients que la partie commence.
	GameStart

	// GameEnd est un message envoyé pour informer les clients que la partie est terminée.
	GameEnd
)
