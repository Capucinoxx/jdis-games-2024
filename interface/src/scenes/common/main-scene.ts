import Phaser from 'phaser';
import { GameManager } from '../../lib';


class MainScene extends Phaser.Scene {
  private manager: GameManager | undefined;
  private cursors: Phaser.Types.Input.Keyboard.CursorKeys | undefined;

  constructor() { super({ key: 'MainScene' }); }

  create() {
    this.manager = new GameManager(this);

    const cam = this.cameras.main;
    cam.setZoom(1);
    
    this.cursors = this.input.keyboard?.createCursorKeys();

    this.input.on('wheel', (_: Phaser.Input.Pointer, __: any, ___: number, dy: number, ____: number) => {
      if (dy > 0) cam.zoom = Math.max(0.5, cam.zoom - 0.1);
      else if (dy < 0) cam.zoom = Math.min(1.1, cam.zoom + 0.1);
    });
  }

  private handle_input(dt: number) {
    if (!this.cursors)
      return;
  
    const cam = this.cameras.main;
    const speed = 500;

    if (this.cursors!.left.isDown) cam.scrollX -= speed * (dt / 1000);
    else if (this.cursors!.right.isDown) cam.scrollX += speed * (dt / 1000);

    if (this.cursors!.up.isDown) cam.scrollY -= speed * (dt / 1000);
    else if (this.cursors!.down.isDown) cam.scrollY += speed * (dt / 1000);
  }

  update(_: number, delta: number) {
    this.manager!.update_players_movement(delta);
    this.handle_input(delta);
  }
};

export { MainScene };
