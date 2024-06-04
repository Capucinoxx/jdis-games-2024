interface Position { x: number, y: number };

type PlayerData = {
  name: string;
  health: number;
  pos: Position;
  projectiles: Array<{ pos: Position, dest: Position }>
};

type ServerMapState = {
  type: 4;
  map: Array<Array<number>>;
};

type ServerGameState = {
  type: 1;
  players: Array<PlayerData>;
};

type Empty = Record<string, never>;


type ServerMessage = ServerMapState | ServerGameState | Empty;
