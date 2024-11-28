import { Buffer } from 'buffer';
import { PlayerInfo, PlayerWeapon, Projectile, Blade, GameState, Coin } from '../core/game_state.js';
import { Point, Collider, CollderType, MapState } from '../core/map_state.js';


const read_str = (buf, offset) => {
  const idx = buf.indexOf(0, offset);
  if (idx === -1)
    return [null, offset];
  const str = buf.toString('utf8', offset, idx);
  return [str, idx + 1];
};

const read_uuid = (buf, offset) => {
  const uuid = buf.slice(offset, offset + 16);
  return uuid.toString('hex').match(/.{1,8}/g).join('-');
};

class JDISDecoder {
  decode_point(data, offset) {
    const x = data.readBoudleLE(offset);
    const y = data.readBoudleLE(offset + 8);
    return [{ x, y }, offset + 16];
  }

  decode_colliders(pos_size, data, offset) {
    c = new Collider();

    for (let i = 0; i < pos_size; i++) {
      const [pos, new_offset] = this.decode_point(data, offset);
      c.points.push(pos);
      offset = new_offset;
    }

    c.collider_type = CollderType[data.readUInt8(offset)];
    offset += 1;

    return [c, offset];
  }

  decode_map_state(data) {
    const m = new MapState();
    m.size = data.readUInt8(0);
    m.discrete_grid = Array.from({ length: size }, (_, i) => 
      Array.from(data.slice(1 + i * size, 1 + (i + 1) * size))
    );

    let offset = 1 + size * size;
    const walls_len = data.readInt32LE(offset);
    offset += 4;

    for (let i = 0; i < walls_len; i++) {
      const pos_size = data.readUInt8(offset);
      offset += 1;

      const [c, new_offset] = this.decode_colliders(pos_size, data, offset);
      m.walls.push(c);
      offset = new_offset;
    }

    m.save = data.slice(offset, offset + 100);
    return m;
  }

  decode_player_info(data) {
    const p = new PlayerInfo();
    const [name, offset] = read_str(data, 0);
    p.name = name;

    p.color = data.readUInt32LE(offset);
    p.health = data.readFloatLE(offset + 4);
    p.score = data.readFloatLE(offset + 8);
    offset += 16;

    [p.pos, offset] = this.decode_point(data, offset);

    const has_dest = data.readUInt8(offset) !== 0;
    offset += 1;

    if (has_dest) {
      const [dest, new_offset] = this.decode_point(data, offset);
      p.dest = dest;
      offset = new_offset;
    }

    p.player_weapon = PlayerWeapon[data.readUInt8(offset)];
    offset += 1;

    const projectiles_len = data.readInt32LE(offset);
    offset += 4;

    for (let i = 0; i < projectiles_len; i++) {
      const projectile = new Projectile();
      [projectile.pos, offset] = this.decode_point(data, offset);
      [projectile.dest, offset] = this.decode_point(data, offset);
      p.projectiles.push(projectile);
    }

    [p.blade.start, offset] = this.decode_point(data, offset);
    [p.blade.end, offset] = this.decode_point(data, offset);

    return [p, offset];
  }

  decode_game_state(data) {
    const g = new GameState();

    g.current_tick = data.readInt32LE(0);
    g.current_round = data.readInt8(4);
    let offset = 5;

    const player_len = data.readInt32LE(offset);
    offset += 4;

    for (let i = 0; i < player_len; i++) {
      const [player, new_offset] = this.decode_player_info(data, offset);
      g.players.push(player);
      offset = new_offset;
    }

    const coins_len = data.readInt32LE(offset);
    for (let i = 0; i < coins_len; i++) {
      const coin = new Coin();
      coin.uid = read_uuid(data, offset);
      offset += 16;

      [coin.pos, offset] = this.decode_point(data, offset);
      coin.value = data.readInt32LE(offset);
      offset += 4;

      g.coins.push(coin);
    }

    return g;
  }
};

export { JDISDecoder };
