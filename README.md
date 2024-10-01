пример запроса 2FA аутентификации: curl -XPOST localhost:8000/auth -H 'Content-Type: application/json' -d '{"CODE":""}'

Пример ENV:

APP_ID=21837465
APP_HASH=9f5402c60e1ee6dbb9e9db64b6bfe621
PG_HOST="postgres://postgres:password@127.0.0.1:5432/tg"
CHAT_ID=123451254,1243246565436,2096696914,-1002065788839,-1002096696914,2065788839
SESSION_FILE="utils/ss.json"
PHONE_NUMBER="+79001411695"
VT_MOUNT_PATH="kv"
VT_READ_PATH="my-secret-password"
VT_WRITE_PATH="my-secret-password"
VT_HOST="http://127.0.0.1:8200"
VT_TOKEN="hvs.xOkDyN9Ygg4ApFc5wDRynbZs"
PORT="8231"
CODE="56144"