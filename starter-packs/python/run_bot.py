import argparse

from network.network import Socket

def main():
    parser = argparse.ArgumentParser(description="Starts the bot")
    parser.add_argument("-t", "--token", help="The token to authenticate yout bot", required=True)
    parser.add_argument("-r", "--rank", action="store_true" ,help="If set, the bot will play ranked games")

    args = parser.parse_args()

    channel = "wss://jdis-ia.dinf.fsci.usherbrooke.ca:8088/echo"
    if args.rank:
        channel = "wss://jdis-ia.dinf.fsci.usherbrooke.ca:8087/echo"
    
    Socket(channel, args.token).run()

if __name__ == "__main__":
    main()
