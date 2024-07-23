import Phaser from 'phaser';
import { Bullet, Coin, GameObject, Payload, Player } from './object';
import { CameraController } from '../lib';


interface Constructor<T> { new(...args: any[]): T; };

const create_instance = <T>(ctor: Constructor<T>, ...args: any[]): T => new ctor(...args);

class Manager<T extends GameObject> {
  protected scene: Phaser.Scene;
  private objects: Phaser.Physics.Arcade.Group;
  private ctor: Constructor<T>;

  protected cache: Map<string, T>;
  protected curr_cache: Map<string, Payload>;
  protected next_cache: Map<string, Payload>;

  constructor(scene: Phaser.Scene, ctor: Constructor<T>) {
    this.scene = scene;
    this.objects = this.scene.physics.add.group({ classType: Phaser.GameObjects.Graphics, runChildUpdate: true });
    this.ctor = ctor;

    this.cache = new Map<string, T>();
    this.curr_cache = new Map<string, Payload>();
    this.next_cache = new Map<string, Payload>();
  }

  protected handle_new_entry(key: string): void { }

  private get_key(p: Payload): string {
    return ('id' in p) ? p.id : p.name;
  }

  public sync(payloads: Payload[]): void {
    this.curr_cache = new Map(this.next_cache);
    console.log(this.cache, this.curr_cache);


    this.cache.forEach((obj, uuid) => {
      if (!this.curr_cache.has(uuid)) {
        console.log(`delete ${uuid}`);
        obj.destroy(true);
        this.cache.delete(uuid);
      }
    });

    this.curr_cache.forEach((value: Payload, key: string) => {
      if (!this.cache.has(key)) {
        const obj = create_instance(this.ctor, this.scene, value);
        this.objects.add(obj);
        this.cache.set(key, obj);
        this.handle_new_entry(key);
      }
    });

    this.next_cache.clear();
    payloads.forEach((p) => this.next_cache.set(this.get_key(p), p));
  }

  public clear(): void {
    this.cache.forEach((obj, _) => obj.destroy());
    this.cache.clear();
    this.curr_cache.clear();
    this.next_cache.clear();
    this.objects.clear(true, true);
  }
};


class BulletManager extends Manager<Bullet> {
  constructor(scene: Phaser.Scene) { super(scene, Bullet); }
  public move(dt: number) { this.cache.forEach((b, _) => b.move(dt)); }
};


class CoinManager extends Manager<Coin> {
  constructor(scene: Phaser.Scene) { super(scene, Coin); }
};


class PlayerManager extends Manager<Player> {
  private cam: CameraController;
  private container: HTMLElement;

  constructor(scene: Phaser.Scene, cam: CameraController) {
    super(scene, Player);

    this.cam = cam;
    this.container = document.querySelector('#players-list ul')!;
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
      this.cam.follow(target, id);
    });

    if (localStorage.getItem('follow') === id) {
      this.container.querySelectorAll('li').forEach((el) => el.className = '');
      li.className = 'active';
      this.cam.follow(target, id);
    }

    this.container.appendChild(li);
  }

  private calculate_path(curr: PlayerObject, next: PlayerObject | undefined): { pos: Position, dest: Position } {
    if (!next || (curr.pos.x === next.pos.x && curr.pos.y === next.pos.y))
      return { pos: curr.pos, dest: curr.pos };
    return { pos: curr.pos, dest: next.dest };
  }

  public sync(payloads: Payload[]) {
    super.sync(payloads);

    payloads.forEach((p) => {
      const player = this.cache.get((p as PlayerObject).name);
      const curr_player = this.curr_cache.get((p as PlayerObject).name) as PlayerObject | undefined;
      const next_player = this.next_cache.get((p as PlayerObject).name) as PlayerObject | undefined;

      if (player && curr_player) {
        const { pos, dest } = this.calculate_path(curr_player, next_player);
        player.set_movement(new Phaser.Math.Vector2(pos.x, pos.y), new Phaser.Math.Vector2(dest.x, dest.y));
        player.rotate_blade(curr_player.blade.rotation);
      }
    });
  }

  public move(dt: number) {
    this.cache.forEach((p, _) => p.move(dt));
  }

  public clear(): void {
    super.clear();

    this.container.innerHTML = '';
  }

  public get(name: string): Player | undefined {
    return this.cache.get(name);
  }
};

export { BulletManager, CoinManager, PlayerManager };
