import Phaser from 'phaser';
import { CameraController, GameManager } from '../../lib';


class MainScene extends Phaser.Scene {
  private manager: GameManager | undefined;
  private cam: CameraController | undefined;

  constructor() { super({ key: 'MainScene' }); }

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