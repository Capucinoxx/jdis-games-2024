import Phaser from 'phaser';
import { CameraController, GameManager } from '../../lib';
import { URL_BASE } from '../../config';


class MainScene extends Phaser.Scene {
  private manager: GameManager | undefined;
  private cam: CameraController | undefined;

  constructor() { super({ key: 'MainScene' }); }

  preload() {
    this.load.image('coin', `https://${URL_BASE}/assets/coin.png`);
    this.load.image('big-coin', `https://${URL_BASE}/assets/big-coin.png`);
    this.load.image('agent', `https://${URL_BASE}/assets/agent.png`);
  }

  create() {
    this.cam = new CameraController(this.cameras.main, this.input);
    this.manager = new GameManager(this, this.cam);
    this.manager.generate_admin_form(document.body);
  }

  update(_: number, delta: number) {
    this.manager!.update_players_movement(delta);
    this.cam!.update(delta);
  }
};

export { MainScene };
