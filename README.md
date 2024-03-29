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

PROVIDER_POSITION_REFRESH_DURATION=10s
PROVIDER_POSITION_REFRESH_DURATION_FLOATING=false
PROVIDER_TRADERS_REFRESH_DURATION=24h

PROXY_GATEWAY=10.0.0.1:1080
PROXY_LIST_PATH=./proxy_list

TELEGRAM_TOKEN=token
TELEGRAM_CHAT_ID=id

METRIC_HTTP_PORT=8000
METRIC_PATH=/
```

## K8s

### Secrets

```
kubectl create secret generic signaller \
  --from-literal=db_database='signaller' \
  --from-literal=db_username='username' \
  --from-literal=db_password='password' \
  --from-literal=telegram_token='token' \
  --from-literal=telegram_chat_id='chat_id'
```