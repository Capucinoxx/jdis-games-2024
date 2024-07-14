class MoveTo {
    constructor(destination) {
        if (typeof position.x !== 'number' || typeof position.y !== 'number') {
            console.error('Action "MoveTo" rejected: Expected "destination" with numeric "x" and "y" properties.');
            return null;
        }
        this.type = 'move_to';
        this.pos = destination;
    }
};


class ShootAt {
    constructor(position) {
        if (typeof position.x !== 'number' || typeof position.y !== 'number') {
            console.error('Action "ShootAt" rejected: Expected "position" with numeric "x" and "y" properties.');
            return null;
        }
        this.type = 'shoot_at';
        this.pos = position;
    }
};


class Store {
    constructor(data) {
        if (!(data instanceof Uint8Array)) {
            console.error('Action "Store" rejected: Expected "data" to be a Uint8Array.');
            return null;
        }
        this.type = 'store';
        this.data = data.slice(0, 100);
    }
};


class SwitchWeapon {
    constructor(weapon) {
        if (typeof value !== 'number') {
            console.error('Action "Switch" rejected: Expected "value" to be a number.');
            return null;
        }
        this.type = 'switch';
        this.weapon = weapon;
    }
};

export { MoveTo, ShootAt, Store, SwitchWeapon };