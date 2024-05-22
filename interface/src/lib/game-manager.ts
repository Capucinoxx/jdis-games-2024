import Phaser from 'phaser';
import { Player } from '.';
import { WS_URL } from './config';

class GameManager {
  private scene: Phaser.Scene;
  private players: Map<string, Player>;
  private ws: WebSocket;

  constructor(scene: Phaser.Scene) {
    this.scene = scene;
    this.players = new Map<string, Player>;
    this.ws = new WebSocket(WS_URL);
  }
  
  public update_from_payload(payload: any) {
    payload.players.forEach((data: any) => {
      let player = this.players.get(data.uuid);

      if (!player) {
        player = new Player(this.scene, {'x': data.x, 'y': data.y}, data.uuid);
        this.scene.add.existing(player);
        this.scene.physics.add.existing(player);
        this.players.set(data.uuid, player);
      } else {
        player.update_position({'x': data.x, 'y': data.y});
      }
    });
  }

  public update_players_movement(delta: number) {
    this.players.forEach((player) => player.move_to_target(delta));
  }

  public get_player(uuid: string): Player | undefined {
    return this.players.get(uuid);
  }

  private handle_ws_messages(): void {

  }
};

export { GameManager };