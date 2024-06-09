import Phaser from 'phaser';
import { Player } from '.';
import { WS_URL } from './config';
import { GridManager } from './grid-manager';
import { BulletManager } from './bullet-manager';
import '../types/index.d.ts';

class GameManager {
  private scene: Phaser.Scene;
  private players: Map<string, Player>;
  private ws: WebSocket;
  private grid: GridManager;
  private bullets: BulletManager; 

  constructor(scene: Phaser.Scene) {
    this.scene = scene;
    this.grid = new GridManager(scene);
    this.bullets = new BulletManager(scene);
    this.players = new Map<string, Player>;

    this.ws = new WebSocket(WS_URL);
    this.ws.binaryType = 'arraybuffer';

    this.handle_ws_messages();
  }

  public handle_game_state(payload: ServerGameState): void {
    const payload_bullets: Projectile[] = [];

    payload.players.forEach((data: PlayerData) => {
      let player = this.players.get(data.name);
      payload_bullets.push(...data.projectiles);

      if (player) {
        player.set_movement(
          new Phaser.Math.Vector2(data.pos.x * 30, data.pos.y * 30),
          new Phaser.Math.Vector2(data.dest.x * 30, data.dest.y * 30));
        return;
      }

      player = new Player(this.scene, data.pos.x * 30, data.pos.y * 30, data.name, 0x7f7287);
      this.scene.add.existing(player);
      this.scene.physics.add.existing(player);
      this.players.set(data.name, player);
    });

    this.bullets.sync(payload_bullets); 
  }

  public update_players_movement(delta: number) {
    this.players.forEach((player) => player.move(delta));
  }

  public get_player(uuid: string): Player | undefined {
    return this.players.get(uuid);
  }

  private handle_ws_messages(): void {
    this.ws.onmessage = (event: MessageEvent<ArrayBuffer>) => {
      const data = window.getInformations(event.data);
      if (!('type' in data))
        return;

      switch (data.type) {
        case 4:
         this.grid.tiles = data.map;
          break;
        case 1:
          console.log(data);
          this.handle_game_state(data);
          break;
      }
    }

    this.ws.onclose = (event: CloseEvent) => {
      console.log('Disconnected from the server', event);
    };
  }
};

export { GameManager };
