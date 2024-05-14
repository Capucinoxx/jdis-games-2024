import Phaser from 'phaser';
import { Player } from './player';

class GameManager {
  private scene: Phaser.Scene;
  private players: Map<string, Player>;

  constructor(scene: Phaser.Scene) {
    this.scene = scene;
    this.players = new Map<string, Player>;
  }
  
  public update_from_payload(payload: any) {
    payload.players.forEach((data: any) => {
      let player = this.players.get(data.uuid);

      if (!player) {
        player = new Player(this.scene, data.x, data.y, data.uuid);
        this.scene.add.existing(player);
        this.scene.physics.add.existing(player);
        this.players.set(data.uuid, player);
      } else {
        player.update_position(data.x, data.y);
      }
    });
  }

  public update_players_movement(delta: number) {
    this.players.forEach((player) => player.move_to_target(delta));
  }

  public get_player(uuid: string): Player | undefined {
    return this.players.get(uuid);
  }
};

export { GameManager };