import { ArgumentParser } from 'argparse';

const DEFAULT_BASE_URL = 'ws://localhost:8086/echo';



const parser = new ArgumentParser({ description: 'Starts the bot' });
parser.add_argument('-s', '--secret', { help: 'The secret that authenticates your bot', required: true });
parser.add_argument('-r', '--rank', { help: 'If set, the bot will play ranked game', action: 'store_true' });
parser.add_argument('-u', '--url', { help: 'The url of the game server', default: DEFAULT_BASE_URL });

const { url, rank, secret } = parser.parse_args();

const channel = `${DEFAULT_BASE_URL}`;
