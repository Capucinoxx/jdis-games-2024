import Phaser from 'phaser';
import { MainScene } from './scenes';
import { load_wasm } from './wasm/loader';
import { handle_modal_interraction } from './modal';
import { handle_forms } from './form';

handle_modal_interraction();
handle_forms();

const config: Phaser.Types.Core.GameConfig = {
  type: Phaser.AUTO,
  width: window.innerWidth,
  height: window.innerHeight,
  backgroundColor: '#F0F0F0',
  physics: {
    default: 'arcade',
    arcade: {
      debug: false
    }
  },
  scene: [MainScene],
  scale: {
    mode: Phaser.Scale.RESIZE,
    autoCenter: Phaser.Scale.CENTER_BOTH
  }
};

load_wasm().then(() => new Phaser.Game(config));
