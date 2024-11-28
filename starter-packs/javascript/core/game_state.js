import { Point } from './map_state.js';

class Projectile {
  constructor(uuid = '', pos = new Point(), dest = new Point()) {
    this.uuid = uuid;
    this.pos = pos;
    this.dest = dest;
  }

  toString() {
    return `{"uid":"${this.uid}","pos":${this.pos.toString()},"dest":${this.dest.toString()}}`;
  }
};

class Blade {
  constructor(start = new Point(), end = new Point(), rotation = 0.0) {
    this.start = start;
    this.end = end;
    this.rotation = rotation;
  }

  toString() {
    return `{"start":${this.start.toString()},"end":${this.end.toString()},"rotation":${this.rotation}}`;
  }
};

class Coin {
  constructor(uid = '', value = 0, pos = new Point()) {
    this.uid = uid;
    this.value = value;
    this.pos = pos;
  }

  toString() {
    return `{"uid":"${this.uid}","value":${this.value},"pos":${this.pos.toString()}}`;
  }
};

const PlayerWeapon = Object.freeze({
  PlayerWeaponNone: 0,
  PlayerWeaponCanon: 1,
  PlayerWeaponBlade: 2,
  toString(value) {
    return Object.keys(this).find(key => this[key] === value) || '';
  }
});

class PlayerInfo {
  constructor(
    name = '',
    color = 0,
    health = 0,
    score = 0,
    pos = new Point(),
    dest = null,
    player_weapon = PlayerWeapon.PlayerWeaponNone,
    projectiles = [],
    blade = null,
  ) {
    this.name = name;
    this.color = color;
    this.health = health;
    this.score = score;
    this.pos = pos;
    this.dest = dest;
    this.player_weapon = player_weapon;
    this.projectiles = projectiles;
    this.blade = blade;
  }

  isAlive() { return this.health > 0; }

  toString() {
    const projectiles_str = this.projectiles.map(p => p.toString()).join(',');
    const blade_str = this.blade ? this.blade.toString() : null;

    return `{
      "name":"${this.name}",
      "color":${this.color},
      "health":${this.health},
      "score":${this.score},
      "pos":${this.pos.toString()},
      "dest":${this.dest.toString()},
      "player_weapon":${this.player_weapon},
      "projectiles":[${projectiles_str}],
      "blade":${blade_str}
    }`;
  }
};

class GameState {
  constructor(current_tick = 0, current_round = 0, players = [], coins = []) {
    this.current_tick = current_tick;
    this.current_round = current_round;
    this.players = players;
    this.coins = coins;
  }

  toString() {
    const players_str = this.players.map(p => p.toString()).join(',');
    const coins_str = this.coins.map(c => c.toString()).join(',');

    return `{
      "current_tick":${this.current_tick},
      "current_round":${this.current_round},
      "players":[${players_str}],
      "coins":[${coins_str}]
    }`;
  }
};

export { GameState, PlayerInfo, PlayerWeapon, Coin, Blade, Projectile };
