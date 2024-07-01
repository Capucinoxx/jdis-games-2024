import Phaser from 'phaser';
import { GameManager } from '../../lib';


class MainScene extends Phaser.Scene {
  private manager: GameManager | undefined;
  private cursors: Phaser.Types.Input.Keyboard.CursorKeys | undefined;
  private wasd_keys: { w: Phaser.Input.Keyboard.Key, a: Phaser.Input.Keyboard.Key, s: Phaser.Input.Keyboard.Key, d: Phaser.Input.Keyboard.Key };

  constructor() { super({ key: 'MainScene' }); }

  create() {
    this.manager = new GameManager(this);

    const cam = this.cameras.main;
    cam.setZoom(1);
   
    if (!this.input.keyboard)
      return;

    this.cursors = this.input.keyboard.createCursorKeys();
    this.wasd_keys = {
      w: this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.W),
      a: this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.A),
      s: this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.S),
      d: this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.D),
    };

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

    if (this.cursors.left.isDown || this.wasd_keys.a.isDown) cam.scrollX -= speed * (dt / 1000);
    else if (this.cursors.right.isDown || this.wasd_keys.d.isDown) cam.scrollX += speed * (dt / 1000);

    if (this.cursors.up.isDown || this.wasd_keys.w.isDown) cam.scrollY -= speed * (dt / 1000);
    else if (this.cursors!.down.isDown || this.wasd_keys.s.isDown) cam.scrollY += speed * (dt / 1000);
  }

  update(_: number, delta: number) {
    this.manager!.update_players_movement(delta);
    this.handle_input(delta);
  }
};

export { MainScene };
