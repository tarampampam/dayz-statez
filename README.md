# StateZ

## Installing and starting `statez-server`

Place the binary at `/usr/local/bin/statez-server` first, then create a system user and the config file (run
everything as root):

```bash
useradd --system --no-create-home --shell /usr/sbin/nologin statez
install -d -m 0755 /etc/statez-server
echo 'STATEZ_SERVER_LISTEN=127.0.0.1:8080' > /etc/statez-server/.env
echo 'STATEZ_SERVER_MOD_DATA_DIR=/path/to/statez/data' >> /etc/statez-server/.env
echo 'STATEZ_SERVER_API_TOKEN=!!!CHANGE_ME!!!' >> /etc/statez-server/.env
chown root:statez /etc/statez-server/.env
chmod 0640 /etc/statez-server/.env
```

The `statez` user must be able to read the data directory. Make the `statez` group the owner of the data, give it read
access, make new files inherit the group (`g+s`), and grant traverse (`x`) on every parent directory (list all of them,
e.g. `/path` and `/path/to`):

```bash
chgrp -R statez /path/to/statez/data
chmod -R g+rX /path/to/statez/data
chmod g+s /path/to/statez/data
chmod o+x /path /path/to
```

New files keep the permissions of the process that writes them, so it must use a umask that allows group read
(e.g. `022`).

Check that `statez` can read the data:

```bash
runuser -u statez -- ls -l /path/to/statez/data
runuser -u statez -- cat /path/to/statez/data/world.json
```

Create the unit file `/etc/systemd/system/statez-server.service`:

```ini
[Unit]
Description=StateZ server
After=network-online.target
Wants=network-online.target

[Service]
User=statez
Group=statez
EnvironmentFile=/etc/statez-server/.env
ExecStart=/usr/local/bin/statez-server
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

Start the service:

```bash
systemctl daemon-reload
systemctl enable --now statez-server
systemctl status statez-server
journalctl -u statez-server -n 50 --no-pager
```

Check it (use the token from `.env`):

```bash
curl -i http://127.0.0.1:8080/healthz
curl -i -H "Authorization: Bearer !!!CHANGE_ME!!!" http://127.0.0.1:8080/world
```

### Logs

Logs go to journald:

```bash
journalctl -u statez-server -f
journalctl -u statez-server --since "1 hour ago" --no-pager
```

Log size is limited by journald settings (`SystemMaxUse=` in `/etc/systemd/journald.conf`). By default, it is 10%
of the filesystem, capped at 4G.

### Uninstalling

```bash
systemctl disable --now statez-server
rm /etc/systemd/system/statez-server.service /usr/local/bin/statez-server
rm -r /etc/statez-server
systemctl daemon-reload
userdel statez
```
