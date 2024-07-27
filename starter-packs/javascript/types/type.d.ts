
/**
 * (fr) Constantes relatives à la carte du jeu.
 * (en) Constants related to the game map.
 */
export const Consts: {
    Map: {
        /**
         * (fr) La largeur de la carte en nombre de cellules.
         * (en) The width of the map in number of cells.
         */
        WIDTH: number;
        /**
         * (fr) La hauteur de la carte en nombre de cellules.
         * (en) The height of the map in number of cells.
         */
        HEIGHT: number;
        /**
         * (fr) La largeur de chaque cellule sur la carte.
         * (en) The width of each cell in the map.
         */
        CELL_WIDTH: number;
        /**
         * (fr) La hauteur de chaque cellule sur la carte.
         * (en) The height of each cell in the map.
         */
        CELL_HEIGHT: number;
    };
    Player: {
        /**
         * (fr) La taille du joueur.
         * (en) The size of the player.
         */
        SIZE: number;
        /**
         * (fr) La vitesse du joueur (par seconde).
         * (en) The speed of the player (per second).
         */
        SPEED: number;
        /**
         * (fr) La vie maximale du joueur.
         * (en) The maximum health of the player.
         */
        MAX_HEALTH: number;
        /**
         * (fr) Le temps nécessaire pour que le joueur réapparaisse après avoir été éliminé.
         * (en) The time it takes for the player to respawn after being eliminated.
         */
        RESPAWN_TIME: number;
    };
    Projectile: {
        /**
         * (fr) La taille du projectile (rayon = SIZE / 2).
         * (en) The size of the projectile (radius = SIZE / 2).
         */
        SIZE: number;
        /**
         * (fr) La vitesse du projectile (par seconde).
         * (en) The speed of the projectile (per second).
         */
        SPEED: number;
        /**
         * (fr) Les dégâts infligés par le projectile.
         * (en) The damage dealt by the projectile.
         */
        DAMAGE: number;
        /**
         * (fr) La durée de vie du projectile (en secondes).
         * (en) The time-to-live of the projectile (in seconds).
         */
        TTL: number;
    };
    Blade: {
        /**
         * (fr) La longueur de la lame.
         * (en) The length of the blade.
         */
        LENGTH: number;
        /**
         * (fr) L'épaisseur de la lame.
         * (en) The thickness of the blade.
         */
        THICKNESS: number;
        /**
         * (fr) Les dégâts infligés par la lame.
         * (en) The damage dealt by the blade.
         */
        DAMAGE: number;
    };
    Coin: {
        /**
         * (fr) La taille d'une pièce (rayon = SIZE / 2).
         * (en) The size of the coin (radius = SIZE / 2).
         */
        SIZE: number;
        /**
         * (fr) La valeur d'une pièce lorsque ramassée.
         * (en) The value of the coin when collected.
         */
        VALUE: number;
        /**
         * (fr) La quantité de pièces dans la carte.
         * (en) The quantity of coins in the map.
         */
        QUANTITY: number;
    };
    Treasure: {
        /**
         * (fr) La taille d'un trésor (rayon = SIZE / 2).
         * (en) The size of the treasure (radius = SIZE / 2).
         */
        SIZE: number;
        /**
         * (fr) La valeur du trésor lorsque ramassé.
         * (en) The value of the treasure when collected.
         */
        VALUE: number;
    };
};



declare namespace Model {
    interface Player {
        name: string;
        color: number;
        health: number;
        pos: Point;
        dest: Point;
        current_weapon: number;
        projectiles: Projectile[];
        blade: Blade;
    }

    interface Point {
        x: number;
        y: number;
    }

    interface Projectile {
        id: string;
        pos: Point;
        dest: Point;
    }

    
    interface Blade {
        start: Point;
        end: Point;
        rotation: number;
    }


    interface Coin {
        id: string;
        pos: Point;
        value: number;
    }


    interface MapState {
        map: number[][];
        walls: Point[][];
        size: number;
        save: Uint8Array;
    }


    interface GameState {
        tick: number;
        round: number;
        players: Player[];
        coins: Coin[];
    }


    type Action = MoveTo | ShootAt | Store | SwitchWeapon;


    type Actions = Action[];
}


declare class MoveAction{
    type: 'dest';
    dest: Model.Point;
    constructor(destination: Model.Point);
}


declare class ShootAction {
    type: 'shoot';
    pos: Model.Point;
    constructor(position: Model.Point);
}


declare class SaveAction {
    type: 'save';
    data: Uint8Array;
    constructor(data: Uint8Array);
}


declare class SwitchWeaponAction {
    type: 'switch';
    weapon: number;
    constructor(weapon: number);
}

declare class BladeRotateAction {
    type: 'rotate_blade';
    rad: number;
    constructor(rad: number);
};

declare class MyBot {
    on_start(state: Model.MapState): void;
    on_end(): void;
    on_tick(state: Model.GameState): Model.Actions;
}


declare namespace Weapon {
    const None = 0;
    const Canon = 1;
    const Blade = 2;
}
