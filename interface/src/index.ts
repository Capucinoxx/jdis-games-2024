import { Application } from 'pixi.js';
import { Vector } from './models/position';
import { Player } from './models/player';
import { MazeManager, PlayerManager } from './manager';
import { colliders } from './models/fake_maze';
import { Camera } from './models/camera';
import { WebsocketService } from './service';

const players = [
  new Player('player1', 'red', 7, new Vector(300, 300)),
];

players[0].set_destination(new Vector(500, 3000));

const ws = new WebsocketService();
ws.connect('ws://localhost:8080');

ws.subscribe((msg: ArrayBuffer): void => {
  console.log(`Receive message: ${msg.byteLength} bytes`);
});

(async () => {
  const app = new Application();
  await app.init({ resizeTo: window });
  document.body.appendChild(app.canvas);

  const maze_manager = new MazeManager(app);
  // maze_manager.set(colliders);
  maze_manager.set(await maze_manager.retrieve());

  const player_manager = new PlayerManager(app);
  player_manager.append(...players);

  const size = maze_manager.size;

  app.stage.x = 0;
  app.stage.y = 0;

  const camera = new Camera([maze_manager.view, player_manager.view], 
    players[0], 
    app.screen.width, app.screen.height,
    size.x, size.y
  );

  app.ticker.add(({ deltaTime }) => {
    player_manager.update(deltaTime);
    camera.update();
  });
  app.ticker.maxFPS = 15;

  globalThis.__PIXI_APP__ = app;
})();
