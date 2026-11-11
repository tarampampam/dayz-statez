# Server configuration

Flag values have the highest priority, then environment variables, and finally defaults.

| Flag            | Env                          | Default          | Meaning                                                        |
|-----------------|------------------------------|------------------|----------------------------------------------------------------|
| `-listen`       | `STATEZ_SERVER_LISTEN`       | `0.0.0.0:8080`   | Address to listen on                                           |
| `-mod-data-dir` | `STATEZ_SERVER_MOD_DATA_DIR` | required         | The mod's `$profile:statez/data/` directory                    |
| `-api-token`    | `STATEZ_SERVER_API_TOKEN`    | empty (auth off) | Token for all endpoints except `/healthz`                      |
| `-tls-cert`     | `STATEZ_SERVER_TLS_CERT`     | empty            | Certificate file. Set together with `--tls-key` to serve HTTPS |
| `-tls-key`      | `STATEZ_SERVER_TLS_KEY`      | empty            | Private key file                                               |

## API

All endpoints are `GET` and return JSON. Each one returns the latest snapshot written by the mod. Send the token as
`Authorization: Bearer <token>` or `X-API-Key: <token>`.

`/healthz` and `/` endpoints do not need a token.

| Endpoint       | Returns                                                                                                                                        |
|----------------|------------------------------------------------------------------------------------------------------------------------------------------------|
| `/`            | JSON index of all endpoints                                                                                                                    |
| `/animals`     | Alive animals per type                                                                                                                         |
| `/areas`       | Effect areas (static gas zones and dynamic events): type, center, radius, players inside                                                       |
| `/events`      | Engine event counters since server start - connects, disconnects, deaths, chat, etc.                                                           |
| `/hits`        | Hit counts since server start, split by player/creature and ammo type                                                                          |
| `/players`     | Every online player - identity, character, position, vitals, state flags, network and lifetime stats                                           |
| `/territories` | Base territory flags - position, refresh-timer fraction, radius                                                                                |
| `/vehicles`    | Every tracked vehicle - class, position, health, crew, engine, lights, fuel, coolant, cargo count, per-part health                             |
| `/weapons`     | Shots per weapon type, live weapon entity counts, kills by source                                                                              |
| `/world`       | Map name, uptime, in-game clock, population counts, vehicle counts, loot and entity breakdowns, server FPS, weather, food/mining/fishing stats |
| `/zombies`     | Infected count per type and per AI state (calm, disturbed, alerted, chase, fight, dead)                                                        |
| `/healthz`     | Liveness probe, always `200`                                                                                                                   |

Response headers on data endpoints: `X-Cache-Hit` (served from cache), `Last-Modified` (mod file time), and
`X-Data-Age-Seconds` (how old the data is).

Example:

```shell
curl -s -H "Authorization: Bearer !!!CHANGE_ME!!!" http://127.0.0.1:12301/world
```

## Installing and starting

Run all the following as **root** to download the Linux `amd64` server archive, extract it, and install the binary:

```shell
curl -fLO https://github.com/tarampampam/dayz-statez/releases/latest/download/statez-server-linux-amd64.tar.gz
tar -xzf ./statez-server-linux-amd64.tar.gz
install -o root -g root -m 0755 ./statez-server /usr/local/bin/statez-server
rm ./statez-server-linux-amd64.tar.gz
```

> [!IMPORTANT]
> Replace "port", "mod data dir path" and "API token" with your own values!

```shell
install -d -m 0755 /etc/statez-server
echo 'STATEZ_SERVER_LISTEN=0.0.0.0:12301' > /etc/statez-server/.env # change port if needed
echo 'STATEZ_SERVER_MOD_DATA_DIR=/path/to/statez/data' >> /etc/statez-server/.env # with `.../data` at the end
echo 'STATEZ_SERVER_API_TOKEN=!!!CHANGE_ME!!!' >> /etc/statez-server/.env
chmod 0600 /etc/statez-server/.env
```

Since we're running API runs as the DayZ user, it has access to the snapshots without any additional ownership or
permission changes. Create the unit file `/etc/systemd/system/statez-server.service` (replace `dayz` with the user
that runs the DayZ server):

```ini
[Unit]
Description=StateZ server
After=network-online.target
Wants=network-online.target
#After=dayz.service # start after the DayZ process starts
#Wants=dayz.service # start DayZ together with StateZ; a DayZ failure does not stop StateZ
#PartOf=dayz.service # stop/restart StateZ when DayZ is stopped/restarted; one-way, does not start DayZ
#Requires=dayz.service # fails StateZ start if DayZ fails to start

[Service]
User=dayz
Group=dayz
EnvironmentFile=/etc/statez-server/.env
ExecStart=/usr/local/bin/statez-server
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

Start the service:

```shell
systemctl daemon-reload
systemctl enable --now statez-server
systemctl status statez-server
journalctl -u statez-server -n 50 --no-pager
```

Check it (use the token and port you configured in `.env`):

```shell
ss -ltnp | grep 12301 # use your port
curl -i http://127.0.0.1:12301/healthz
curl -i -H "Authorization: Bearer !!!CHANGE_ME!!!" http://127.0.0.1:12301/world
```

### Logs

Logs go to journald:

```shell
journalctl -u statez-server -f
journalctl -u statez-server --since "1 hour ago" --no-pager
```

> [!NOTE]
> Log size is limited by journald settings, not by StateZ. See `journald.conf(5)`, for example `SystemMaxUse=` in
> `/etc/systemd/journald.conf`.

### Uninstalling

```shell
systemctl disable --now statez-server
rm /etc/systemd/system/statez-server.service /usr/local/bin/statez-server
rm -r /etc/statez-server
systemctl daemon-reload
```
