# moon-monorepo

This is moon monorepo, a privately hosted personal website (the private cloud).

**Frontend Projects:**

- moon (Angular)

**Backend Projects:**

- user-vault (User Authentication and Authorization Service)
- gatekeeper (Application Gateway)
- logbot (Error Log Watching Service)
- mini-fstore (Simple File Storage Service)
- vfm (Virtual File Management Service)
- acct (Simple Personal Accounting Service) \[**Developing, not deployed yet**\]

## Original Repositories

- [https://github.com/CurtisNewbie/moon](https://github.com/CurtisNewbie/moon)
- [https://github.com/CurtisNewbie/event-pump](https://github.com/CurtisNewbie/event-pump)
- [https://github.com/CurtisNewbie/gatekeeper](https://github.com/CurtisNewbie/gatekeeper)
- [https://github.com/CurtisNewbie/logbot](https://github.com/CurtisNewbie/logbot)
- [https://github.com/CurtisNewbie/mini-fstore](https://github.com/CurtisNewbie/mini-fstore)
- [https://github.com/CurtisNewbie/user-vault](https://github.com/CurtisNewbie/user-vault)
- [https://github.com/CurtisNewbie/vfm](https://github.com/CurtisNewbie/vfm)
- [https://github.com/CurtisNewbie/acct](https://github.com/CurtisNewbie/acct)

## Development Environment Preparation

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
4. For each backend project, open the file `conf.yml` and change the configuration:

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

5. For each backend service, go to `**/schema/` folder and execute the `schema.sql` DDL script.
6. Bootstrap each of the backend servers:

```sh
# e.g,.
go run main.go

# or
go run cmd/main.go
```

7. Finally, if everything goes right, you can visit the website via `https://localhost:4200`.

## Deployment

> [!NOTE]
>
> The whole setup is based on docker-compose and this project is privately hosted on my machine. Below includes examples of some of the configuration files.
>
> Instead of using `conf.yml`, you should be using `conf-prod.yml` instead, because most of the configurations are externalized.

The `conf-prod.yml` configuration file uses `${...}` syntax to read values from environment variables. E.g., the example below loads `rabbitmq.host`, `rabbitmq.username`, `rabbitmq.password` from environment variables `RABBITMQ_ADDR=`, `RABBITMQ_USERNAME=`, and `RABBITMQ_PASSWORD=`.

```yaml
rabbitmq:
  enabled: true
  host: "${RABBITMQ_ADDR}"
  port: 5672
  username: "${RABBITMQ_USERNAME}"
  password: "${RABBITMQ_PASSWORD}"
  vhost: "/"
```

- [Docker-Compose Conf](./deploy/docker-compose.yml)
- [Environment Variables And Secrets](./deploy/backend.env)
- [Nginx Conf](./deploy/nginx.conf)
- [Prometheus Conf](./deploy/prometheus.yml)

TODO ...