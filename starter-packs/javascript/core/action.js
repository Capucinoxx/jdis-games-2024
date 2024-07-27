class MoveAction {
    constructor(destination) {
        if (typeof destination.x !== 'number' || typeof destination.y !== 'number') {
            console.error('Action "MoveTo" rejected: Expected "destination" with numeric "x" and "y" properties.');
            return null;
        }
        this.type = 'dest';
        this.destination = destination;
    }
};


class ShootAction {
    constructor(position) {
        if (typeof position.x !== 'number' || typeof position.y !== 'number') {
            console.error('Action "ShootAt" rejected: Expected "position" with numeric "x" and "y" properties.');
            return null;
        }
        this.type = 'shoot';
        this.pos = position;
    }
};


class SaveAction {
    constructor(data) {
        if (!(data instanceof Uint8Array)) {
            console.error('Action "Store" rejected: Expected "data" to be a Uint8Array.');
            return null;
        }
        this.type = 'save';
        this.data = data.slice(0, 100);
    }
};


class SwitchWeaponAction {
    constructor(weapon) {
        if (typeof weapon !== 'number') {
            console.error('Action "Switch" rejected: Expected "weapon" to be a number.');
            return null;
        }
        this.type = 'switch';
        this.weapon = weapon;
    }
};


class BladeRotateAction {
    constructor(rad) {
        if (typeof rad !== 'number') {
            console.error('Action "BladeRotate" rejected: Expected "rad" to be a number');
            this.type = 'rotate_blade'
            this.rad = rad;
        }
    }
};

const Weapon = {
    None: 0,
    Canon: 1,
    Blade: 2
}; 

export { MoveAction, ShootAction, SaveAction, SwitchWeaponAction, BladeRotateAction, Weapon };