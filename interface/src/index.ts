import { Application, Text, Graphics, Point, Container } from 'pixi.js';
import { Vector } from './models/position';
import { Player } from './models/player';


const players = [
  new Player('player1', 'red', 7, new Vector(300, 300)),
];

players[0].set_destination(new Vector(100, 100));

(async () => {
  const app = new Application();
  await app.init({ background: 'blue', resizeTo: window });
  document.body.appendChild(app.canvas);

  const players_container = new Container();


  players.forEach(player => players_container.addChild(player.graphics));

  app.stage.addChild(players_container);

  app.ticker.add((t) => {
    players.forEach(player => player.update(t.deltaTime));
  });

  app.ticker.maxFPS = 15;
})();