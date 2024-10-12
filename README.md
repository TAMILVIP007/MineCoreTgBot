# MineCoreBot

A simple bot for Minecraft servers that allows you to manage the server and interact with players through Telegram.

## Features

- Start/Stop the server
- Execute commands on the server
- Receive chat messages from the server in Telegram
- Send chat messages from Telegram to the server
- View live console output from the server

## Installation

1. Clone the repository:

    ```bash
    git clone github.com/tamilvip007/MineCoreBot
    ```

2. Create a Telegram bot.
3. Create a Telegram group.
4. Add the bot to the group.
5. Get the group ID.
6. Create a configuration file:

    ```bash
    nano .env
    ```

    Add the following content to the file:

    ```bash
    RCONN_PORT=12345
    RCONN_HOST=1.1.1.1
    RCONN_PASSWORD=12345
    BOT_TOKEN=123214
    CHAT_TOPIC=4
    DEATH_MSG_TOPIC=6
    CHAT_ID=-10023
    ACHIEVEMENT_TOPIC=10
    LOG_FILE_PATH=/root/path/to/server.log
    OWNER_ID=1234
    EVENT_TOPIC=8
    DEFAULT_TOPIC=12
    ```

7. Run the bot:

    ```bash
    go run main.go
    ```

## TODO

- [ ] Add more events
- [ ] Improve handling of death messages; it’s currently broken
- [ ] Implement a way to read only new logs written to the log file instead of reading the whole file
- [ ] Support multiple servers

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
