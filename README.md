# snippetbox

A small Go web application (snippets + user accounts) backed by MySQL.

## Running locally

```
go run ./cmd/web
```

By default the app listens on `:4000` over HTTPS using the self-signed
certificate in `assets/tls/`, and connects to a local MySQL database with the
DSN `web:acosta@/snippetbox?parseTime=true`.

Flags:

| Flag | Default | Description |
| --- | --- | --- |
| `-addr` | `:4000`, or `:$PORT` when `PORT` is set | HTTP network address |
| `-dsn` | local DSN, or `$DSN` when set | MySQL data source name |
| `-tls` | `true`, `false` when `PLATFORM_APPLICATION` is set | Serve HTTPS with the local certificate |

## Deploy to Upsun

The repository ships with the configuration needed to deploy on
[Upsun](https://upsun.com):

- `.upsun/config.yaml` — one `golang` application container (`snippetbox`) and a
  MariaDB service, plus the routes.
- `.environment` — builds the `DSN` environment variable from the `MYSQL_*`
  variables Upsun injects for the `mysql` relationship.
- `db/schema.sql` — idempotent schema (`snippets`, `users` and the `sessions`
  table required by `scs/mysqlstore`), applied by the deploy hook on every
  deploy.

Steps:

```
upsun project:create
upsun push
```

The build hook compiles the app to `bin/app`; the web command starts it with
`-addr=":$PORT"`. Upsun's router terminates TLS, so the app itself serves plain
HTTP — `-tls` defaults to `false` on the platform because
`PLATFORM_APPLICATION` is set there.

### Troubleshooting

```
upsun log app        # application logs
upsun activity:log   # build and deploy hook output
upsun ssh            # shell into the container
```

If the deploy fails at boot, the most common cause is the database being
unreachable: the app pings MySQL at startup and exits non-zero when that fails.
Check that the `mysql` relationship is declared and that `DSN` is exported
(`upsun ssh -- env | grep -E 'DSN|MYSQL_'`).
