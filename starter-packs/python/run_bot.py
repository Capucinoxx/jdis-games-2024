import argparse

import asyncio

DEFAULT_BASE_URL = "ws://jdis-ia.dinf.fsci.usherbrooke.ca"

# async def loop(token: str, url: str):
#     await Socket(url, token).run()

def main():
    parser = argparse.ArgumentParser(description="Starts the bot")
    parser.add_argument("-u", "--url", help="The base URL of the server", default=DEFAULT_BASE_URL)
    parser.add_argument("-t", "--token", help="The token to authenticate yout bot", required=True)
    parser.add_argument("-r", "--rank", action="store_true" ,help="If set, the bot will play ranked games")

    args = parser.parse_args()

    args.url = args.url.rstrip()
    channel = "/unranked/game"
    if args.rank:
        channel = "/ranked/game"

    print(f"Starting bot with base URL: {args.url}, token: {args.token}, is ranked: {args.rank}")
    
    # asyncio.run(loop(args.token, args.url + channel))

if __name__ == "__main__":
    main()