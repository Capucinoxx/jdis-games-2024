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
    };


    interface Point {
        x: number;
        y: number;
    };


    interface Projectile {
        id: string;
        pos: Point;
        dest: Point;
    };


    interface Blade {
        start: Point;
        end: Point;
        rotation: number;
    };


    interface Coin {
        id: string;
        pos: Point;
        value: number;
    };


    interface MapState {
        map: number[][];
        walls: Point[][];
    };


    interface GameState {
        tick: number;
        round: number;
        players: Player[];
        coins: Coin[];
    };


    type Action = MoveTo | ShootAt | Store | SwitchWeapon;


    type Actions = Action[];
};


declare class MoveTo {
    type: 'move_to';
    dest: Model.Point;
    constructor(destination: Model.Point);
};


declare class ShootAt {
    type: 'shoot_at';
    pos: Model.Point;
    constructor(position: Model.Point);
};


declare class Store {
    type: 'store';
    data: Uint8Array;
    constructor(data: Uint8Array);
};


declare class SwitchWeapon {
    type: 'switch';
    weapon: number;
    constructor(weapon: number);
};


declare class MyBot {
    on_start(state: Model.MapState): void;
    on_tick(state: Model.GameState): Model.Actions;
};


declare namespace Weapon {
    const None = 0;
    const Gun = 1;
    const Blade = 2;
};

