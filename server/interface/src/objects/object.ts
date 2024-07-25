import Phaser from 'phaser';
import { PROJECTILE_SIZE, PROJECTILE_SPEED, COIN_SIZE, PLAYER_SPEED, PLAYER_SIZE, BLADE_ROTATION_SPEED, BLADE_DISTANCE, BLADE_LENGTH } from '../config'; 
import { MovableObject } from './movable';
import '../types/index.d.ts';

type Payload = ProjectileObject | PlayerObject | ScorerObject;

interface GameObject extends Phaser.GameObjects.Container { };


class Bullet extends MovableObject implements GameObject {
  constructor(scene: Phaser.Scene, payload: Payload) {
    const { pos, dest } = payload as ProjectileObject;
    const circle = scene.add.circle(0, 0, PROJECTILE_SIZE / 2, 0xffffff);

    super(scene, pos.x, pos.y, new Phaser.Math.Vector2(dest?.x, dest?.y), PROJECTILE_SPEED, [circle]);

    this.setDepth(4);
  };

  public set_movement(pos: Phaser.Math.Vector2, dest: Phaser.Math.Vector2): void {
    this.destination = dest;

    if (Math.abs(this.x - pos.x) > 0.01 || Math.abs(this.y - pos.y) > 0.01) {
      this.setPosition(pos.x, pos.y);
    }
  }
};


class Coin extends Phaser.GameObjects.Container implements GameObject {
  public static size = COIN_SIZE;

  constructor(scene: Phaser.Scene, payload: Payload) {
    const { pos } = payload as ScorerObject;
    const img = scene.add.image(0, 0, Coin.size === COIN_SIZE ? 'coin' : 'big-coin');
    img.setOrigin(0.5, 0.5);
    img.setDisplaySize(Coin.size, Coin.size);

    super(scene, pos.x, pos.y, [img]);

    this.scene.physics.world.enable(this);
    scene.add.existing(this);

    this.setDepth(3);
  }
};


class Player extends MovableObject implements GameObject {
  private blade: Blade;
  private custom_color: number;

  constructor(scene: Phaser.Scene, payload: Payload) {
    const { pos, name, color } = payload as PlayerObject;

    const circle_size = PLAYER_SIZE / 2 * 0.6;
    const circle = scene.add.circle(0, 0, circle_size, color).setOrigin(0.5, 0.5);

    const img = scene.add.image(0, 0, 'agent').setOrigin(0.5, 0.5);
    img.setDisplaySize(PLAYER_SIZE, PLAYER_SIZE);
    const label = scene.add.text(0, -PLAYER_SIZE / 2 - 10, name, { fontSize: '16px', color: 'white' }).setOrigin(0.5, 0.5);

    super(scene, pos.x, pos.y, new Phaser.Math.Vector2(pos.x, pos.y), PLAYER_SPEED, [img, circle, label]);

    this.blade = new Blade(scene, this, color);
    this.custom_color = color;
    this.setDepth(5);
  }

  public set_movement(pos: Phaser.Math.Vector2, dest: Phaser.Math.Vector2): void {
    this.destination = dest;

    if (Math.abs(this.x - pos.x) > 0.01 || Math.abs(this.y - pos.y) > 0.01) {
      this.setPosition(pos.x, pos.y);
      this.blade.setPosition(pos.x, pos.y);
    }
  }

  public get color(): string {
    const r = (this.custom_color >> 16) & 0xFF;
    const g = (this.custom_color >> 8) & 0xFF;
    const b = this.custom_color & 0xFF;
    return `rgba(${r}, ${g}, ${b}, 1)`
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
  private graphics: Phaser.GameObjects.Graphics;

  constructor(scene: Phaser.Scene, player: Player, color: number) {
    super(scene, player.x, player.y);
    this.owner = player;

    this.graphics = scene.add.graphics();
    this.graphics.fillStyle(color, 1);
    this.graphics.fillRect(-BLADE_LENGTH / 2, -2, BLADE_LENGTH, 4);

    this.add(this.graphics);
    scene.add.existing(this);
  }

  public update(): void {
    this.x = this.owner.x + (BLADE_DISTANCE * Math.cos(this.rotation));
    this.y = this.owner.y + (BLADE_DISTANCE * Math.sin(this.rotation));
  }

  public rotate(theta: number): void {
    this.setRotation(-theta);
  }
};

export { GameObject, Bullet, Coin, Player, Payload };
