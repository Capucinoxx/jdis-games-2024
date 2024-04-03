import { Application } from 'pixi.js';
import { Vector } from './models/position';
import { Player } from './models/player';
import { PlayerManager } from './manage/player_manager';
import { MazeManager } from './manage/maze_manager';
import { colliders } from './models/fake_maze';

const players = [
  new Player('player1', 'red', 7, new Vector(300, 300)),
];

players[0].set_destination(new Vector(100, 100));

(async () => {
  const app = new Application();
  await app.init({ background: 'blue', resizeTo: window });
  document.body.appendChild(app.canvas);

  const maze_manager = new MazeManager(app);
  maze_manager.set(colliders);

  const player_manager = new PlayerManager(app);
  player_manager.append(...players);

  app.ticker.add(({ deltaTime }) => player_manager.update(deltaTime));
  app.ticker.maxFPS = 15;
})();
