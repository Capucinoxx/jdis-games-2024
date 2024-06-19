import Phaser from 'phaser';
import { PROJECTILE_SIZE, PROJECTILE_SPEED, COIN_SIZE } from '../config'; 
import { MovableObject } from '../lib/movable';
import '../types/index.d.ts';

type Payload = {
  id: string;
  pos: Position;
  dest?: Position;
};

interface GameObject extends Phaser.GameObjects.Graphics { };

class Bullet extends MovableObject implements GameObject {

  constructor(scene: Phaser.Scene, payload: Payload) {
    const { pos, dest } = payload;
    super(scene, pos.x, pos.y, new Phaser.Math.Vector2(dest?.x, dest?.y), PROJECTILE_SPEED);

    this.fillStyle(0x00ff00, 1);
    this.fillCircle(0, 0, PROJECTILE_SIZE);
  };
};


class Coin extends Phaser.GameObjects.Graphics implements GameObject {
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


export { GameObject, Bullet, Coin, Payload };
