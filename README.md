# moon-monorepo

moon monorepo - a privately hosted website developed for myself for fun :D

**Frontend Projects:**

- moon (Angular)

**Backend Projects:**

- user-vault (User Authentication and Authorization Service)
- gatekeeper (Application Gateway)
- logbot (Error Log Watching Service)
- mini-fstore (Simple File Storage Service)
- vfm (Virtual File Management Service)
- acct (Simple Personal Accounting Service, `acct` only supports WeChat, it's completely optional)

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

### Script Example

For example, to bootstrap all backend and frontend services:

```bash
for r in $(ls "./moon-monorepo/backend");
do
    (
      cd "./moon-monorepo/backend/$r";
      if [ -f "main.go" ]; then
          go run main.go "logging.rolling.file=./logs/$r.log" 'logging.file.max-backups=1' 'logging.file.max-size=30' > /dev/null 2>&1 &
      else
          go run cmd/main.go "logging.rolling.file=./logs/$r.log" 'logging.file.max-backups=1' 'logging.file.max-size=30' > /dev/null 2>&1 &
      fi
    );
done;

( cd "./moon-monorepo/frontend/moon"; ng serve > /dev/null 2>&1 & )
```

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
- [Grafana Dashboard](./deploy/grafana_dashboard.json)

## Compatibility

Before moon-monorepo, backend projects are maintained in separate repos. Both v0.0.1 and v0.0.2 are compatible with following releases (in terms of database schema):

- [acct >= v0.0.0](https://github.com/CurtisNewbie/acct)
- [vfm >= v0.1.21](https://github.com/CurtisNewbie/vfm/tree/v0.1.21)
- [user-vault >= v0.0.26](https://github.com/CurtisNewbie/user-vault/tree/v0.0.26)
- [event-pump >= v0.0.14](https://github.com/CurtisNewbie/event-pump/tree/v0.0.14)
- [gatekeeper >= v0.0.23](https://github.com/CurtisNewbie/gatekeeper/tree/v0.0.23)
- [logbot >= v0.0.9](https://github.com/CurtisNewbie/logbot/tree/v0.0.9)
- [mini-fstore >= v0.1.21](https://github.com/CurtisNewbie/mini-fstore/tree/v0.1.21)

Meaning that you can directly upgrade to the code in moon-monorepo without worrying about the data migration if your deployment is already up-to-date. However, if you are not using the latests releases in these repos, you may need to consider upgrading to the latest versions before migrating to this monorepo version.

> [!IMPORTANT]
>
> Previously, some of the backend projects rely on `svc` to automatically upgrade schema, this functionality is now removed. DDL changes for each release is maintained in separate SQL files, you will have to execute the DDL scripts yourself based on the version you are using.

TODO ...
