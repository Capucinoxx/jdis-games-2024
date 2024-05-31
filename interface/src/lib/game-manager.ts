import Phaser from 'phaser';
import { Player } from '.';
import { WS_URL, MESSAGE_TYPE } from './config';
import { GridManager } from './grid-manager'; 

class GameManager {
  private scene: Phaser.Scene;
  private players: Map<string, Player>;
  private ws: WebSocket;
  private grid: GridManager;

  constructor(scene: Phaser.Scene) {
    this.scene = scene;
    this.grid = new GridManager(scene);
    this.players = new Map<string, Player>;

    this.ws = new WebSocket(WS_URL);
    this.ws.binaryType = 'arraybuffer';

    this.handle_ws_messages();
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
    this.ws.onmessage = (event: MessageEvent<ArrayBuffer>) => {
      const data = window.getInformations(event.data);
      if (!('type' in data))
        return;

      switch (data.type) {
        case 4:
         this.grid.tiles = data.map;    
      }
      //console.log(window.getInformations(event.data));
    }

    this.ws.onclose = (event: CloseEvent) => {
      console.log('Disconnected from the server', event);
    };
  }
};

export { GameManager };
