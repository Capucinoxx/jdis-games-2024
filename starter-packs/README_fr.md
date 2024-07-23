# Magellan - JDIS Games 2024

## Mise en situation

Ferdinand Magellan était un explorateur portugais connu pour avoir dirigé la première expédition à faire le tour du monde au 16e siècle. Le but principal de l'expédition de Magellan était de trouver une route maritime vers les îles aux épices. À cette époque, les épices étaient extrêmement précieux en Europe. Les astrolabes étaient des instruments de navigation essentiels utilisés par explorateurs afin dde mesurer l'altitude des étoiles et des planètes au-dessus de l'horizon pour ainsi résoudre des problèmes de navigation.

## Votre objectif

Tout comme Magellan, vous serez des explorateurs devant naviguer des terres inconnues pour s'approprier d'un précieux trésor. Vous connaissez sa coordonnée, mais l'emplacement des obstacles reste floue. Il y a aura deux phases dans votre aventure : exploration et ...

## Éléments de la carte

Les éléments suivants se trouveront sur la carte : 
- Murs 🧱
- Pièces 🪙
- Trésor

### Pointage

| Item       | Points |
| ---------- | ------ |
| Projectile | 15     |
| Pièce      | 40     |
| Épée       | 40     |
| Trésor     | 1200   |

## Déroulement du jeu

### Tour 1

Durée : x temps
Les agents et les pièces sont placées de manière alétoire sur la carte 🗺️

Lorsqu'un agent prend une pièce, celle-ci reaparrait de manière aléatoire sur la carte.

### Tour 2

Durée : x temps
Un trésor est placé au centre de la carte, les agents apparaitront à équidistance de déplacement du trésor. Dans ce tour, il n'y aura pas de pièces sur la carte.

## Actions

### Déplacement



###  Attaque

#### Pistolet



#### Épée

### Sauvegarde

## Données reçues

### État de la carte

L'état de carte est envoyée lors de la connexion au serveur. 


#### Grille discrète

La grille discrète contient le nombre de murs par 4 cases.

METTRE IMAGE DE LA MAP AVEC MURS ET SANS MURS

### État du jeu

```
GameState {
    DiscretGrid : [][]Number    // Grille contenant le nombre de murs par 4 cases
    Players : [                 // Tableau contenant les informations des joueurs 
        name : String,
        color : String
    ] ...
}
```