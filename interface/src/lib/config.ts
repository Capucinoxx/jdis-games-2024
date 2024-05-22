const WS_URL = 'ws://localhost:8087/echo';

const MESSAGE_TYPE = {
  0: 'SPAWN',
  1: 'POSITION',
  2: 'REGISTER',
  3: 'FIRE',
  4: 'GAME_START',
  5: 'GAME_END'
};

export { WS_URL, MESSAGE_TYPE };