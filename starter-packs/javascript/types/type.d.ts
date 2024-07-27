
/**
 * Model contient l'ensemble de types
 * utilisés pour représenter l'état du jeu.
 */
declare namespace Model {
    /**
    * Représente un agent (bot) dans le jeu.
    * Chaque agent est unique par son nom.
    *
    * Chaque bot possède deux armes soit un pistolet et une épée.
    * Lorsqu'un bot tire, un projectile est créé et ajouté à la liste des projectiles.
    * L'attribut current_weapon indique l'arme actuellement utilisée (0 = aucun, 1 = pistolet, 2 = épée).
    */
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

    /**
    * Point représente une position dans la map. 
    * La position (0, 0) est le coin en haut à gauche.
    */
    interface Point {
        x: number;
        y: number;
    }


    /**
    * Chaque joueur possède une liste de projectiles. 
    * Un projectile est tiré par le pistolet du joueur.
    */
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


declare class MoveTo {
    type: 'move_to';
    dest: Model.Point;
    constructor(destination: Model.Point);
}


declare class ShootAt {
    type: 'shoot_at';
    pos: Model.Point;
    constructor(position: Model.Point);
}


declare class Store {
    type: 'store';
    data: Uint8Array;
    constructor(data: Uint8Array);
}


declare class SwitchWeapon {
    type: 'switch';
    weapon: number;
    constructor(weapon: number);
}

declare class BladeRotate {
    type: 'rotate_blade';
    rad: number;
    constructor(rad: number);
};

declare class MyBot {
    on_start(state: Model.MapState): void;
    on_tick(state: Model.GameState): Model.Actions;
}


declare namespace Weapon {
    const None = 0;
    const Canon = 1;
    const Blade = 2;
}
