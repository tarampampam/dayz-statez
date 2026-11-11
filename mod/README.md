## Installing and configuring

The mod is server-side only. **Clients need nothing installed**.

Download [`mod.tar.gz`](https://github.com/tarampampam/dayz-statez/releases/latest/download/mod.tar.gz)
from the latest release and unpack it into the DayZ server directory as the operating-system **user that runs DayZ**.
Replace `dayz` below with that user. Switch to that user, then run the following commands:

```shell
root$ su - dayz
dayz$ cd /path/to/dayz-server
dayz$ curl -fLO https://github.com/tarampampam/dayz-statez/releases/latest/download/mod.tar.gz
dayz$ tar -xzf ./mod.tar.gz
dayz$ rm ./mod.tar.gz
```

The archive creates `/path/to/dayz-server/@statez`. Start the server with `-serverMod=@statez` (under the same user):

```shell
dayz$ ./DayZServer -config=serverDZ.cfg -serverMod=@statez
```

Notes:

* Use exactly `-serverMod=@statez`, not `-mod=@statez`
* The PBO is unsigned, so set `verifySignatures=0;` in your server config

Config is created on first run at `$profile:statez/config.json`:

| Key                             | Default | Meaning                                                                     |
|---------------------------------|---------|-----------------------------------------------------------------------------|
| `interval`                      | `5`     | Seconds between collection cycles (clamped to 0-900)                        |
| `delete_on_shutdown`            | `true`  | Delete the JSON files on shutdown. Set `false` to inspect them after a stop |
| `thresholds.hit_damage`         | `3`     | Minimum damage to record a hit. `-1` disables the threshold                 |
| `thresholds.hit_damage_vehicle` | `15`    | Minimum vehicle damage to record a hit. `-1` disables the threshold         |

Snapshots go to `$profile:statez/data/`. To store them elsewhere, replace that directory with a symlink.

> [!NOTE]
> Positions use map axes - `pos_x` is engine X, `pos_y` is engine Z, `pos_alt` is engine Y, and `pos_angle` is yaw in
degrees.

### Uninstalling

Stop the server, remove `-serverMod=@statez` from your launch command, then delete the mod:

```shell
rm -rv /path/to/dayz-server/@statez
```

Remove the config and snapshots. `$profile` is your server profiles directory (`-profiles`):

```shell
rm -rv /path/to/profiles/statez
```

The mod deletes `data/*.json` on graceful shutdown when `delete_on_shutdown` is `true`. It never deletes `config.json`,
so step 2 is the only way to remove that file.

## Licence

GPL-3.0-or-later (see [`LICENSE`](./LICENSE)). It is based on WoozyMasta/metricz, Copyright (c) WoozyMasta.
