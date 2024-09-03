# moon-monorepo

This is a moon monorepo.

**Frontend Projects:**

- moon (Angular)

**Backend Projects:**

- user-vault
- gatekeeper
- logbot
- mini-fstore
- vfm
- acct

## Original Repositories

- [https://github.com/CurtisNewbie/moon](https://github.com/CurtisNewbie/moon)
- [https://github.com/CurtisNewbie/acct](https://github.com/CurtisNewbie/acct)
- [https://github.com/CurtisNewbie/event-pump](https://github.com/CurtisNewbie/event-pump)
- [https://github.com/CurtisNewbie/gatekeeper](https://github.com/CurtisNewbie/gatekeeper)
- [https://github.com/CurtisNewbie/logbot](https://github.com/CurtisNewbie/logbot)
- [https://github.com/CurtisNewbie/mini-fstore](https://github.com/CurtisNewbie/mini-fstore)
- [https://github.com/CurtisNewbie/user-vault](https://github.com/CurtisNewbie/user-vault)
- [https://github.com/CurtisNewbie/vfm](https://github.com/CurtisNewbie/vfm)

## Environment Preparation

### Prepare Frontend Environment

1. Install node_modules

```sh
cd moon/

npm ci
```

2. Start Angular Dev Server:

```sh
ng serve

# or if your openssl version isn't compatible
NODE_OPTIONS=--openssl-legacy-provider ng serve
```

### Prepare Backend Environment

1. Install related middlewares:

- RabbitMQ
- MySQL >= 5.7
- Redis
- Consul

```sh
# for example, using brew:
brew install rabbitmq
brew install mysql@5.7
brew install consul
brew install redis
```

2. Run middleware as services:

```sh
# for example, using brew; if you are using linux, you may run command like: `sudo systemctl rabbitmq-server start`
brew services start rabbitmq
brew services start mysql@5.7
brew services start consul
brew services start redis
```

3. Configure middlewares authentication/authorization stuff (e.g., for MySQL, RabbitMQ).
4. For each backend project, open file `conf.yml` and change the configuration for these middlewares:

```yaml
# configuration for consul ...
consul:
  enabled: "true"
  consulAddress: "localhost:8500"

# configuration for redis ...
redis:
  enabled: "true"
  address: "localhost"
  port: "6379"
  username: ""
  password: ""
  database: "0"

# configuration for mysql ...
mysql:
  enabled: "true"
  host: "localhost"
  port: "3306"
  user: "root"
  password: ""
  database: "acct"
  connection:
    parameters:
      - "charset=utf8mb4"
      - "parseTime=True"
      - "loc=Local"
      - "readTimeout=30s"
      - "writeTimeout=30s"
      - "timeout=3s"

# configuration for rabbitmq ...
rabbitmq:
  enabled: "true"
  host: "localhost"
  port: "5672"
  username: "guest"
  password: "guest"
  vhost: ""
```

5. Bootstrap each of the backend servers:

```sh
# e.g,.
go run main.go

# or
go run cmd/main.go
```

## Todo

- [ ] Provide Prometheus and Grafana Configuration Example.
- [ ] Provide Demo Snapshot.
- [ ] Provide Docker-Compose Example.