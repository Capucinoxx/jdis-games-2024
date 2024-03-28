import { Application, Text } from 'pixi.js';

(async () => {
  const app = new Application();
  await app.init({ background: 'red', resizeTo: window });
  document.body.appendChild(app.canvas);

  app.stage.addChild(new Text({ text: 'Hello World', x: 100, y: 100 }));
})();