import Phaser from 'phaser';
import { MainScene } from './scenes';
import { load_wasm } from './wasm/loader';
import { handle_modal_interraction } from './modal';
import { handle_forms } from './form';

handle_modal_interraction();
handle_forms();

const config: Phaser.Types.Core.GameConfig = {
  type: Phaser.AUTO,
  width: 800,
  height: 800,
  backgroundColor: '#F0F0F0',
  pixelArt: true,
  physics: {
    default: 'arcade',
    arcade: {
      debug: false
    }
  },
  scene: [MainScene],
  scale: {
    autoCenter: Phaser.Scale.CENTER_BOTH
  }
};

load_wasm().then(() => new Phaser.Game(config));
