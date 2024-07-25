# Magellan - JDIS Games 2024

## Mise en situation

Imaginez-vous au début du 16e siècle, une époque où les cartes du monde étaient encore remplies de zones inconnues et mystérieuses. Au milieu de cette période de découvertes et d'explorations, Ferdinand Magellan, un explorateur portugais, se préparait. Son objectif : trouver une route maritime vers les îles aux épices dont les richesses étaient convoitées par toute l'Europe.

À bord de son navire, Magellan et son équipage s'aventurèrent dans des eaux inexplorées, armés de patience et de leurs connaissances en navigation. Leur principal allié était l'astrolabe, un instrument capable de mesurer l'altitude des étoiles et des planètes. Cet outil était indispensable pour tracer leur chemin à travers les vastes océans, leur permettant de se frayer un chemin vers l'inconnu.

## Votre objectif

À l'image de Magellan, vous deviendrez des explorateurs naviguant sur des eaux inconnues à la recherche d'un trésor. Vous connaissez sa coordonnée, mais l'emplacement des obstacles reste floue. Il y a aura deux phases dans votre aventure : découverte et prise.

Dans Magellan, vous contrôlez un agent de forme circulaire pouvant se déplacer sur une carte. Vous devrez vous déplacer sur la carte afin de ramasser les pièces et les trésors se trouvant dans votre chemin tout en vous défendant contre les autres agents. 

## Déroulement du jeu

Le jeu est sous format _"Long running"_ ce qui signifie que le jeu n'arrête jamais. Plusieurs parties auront lieu eu courant de toute la journée. Vous accumulerez des points **tout au long de la journée**. Vous devez prévoir les meilleurs moments pour déconnecter votre agent pour y mettre une nouvelle version. Il est primordiale de concevoir son bot de manière incrémentale et de faire de l’amélioration continue.

Un cycle de rafraichissement dure 300 ms.Durant un cycle de rafraichissement, le serveur effectuera 10 boucles d'action. 

Lors d'une partie, l'agent sera placé sur une carte et recevra les informations des [éléments](#éléments-de-la-carte) se trouvant sur cette dernière. L'agent pourra envoyer plusieurs [actions](#actions) par cycle de rafraichissement sous certaines conditions.

À chaque nouvelle partie, tous les murs et les pièces sont placés de manière aléatoire sur la carte. Une partie est composé de plusieurs phases :
- [Phase 1 : Découverte](#phase-1)
- [Phase 2 : Prise](#phase-2)

## Éléments de la carte

Lors des précédentes explorations, certaines informations ont été receuilles.

Les éléments suivants se trouveront sur la carte :   

|        |                                             |
| ------ | ------------------------------------------- |
| Pièce  | ![image](/starter-packs/docs/coin.png)      |
| Trésor | ![image](/starter-packs/docs/astrolab2.png) |
| Mur    | Les murs ne sont pas visibles sur la carte  |

### Phase 1

Lors de la première phase de la partie, les agents connectés au serveur ainsi que le pièces seront placés de manière aléatoire sur la carte. 🗺️ Ceci est la phase de découverte et d'exploration de la carte. 

Lorsqu'un agent prend une pièce, celle-ci reaparrait de manière aléatoire sur la carte.
Lors de la fin de la phase, les agents et les pièces sont enlevés de la carte pour passer à la deuxième phase.

### Phase 2

Lors de la dernière phase de la partie, un trésor sera placé au centre de la carte, les agents apparaitront à équidistance de déplacement du trésor. Dans cette phase, il n'y aura pas de pièces sur la carte.

### Mort

Lorsque l'agent perd toute sa vie, ce dernier disparait de la carte et ne reçoit aucune donnée du serveur pendant un temps défini.

## Actions

Plusieurs actions peuvent être envoyées au serveur dans un même cycle de raifraichissement avec certaines contraintes d'utilisation.

### Déplacement

La position envoyée est celle vers laquele l'agent va se déplacer. L'agent ne peut pas traverser les murs. Il n'y a pas d'algorithme de recherche de chemin (pathfinding) qui a implémenté, Pour vous rendre à une position précise, vous devrez implémenter un algorithme de recherche de chemin (pathfinding).

Cette action n'a pas de contraintes d'utilisation.

###  Attaque

#### Changement d'arme

Par défaut, vous n'êtes équipez d'aucune arme. Pour pouvoir utiliser une arme, vous devrez faire l'action de changement d'arme.

Cette action ne peut pas être accompagnée de l'utilisation d'une arme dans un même cycle de rafraichissement.

#### Pistolet

Pour utiliser le pistolet, il faut envoyer la position de destination souhaitée d'un projectile. Le projectile a une portée définie. Lorsqu'un projectile rentre en collision avec un autre agent, ce dernier reçoit du dégât. Le projectile disparaît par la suite.

Cette action ne peut pas être accompagnée de l'équipement d'une arme dans le même cycle de rafraichissement. 

#### Lame

En début de partie, la lame apparaitra à 0 radians du joueur comme illustré à l'image suivante.  

<div align="center">
  <img width="200" alt="logo" src="./docs/blade2.png">
</div>

Pour changer la rotation de la lame, il faut envoyer la nouvelle rotation en radians. Lorsque la lame rentre en collision avec un agent, celui-ci reçoit maximum 40 points de dégât. Le nombre de points de dégâts dependra de la durée que la lame sera en contact avec l'agent durant un cycle de raffraichissement.

### Sauvegarde

Une quantité limitée d'octects pourra être envoyée au serveur. Ces octets seront sauvegardé le temps d'une partie. Ces données seront reçues par le joueur à chaque fois qu'il va se connecter au serveur. Cette action vous permettra donc de sauvegardé de l'information dont vous aurait accès dans une même partie lorsque votre bot sera déconnecté.

## Données reçues

### État de la carte

L'agent reçoit l'état de carte lorsqu'il se connecte au serveur. 

```
MapState {
    DiscretGrid : [][]Integer,      // Grille contenant le nombre de murs par 4 cases
    size: int,                      // La taille de la carte
    walls: List[Collider],          // La liste des murs, cette liste devrait être vide
    save: bytearray                 // La mémoire que vous pouvez sauvegardé dans le jeu
}
```

### Grille discrète

La carte est reçue sous forme de grille discrète. La grille discrète contient seulement le nombre de murs par 4 cases. Les murs extérieurs détimitant la carte sont aussi comptés. La grille est envoyé à chaque début de partie dans l'état de la carte. La guille ne change pas au long d'une partie. Les murs d'une grillés sont placés de manière aléatoire à chaque nouvelle  partie.

Voici un exemple de grille :

<div align="center">
  <img width="1000" alt="logo" src="./docs/grille_murs.png">
</div>

### Pointage

| Item       | Points |
| ---------- | ------ |
| Pistolet   | 15     |
| Lame       | 40     |
| Pièce      | 40     |
| Trésor     | 1200   |

### État du jeu

L'agent reçoit l'état du jeu à chaque cycle de rafraichissement.

```
GameState {
    currentTick : Integer,                          // Cycle de rafraichissement courant
    currentRound : Integer,                         // Tour courant
    Players : [                                     // Liste des informations des joueurs 
        name : String,                              // Nom du joueur
        color : String,                             // Couleur du joueur sur la carte
        health : Integer,                           // Quantité de vie du joueur
        score : Integer,                            // Nombre de points du joueur
        position : { x : Float, y : Float },        // Position actuelle
        destination : { x : Float, y : Float },     // Position de destination
        playerWeapon : Integer,                     // Type d'arme équipée, aucune(0), pistolet(1), épée(2) 
        Projectiles: [                              // Informations des projectiles du pistolet
            uid : String,                           // Identifiant du projectile
            position : { x : Float, y : Float },    // Position actuelle
            destination : { x : Float, y : Float }, // Position de destination
        ],
        Blade: {                                    // Informations sur l'épée
            start : { x : Float, y : Float },       // Position du début de l'épée
            end : { x : Float, y : Float },         // Position de fin de l'épée
            rotation : Integer                      // Angle en radians
        }
    ],
    Coins : [                                       // Liste des pièces ou trésor
        uid : String,                               // Identifiant de la pièce
        value : Integer,                            // Valeur de la pièce ou du trésor
        position : { x : Float, y : Float }         // Position de la pièce ou du trésor
    ]
}
```




## Constantes

AJOUTER TABLEAU