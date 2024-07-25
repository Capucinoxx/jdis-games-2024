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
  protected handle_remove_entry(key: string): void { }

  private get_key(p: Payload): string {
    return ('id' in p) ? p.id : p.name;
  }

  public sync(payloads: Payload[]): void {
    this.curr_cache = new Map(this.next_cache);

    this.cache.forEach((obj, uuid) => {
      if (!this.curr_cache.has(uuid)) {
        obj.destroy(true);
        this.cache.delete(uuid);
        this.handle_remove_entry(uuid);
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

  private calculate_path(curr: ProjectileObject, next: ProjectileObject | undefined): { pos: Position, dest: Position } {
    if (!next || (curr.pos.x === next.pos.x && curr.pos.y === next.pos.y))
      return { pos: curr.pos, dest: curr.pos };
    return { pos: curr.pos, dest: next.dest };
  }

  public sync(payloads: Payload[]) {
    super.sync(payloads);

    payloads.forEach((p) => {
      const bullet = this.cache.get((p as ProjectileObject).id);
      const curr_bullet = this.curr_cache.get((p as ProjectileObject).id) as ProjectileObject | undefined;
      const next_bullet = this.next_cache.get((p as ProjectileObject).id) as ProjectileObject | undefined;

      if (bullet && curr_bullet) {
        const { pos, dest } = this.calculate_path(curr_bullet, next_bullet);
        bullet.set_movement(new Phaser.Math.Vector2(pos.x, pos.y), new Phaser.Math.Vector2(dest.x, dest.y));
      }
    });
  }

  public move(dt: number) { this.cache.forEach((b, _) => b.move(dt)); }
};


class CoinManager extends Manager<Coin> {
  constructor(scene: Phaser.Scene) { super(scene, Coin); }
};


class PlayerManager extends Manager<Player> {
  private cam: CameraController;
  private container: HTMLElement;
  private filter: HTMLInputElement;

  constructor(scene: Phaser.Scene, cam: CameraController) {
    super(scene, Player);

    this.cam = cam;
    
    const btn = document.querySelector('#players-list')! as HTMLElement;
    btn.addEventListener('click', () => btn.classList.toggle('focus'));

    this.container = document.querySelector('#players-list ul')!;
    this.filter = this.container.querySelector('input')!;
    this.filter.addEventListener('click', (e) => { e.preventDefault(); e.stopPropagation(); });
    this.filter.addEventListener('input', this.hande_filter_input.bind(this));
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

  protected handle_remove_entry(id: string) {
    const els = this.container.children;
    for (let i = 0; i < els.length; i++) {
      const el = els.item(i) as HTMLElement | null;
      if (!el) continue;

      if (el.classList.contains('actiobe')) this.cam.unfollow();
      if (el.textContent === id) el.remove();
    }
  }

  private hande_filter_input(e: Event) {
    e.preventDefault();
    e.stopPropagation();
    const filter = this.filter.value.toLowerCase();
    const els = this.container.children;

    for (let i = 1; i < els.length; i++) {
      const el = els.item(i) as HTMLElement;

      if (el.textContent === '' || el.textContent?.toLowerCase().includes(filter)) {
        el.style.display = '';
      } else {
        el.style.display = 'none';
      }
    }
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
    const li = document.createElement('li');
    li.append(this.filter);
    this.container.append(li);
  }

  public get(name: string): Player | undefined {
    return this.cache.get(name);
  }
};

export { BulletManager, CoinManager, PlayerManager };
