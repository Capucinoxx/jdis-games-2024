import Phaser from 'phaser';
import { PROJECTILE_SIZE, PROJECTILE_SPEED, COIN_SIZE } from '../config'; 
import '../types/index.d.ts';

type Payload = {
  id: string;
  pos: Position;
  dest?: Position;
};


class Bullet extends Phaser.GameObjects.Graphics {
  private destination: Phaser.Math.Vector2;

  constructor(scene: Phaser.Scene, payload: Payload) {
    super(scene);
    this.scene = scene;
    const { pos, dest } = payload;

    this.destination = new Phaser.Math.Vector2(dest?.x, dest?.y);

    this.setPosition(pos.x, pos.y);
    this.fillStyle(0x00ff00, 1);

    this.fillCircle(0, 0, PROJECTILE_SIZE);
    this.scene.add.existing(this);
    this.scene.physics.world.enable(this);
  };

  public update(): void {
    const angle = Phaser.Math.Angle.Between(this.x, this.y, this.destination.x, this.destination.y);
    (this.body as Phaser.Physics.Arcade.Body).setVelocity(Math.cos(angle) * PROJECTILE_SPEED, Math.sin(angle) * PROJECTILE_SPEED);

    if (Phaser.Math.Distance.Between(this.x, this.y, this.destination.x, this.destination.y) < PROJECTILE_SIZE)
      this.destroy();
  }
};


class Coin extends Phaser.GameObjects.Graphics {
  constructor(scene: Phaser.Scene, payload: Payload) {
    super(scene);
    this.scene = scene;
    const { pos } = payload;


    this.setPosition(pos.x, pos.y);
    this.fillStyle(0xffff00, 1);

    this.fillCircle(0, 0, COIN_SIZE);
    this.scene.add.existing(this);
    this.scene.physics.world.enable(this);
  }
};


export { Bullet, Coin, Payload };
