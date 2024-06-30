import Phaser from 'phaser';
import { PROJECTILE_SIZE, PROJECTILE_SPEED, COIN_SIZE, PLAYER_SPEED, PLAYER_SIZE } from '../config'; 
import { MovableObject } from '../lib/movable';
import '../types/index.d.ts';

type Payload = ProjectileObject | PlayerObject;

interface GameObject extends Phaser.GameObjects.Container { };


class Bullet extends MovableObject implements GameObject {
  constructor(scene: Phaser.Scene, payload: Payload) {
    const { pos, dest } = payload as ProjectileObject;
    const circle = scene.add.circle(0, 0, PROJECTILE_SIZE, 0x00ff00);

    super(scene, pos.x, pos.y, new Phaser.Math.Vector2(dest?.x, dest?.y), PROJECTILE_SPEED, [circle]);
  };
};


class Coin extends Phaser.GameObjects.Container implements GameObject {
  constructor(scene: Phaser.Scene, payload: Payload) {
    const { pos } = payload;
    const circle = scene.add.circle(0, 0, COIN_SIZE, 0xffff00);

    super(scene, pos.x, pos.y, [circle]);
  }
};


class Player extends MovableObject implements GameObject {
  constructor(scene: Phaser.Scene, payload: Payload) {
    const { pos, name, color } = payload as PlayerObject;
    const rect = scene.add.rectangle(0, 0, PLAYER_SIZE, PLAYER_SIZE, color).setOrigin(0.5, 0.5);
    const label = scene.add.text(0, -PLAYER_SIZE / 2 - 10, name, { fontSize: '16px', color: 'red' }).setOrigin(0.5, 0.5);

    super(scene, pos.x, pos.y, new Phaser.Math.Vector2(pos.x, pos.y), PLAYER_SPEED, [rect, label]);

  }

  public set_movement(pos: Phaser.Math.Vector2, dest: Phaser.Math.Vector2): void {
    this.destination = dest;
    if (Math.abs(this.x - pos.x) > 0.01 || Math.abs(this.y - pos.y) > 0.01)
      this.setPosition(pos.x, pos.y);
  }
};

export { GameObject, Bullet, Coin, Player, Payload };
