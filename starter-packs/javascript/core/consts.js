/**
 * Constantes relatives à la carte du jeu.
 * Constants related to the game map.
 */
const Map = {
    WIDTH: 10,       // La largeur de la carte en nombre de cellules. / The width of the map in number of cells.
    HEIGHT: 10,      // La hauteur de la carte en nombre de cellules. / The height of the map in number of cells.
    CELL_WIDTH: 10,  // La largeur de chaque cellule sur la carte. / The width of each cell in the map.
    CELL_HEIGHT: 10  // La hauteur de chaque cellule sur la carte. / The height of each cell in the map.
};

/**
 * Constantes relatives au joueur.
 * Constants related to the player.
 */
const Player = {
    SIZE: 1.0,          // La taille du joueur. / The size of the player.
    SPEED: 1.15,        // La vitesse du joueur (par seconde). / The speed of the player (per second).
    MAX_HEALTH: 100,    // La vie maximale du joueur. / The maximum health of the player.
    RESPAWN_TIME: 5.0   // Le temps nécessaire pour que le joueur réapparaisse après avoir été éliminé. / The time it takes for the player to respawn after being eliminated.
};

/**
 * Constantes relatives aux projectiles.
 * Constants related to projectiles.
 */
const Projectile = {
    SIZE: 0.35,         // La taille du projectile (rayon = SIZE / 2). / The size of the projectile (radius = SIZE / 2).
    SPEED: 3.0,         // La vitesse du projectile (par seconde). / The speed of the projectile (per second).
    DAMAGE: 15,         // Les dégâts infligés par le projectile. / The damage dealt by the projectile.
    TTL: 5.0            // La durée de vie du projectile (en secondes). / The time-to-live of the projectile (in seconds).
};

/**
 * Constantes relatives aux lames.
 * Constants related to blades.
 */
const Blade = {
    LENGTH: 2.0,        // La longueur de la lame. / The length of the blade.
    THICKNESS: 0.25,    // L'épaisseur de la lame. / The thickness of the blade.
    DAMAGE: 4           // Les dégâts infligés par la lame. / The damage dealt by the blade.
};

/**
 * Constantes relatives aux pièces.
 * Constants related to coins.
 */
const Coin = {
    SIZE: 0.5,          // La taille d'une pièce (rayon = SIZE / 2). / The size of the coin (radius = SIZE / 2).
    VALUE: 40,          // La valeur d'une pièce lorsque ramassée. / The value of the coin when collected.
    QUANTITY: 30        // La quantité de pièces dans la carte. / The quantity of coins in the map.
};

/**
 * Constantes relatives au trésor.
 * Constants related to treasure.
 */
const Treasure = {
    SIZE: 4.0,          // La taille d'un trésor (rayon = SIZE / 2). / The size of the treasure (radius = SIZE / 2).
    VALUE: 1200         // La valeur du trésor lorsque ramassé. / The value of the treasure when collected.
};

const Consts = {
    Map,
    Player,
    Projectile,
    Blade,
    Coin,
    Treasure
};

export { Consts };
