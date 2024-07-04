interface Position { x: number, y: number };

type ServerObject = {
  id: string;
  pos: Position;
};

type BladeObject = {
  start: Position;
  end: Position;
  rotation: number;
};

type ProjectileObject = ServerObject & {
  dest: Position;
};

type ScorerObject = ServerObject & {
  value: number;
};

type PlayerObject = {
  name: string;
  color: number;
  pos: Position;
  dest: Position;
  blade: BladeObject;
};

type PlayerData = {
  name: string;
  color: number;
  health: number;
  pos: Position;
  dest: Position;
  blade: BladeObject;
  projectiles: Array<Projectile>
};

type ServerMapState = {
  type: 4;
  map: Array<Array<number>>;
  walls: Array<Array<Position>>;
};

type ServerGameEnd = {
  type: 5;
};

type ServerGameState = {
  type: 1;
  tick: number;
  round: number;
  players: Array<PlayerData>;
  coins: Array<ScorerObject>;
};

type Empty = Record<string, never>;


type ServerMessage = ServerMapState | ServerGameState | ServerGameEnd | Empty;
