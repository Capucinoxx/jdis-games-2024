import argparse

from network.network import Socket

def main():
    parser = argparse.ArgumentParser(description="Starts the bot")
    parser.add_argument("-t", "--token", help="The token to authenticate yout bot", required=True)
    parser.add_argument("-r", "--rank", action="store_true" ,help="If set, the bot will play ranked games")

    args = parser.parse_args()

    args.url = args.url.rstrip()

    channel = "wss://localhost:8088/echo"
    if args.rank:
        channel = "wss://localhost:8087/echo"

    print(f"Starting bot with base URL: {args.url}, token: {args.token}, is ranked: {args.rank}")
    
    Socket(args.url + channel, args.token).run()

if __name__ == "__main__":
    main()
