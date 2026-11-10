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

The `statez` user must be able to read the data directory. Grant read access to the existing files and to files
created later, and traverse access (`x`) on every parent directory (list all of them, e.g. `/path` and `/path/to`):

```bash
setfacl -R -m u:statez:rX /path/to/statez/data
setfacl -R -d -m u:statez:rX /path/to/statez/data
setfacl -m u:statez:x /path /path/to
```

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

### Uninstalling

```bash
systemctl disable --now statez-server
rm /etc/systemd/system/statez-server.service /usr/local/bin/statez-server
rm -r /etc/statez-server
systemctl daemon-reload
userdel statez
rm /usr/local/bin/statez-server
```
