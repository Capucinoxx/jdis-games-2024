import Phaser from 'phaser';
import { GameManager } from '../../lib';


class MainScene extends Phaser.Scene {
  private manager: GameManager | undefined;

  constructor() { super({ key: 'MainScene' }); }

  create() {
    this.manager = new GameManager(this);

    this.cameras.main.setBounds(0, 0, 8000, 6000);
    this.physics.world.setBounds(0, 0, 8000, 6000);

    this.input.on('wheel', (pointer: Phaser.Input.Pointer): void => {
      const dy = pointer.deltaY;
      
      let new_zoom = this.cameras.main.zoom + (dy > 0 ? -.1 : .1);
      if (new_zoom > 0.3 && new_zoom <= 1)
        this.cameras.main.zoom = new_zoom;
    });

    this.input.on('pointerdown', (pointer: Phaser.Input.Pointer): void => {
      this.cameras.main.pan(pointer.x, pointer.y, 500);
    });

    window.addEventListener('resize', () => this.scale.resize(window.innerWidth, window.innerHeight));
  }

  update(_: number, delta: number) {
    this.manager!.update_players_movement(delta);
  }
};

export { MainScene };
