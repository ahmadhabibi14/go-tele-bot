## TELEGRAM BOT with ChatGPT, Built with Go lang

### How to start ?

1. Generate your Telegram bot token at https://t.me/BotFather
2. Generate OpenAI API key at https://platform.openai.com/account/api-keys
3. Create `.env` file in the root directory
4. Add variable named `TELEGRAM_BOT_TOKEN` and `OPENAI_API_KEY` inside of `.env` file,
   and then put your Telegram bot token and OpenAI API key as its value
5. Run `go get` to install required modules
6. Make sure to assign execute permission to `start.sh` file
7. Run `start.sh`

### Wanna modify bot ?

If you want, you can rebuild the whole source code 😂