# JDIS Games 2024 - Magellan

Programming Agents Competition organized as part of the 2024 edition of the JDIS Games at the University of Sherbrooke.

## Installation

### Domain Modification

To modify the domain, follow these steps:

1. Open your command line interface.

2. Set the new base URL:

    ```bash
    base_url="<NEW URL>"
    ```

3. Change the Docker Compose URL:

    ```bash
    sed -i.bak  -e "s|API_URL=.*/rank|API_URL=${base_url}/rank|g" \
                -e "s|API_URL=.*/unrank|API_URL=${base_url}/unrank|g" \
                -e "s|DOMAIN=.*|DOMAIN=${base_url}|g" \
                "docker-compose.yml"
    ```

4. Change the action URL in the HTML file:

    ```bash
    sed -i      -e "s|action='*/create|action='http://${base_url}/create'|g" \
                "server/interface/index.html"
    ```

### Create and Modify the `.env` File

1. Copy the example `.env` file:

    ```sh
    cp .env.example .env
    ```

2. Modify the administrative information in the `.env` file as needed.

## Starting the Services

To start the services, run:

```sh
docker compose up --build
```

Then navigate to:
- `https://<YOUR URL>/rank`
- `https://<YOUR URL>/unrank`

## Creating an Agent

To create an agent, follow these steps:

1. Go to the game page at https://<YOUR URL>/unrank.
2. Click on the menu icon in the top right corner.
3. Enter the name of your agent.
4. Copy the token from the message at the bottom right.
5. Use this token in your starter pack.

## Administrato Actions

Administrators can perform the following actions:

| Action                        | Path                                                               |
| :---------------------------- | :----------------------------------------------------------------- |
| Start a game                  | `https://<URL>/<rank,unrank>/start?tkn=<ADMIN_TOKEN>`              |
| Toggle leaderboard visibility | `https://<URL>/<rank,unrank>/toggle_leaderboard?tkn=<ADMIN_TOKEN>` |
| Freeze a game                 | `https://<URL>/<rank,unrank>/freeze?tkn=<ADMIN_TOKEN>`             |
| Unfreeze a game               | `https://<URL>/<rank,unrank>/unfreeze?tkn=<ADMIN_TOKEN>`           |

