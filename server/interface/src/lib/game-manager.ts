import Phaser from 'phaser';
import { COIN_SIZES, WS_URL } from '../config';
import { GridManager } from './grid-manager';
import { BulletManager, CoinManager, PlayerManager } from '../objects';
import '../types/index.d.ts';
import { Coin } from '../objects/object';
import { CameraController } from './camera-controller';
import { ProgressBar } from '../progress';

class GameManager {
  private ws: WebSocket | undefined;
  private grid: GridManager;
  
  private players: PlayerManager;
  private bullets: BulletManager;
  private coins: CoinManager;

  private progress: ProgressBar

  constructor(scene: Phaser.Scene, cam: CameraController) {
    this.grid = new GridManager(scene);
    this.bullets = new BulletManager(scene);
    this.coins = new CoinManager(scene);
    this.players = new PlayerManager(scene, cam);
    this.progress = new ProgressBar(document.querySelector('#game rect') as SVGRectElement);

    this.ws_connection = '';
  }

  public set ws_connection(token: string) {
    const conn = WS_URL + (token === '' ? '' : `?token=${token}`);
    if (this.ws && this.ws.readyState !== WebSocket.CLOSED)
      this.ws.close();

    this.ws = new WebSocket(conn);
    this.ws.binaryType = 'arraybuffer';

    this.handle_ws_messages();
}

  public generate_admin_form(container: HTMLElement) {
    const input = document.createElement('input');
    const btn = document.createElement('button');
    btn.textContent = 'reload';

    btn.addEventListener('click', (e) => {
      e.preventDefault();
      this.ws_connection = input.value;      
      input.value = '';
    });

    container.appendChild(input);
    container.appendChild(btn);
  }

  public handle_game_state(payload: ServerGameState): void {
    const payload_bullets: ProjectileObject[] = [];
    const payload_players: PlayerObject[] = [];

    payload.players.forEach((data: PlayerData) => {
      const { name, pos, color, dest, blade, current_weapon } = data;
      payload_players.push({ name, pos, color, dest, blade, current_weapon });
      payload_bullets.push(...data.projectiles);
    });

    this.progress.current_value = payload.tick;
    this.players.sync(payload_players);
    this.bullets.sync(payload_bullets);
    this.coins.sync(payload.coins);
  }

  public update_players_movement(delta: number) {
    this.players.move(delta);
    this.bullets.move(delta);
  }


  private handle_ws_messages(): void {
    if (!this.ws)
      return;

    this.ws.onmessage = (event: MessageEvent<ArrayBuffer>) => {
      const data = window.getInformations(event.data);
      if (!('type' in data))
        return;

      switch (data.type) {
        case 4:
          this.grid.map = { cells: data.map, colliders: data.walls };
          break;

        case 5:
          this.clean();
          break;

        case 1:
          Coin.size = COIN_SIZES[data.round];
          this.handle_game_state(data);
          break;
      }
    }

    this.ws.onclose = (event: CloseEvent) => {
      console.log('Disconnected from the server', event);
    };
  }

  private clean(): void {
    this.players.clear();
    this.bullets.clear();
    this.coins.clear();
    this.grid.clear();
  }
};

export { GameManager };
