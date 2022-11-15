## SIGNALLER

### Environment vars required

```
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=signaller
DB_USERNAME=username
DB_PASSWORD=password
```

### Environment vars optional

```
DEBUG=false 

LOG_LEVEL=info

GRPC_PORT=8080

PROXY_GATEWAY=10.0.0.1:1080
PROXY_LIST_PATH=./proxy_list

TELEGRAM_TOKEN=token
TELEGRAM_CHAT_ID=id

METRIC_HTTP_PORT=8000
METRIC_PATH=/
```