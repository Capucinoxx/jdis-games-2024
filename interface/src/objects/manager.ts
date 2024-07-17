import Phaser from 'phaser';
import { GameObject, Bullet, Coin, Player, Payload } from './object';
import '../types/index.d.ts';
import { CameraController } from '../lib';
import { PLAYER_WEAPON } from '../config';

interface Constructor<T> {
  new(...args: any[]): T;
};

const create_instance = <T>(ctor: Constructor<T>, ...args: any[]): T  => new ctor(...args);


class Manager<T extends GameObject> {
  protected scene: Phaser.Scene;
  private objects: Phaser.Physics.Arcade.Group;
  protected cache: Map<string, T>;
  private ctor: Constructor<T>;

  constructor(scene: Phaser.Scene, ctor: Constructor<T>) {
    this.scene = scene;
    this.objects = this.scene.physics.add.group({ classType: Phaser.GameObjects.Graphics, runChildUpdate: true });
    this.cache = new Map<string, T>();
    this.ctor = ctor;
  }

  public sync(payloads: Payload[]): void {
    const current = new Map<string, Payload>();
    payloads.forEach((p) => current.set(this.get_key(p), p));

    this.cache.forEach((obj, uuid) => {
      if (!current.has(uuid)) {
        obj.destroy();
        this.cache.delete(uuid);
      }
    });

    current.forEach((value: Payload, key: string) => {
      if (!this.cache.has(key)) {
        const obj = create_instance(this.ctor, this.scene, value);
        this.objects.add(obj);
        this.cache.set(key, obj);
        this.handle_new_entry(key);
      }
    }); 
  }

  protected handle_new_entry(id: string) {}

  public clear(): void {
    this.cache.forEach((obj, _) => obj.destroy());
    this.cache.clear();
    this.objects.clear(true, true);
  }

  private get_key(p: Payload): string {
    if ('id' in p) {
      return p.id;
    } else {
      return p.name;
    }
  }
};


class BulletManager extends Manager<Bullet> {
  constructor(scene: Phaser.Scene) {
    super(scene, Bullet);
  }

  public move(dt: number) {
    this.cache.forEach((b, _) => b.move(dt)); 
  }
};


class CoinManager extends Manager<Coin> {
  constructor(scene: Phaser.Scene) {
    super(scene, Coin);
  }
};

class PlayerManager extends Manager<Player> {
  private container: HTMLElement;
  private cam: CameraController;

  constructor(scene: Phaser.Scene, cam: CameraController) {
    super(scene, Player);
    this.container = document.querySelector('#players-list ul')!;
    this.cam = cam;
  }

  public sync(payloads: Payload[]) {
    super.sync(payloads);

    payloads.forEach((p) => {
      const pp = p as PlayerObject;
      const player = this.cache.get((p as PlayerObject).name);
      if (player) {
        player.set_movement(new Phaser.Math.Vector2(pp.pos.x, pp.pos.y), new Phaser.Math.Vector2(pp.dest!.x, pp.dest!.y));
        player.blade_visibility = (pp.current_weapon === PLAYER_WEAPON.Blade);
      }
    });
  }

  public move(dt: number) {
    this.cache.forEach((p, _) => p.move(dt));
  }

  protected handle_new_entry(id: string): void {
    const li = document.createElement('li');
    li.textContent = id;

    const target = this.cache.get(id);
    if (!target)
      return;

    li.addEventListener('click', () => {
      if (li.className === 'active') {
        li.className = '';
        this.cam.unfollow();
        return;
      }

      this.container.querySelectorAll('li').forEach((el) => el.className = '');
      li.className = 'active';
      this.cam.follow(target);
    });

    this.container.appendChild(li);
  }

  public clear(): void {
    super.clear();

    this.container.innerHTML = '';
  }
};

export { BulletManager, CoinManager, PlayerManager };

