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
  };
};


class Coin extends Phaser.GameObjects.Container implements GameObject {
  public static size = COIN_SIZE;

  constructor(scene: Phaser.Scene, payload: Payload) {
    const { pos } = payload as ScorerObject;
    const circle = scene.add.circle(0, 0, Coin.size / 2, 0x131313);
    super(scene, pos.x, pos.y, [circle]);

    this.scene.physics.world.enable(this);
    scene.add.existing(this);
  }
};


class Player extends MovableObject implements GameObject {
  private blade: Blade;

  constructor(scene: Phaser.Scene, payload: Payload) {
    const { pos, name, color } = payload as PlayerObject;
    const rect = scene.add.rectangle(0, 0, PLAYER_SIZE, PLAYER_SIZE, color).setOrigin(0.5, 0.5);
    const label = scene.add.text(0, -PLAYER_SIZE / 2 - 10, name, { fontSize: '16px', color: 'red' }).setOrigin(0.5, 0.5);

    super(scene, pos.x, pos.y, new Phaser.Math.Vector2(pos.x, pos.y), PLAYER_SPEED, [rect, label]);

    this.blade = new Blade(scene, this);

  }

  public set blade_visibility(visibility: boolean) {
    const curr_visibility = this.blade.visibility;
    if (curr_visibility != visibility) {
      this.blade.visible = visibility;
      this.blade.visibility = visibility;
    }
  }

  public set_movement(pos: Phaser.Math.Vector2, dest: Phaser.Math.Vector2): void {
    this.destination = dest;

    if (Math.abs(this.x - pos.x) > 0.01 || Math.abs(this.y - pos.y) > 0.01) {
      this.setPosition(pos.x, pos.y);
      this.blade.setPosition(pos.x, pos.y);
    }
  }

  public update(time: number, delta: number): void {
    super.update(time, delta);


    this.blade.update(delta);
  }

  public destroy(): void {
    this.blade.destroy();
    super.destroy(); 
  }
};

class Blade extends Phaser.GameObjects.Container {
  private blade: Phaser.GameObjects.Rectangle;
  private owner: Player;
  private speed: number;
  public visibility: boolean = false; 

  constructor(scene: Phaser.Scene, player: Player) {
    super(scene, player.x, player.y);
    const blade = scene.add.rectangle(0, 0, 4, BLADE_LENGTH, 0x0000);
    this.visible = false;

    this.add(blade);

    this.blade = blade;
    this.owner = player;
    this.speed = Phaser.Math.DegToRad(BLADE_ROTATION_SPEED);
    
    scene.add.existing(this);
  }

  public update(dt: number): void {
    if (!this.visibility)
      return;

    this.angle += (this.speed * (dt / 1000));
    this.angle %= (Math.PI * 2);

    this.x = this.owner.x + (BLADE_DISTANCE * Math.cos(this.angle));
    this.y = this.owner.y + (BLADE_DISTANCE * Math.sin(this.angle));

    const dx = this.owner.x - this.x;
    const dy = this.owner.y - this.y;
    const rotation = Math.atan2(dy, dx);

    this.blade.setRotation(rotation + Math.PI / 2);
  }
};

export { GameObject, Bullet, Coin, Player, Payload };
