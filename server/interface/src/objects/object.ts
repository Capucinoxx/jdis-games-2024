import Phaser from 'phaser';
import { PROJECTILE_SIZE, PROJECTILE_SPEED, COIN_SIZE, PLAYER_SPEED, PLAYER_SIZE, BLADE_ROTATION_SPEED, BLADE_DISTANCE, BLADE_LENGTH } from '../config'; 
import { MovableObject } from './movable';
import '../types/index.d.ts';

type Payload = ProjectileObject | PlayerObject | ScorerObject;

interface GameObject extends Phaser.GameObjects.Container { };


class Bullet extends MovableObject implements GameObject {
  constructor(scene: Phaser.Scene, payload: Payload) {
    const { pos, dest } = payload as ProjectileObject;
    const circle = scene.add.circle(0, 0, PROJECTILE_SIZE / 2, 0x00ff00);

    super(scene, pos.x, pos.y, new Phaser.Math.Vector2(dest?.x, dest?.y), PROJECTILE_SPEED, [circle]);

    this.setDepth(4);
  };
};


class Coin extends Phaser.GameObjects.Container implements GameObject {
  public static size = COIN_SIZE;

  constructor(scene: Phaser.Scene, payload: Payload) {
    const { pos } = payload as ScorerObject;
    const circle = scene.add.circle(0, 0, Coin.size / 2, 0xdee0e9);
    super(scene, pos.x, pos.y, [circle]);

    this.scene.physics.world.enable(this);
    scene.add.existing(this);

    this.setDepth(6);
  }
};


class Player extends MovableObject implements GameObject {
  private blade: Blade;

  constructor(scene: Phaser.Scene, payload: Payload) {
    const { pos, name, color } = payload as PlayerObject;
    const rect = scene.add.rectangle(0, 0, PLAYER_SIZE, PLAYER_SIZE, color).setOrigin(0.5, 0.5);
    const label = scene.add.text(0, -PLAYER_SIZE / 2 - 10, name, { fontSize: '16px', color: 'white' }).setOrigin(0.5, 0.5);

    super(scene, pos.x, pos.y, new Phaser.Math.Vector2(pos.x, pos.y), PLAYER_SPEED, [rect, label]);

    this.blade = new Blade(scene, this, color);
    this.setDepth(5);
  }

  public set_movement(pos: Phaser.Math.Vector2, dest: Phaser.Math.Vector2): void {
    this.destination = dest;

    if (Math.abs(this.x - pos.x) > 0.01 || Math.abs(this.y - pos.y) > 0.01) {
      this.setPosition(pos.x, pos.y);
      this.blade.setPosition(pos.x, pos.y);
    }
  }

  public rotate_blade(theta: number): void {
    this.blade.rotate(theta);
  }

  public update(time: number, delta: number): void {
    super.update(time, delta);
    this.blade.update();
  }

  public destroy(): void {
    this.blade.destroy();
    super.destroy(); 
  }

  protected on_arrival(): void {}
};

class Blade extends Phaser.GameObjects.Container {
  private owner: Player;

  constructor(scene: Phaser.Scene, player: Player, color: number) {
    super(scene, player.x, player.y);
    this.owner = player;

    const graphics = scene.add.graphics();
    graphics.fillStyle(color, 1);
    graphics.fillRect(-BLADE_LENGTH / 2, -2, BLADE_LENGTH, 4);

    this.add(graphics);
    
    scene.add.existing(this);
  }

  public update(): void {
    this.x = this.owner.x + (BLADE_DISTANCE * Math.cos(this.angle));
    this.y = this.owner.y + (BLADE_DISTANCE * Math.sin(this.angle));
  }

  public rotate(theta: number): void {
    this.rotation = theta;
  }
};

export { GameObject, Bullet, Coin, Player, Payload };
