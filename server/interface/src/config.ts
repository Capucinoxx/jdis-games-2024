const WS_URL = `wss://${window.location.hostname}${window.location.pathname}/echo`;
let TYPE = window.location.href.endsWith('unrank') ? 'unrank' : 'rank';


const MESSAGE_TYPE = {
  1: 'POSITION',
  4: 'GAME_START',
  5: 'GAME_END'
};

const PLAYER_WEAPON = {
  'None': 0,
  'Canon': 1,
  'Blade': 2
};

const SCALE = 30;

const TICK_ROUND_TWO_START = 4 * 60 * 3;

const PLAYER_SIZE = 1 * SCALE;
const PLAYER_SPEED = 1.15 * SCALE;

const PROJECTILE_SIZE = 0.35 * SCALE;
const PROJECTILE_SPEED = 3.0 * SCALE;

const COIN_SIZE = 0.5 * SCALE;
const BIG_COIN_SIZE = 4.0 * SCALE;
const COIN_SIZES = [COIN_SIZE, BIG_COIN_SIZE];

const NUM_CELLS = 10;
const CELL_WIDTH = 10.0 * SCALE;


const BLADE_ROTATION_SPEED = 230;
const BLADE_DISTANCE = 1.5 * SCALE;
const BLADE_LENGTH = 1.5 * SCALE;

export {
  PLAYER_WEAPON, PLAYER_SIZE, PLAYER_SPEED,
  PROJECTILE_SIZE, PROJECTILE_SPEED,
  BLADE_ROTATION_SPEED, BLADE_DISTANCE, BLADE_LENGTH,
  COIN_SIZE, BIG_COIN_SIZE, COIN_SIZES,
  CELL_WIDTH, NUM_CELLS,
  TICK_ROUND_TWO_START,
};

export { WS_URL, MESSAGE_TYPE, TYPE as GAME_TYPE };
