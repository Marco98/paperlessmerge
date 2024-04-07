# PaperlessMerge

> [!NOTE]  
> Archived - obsolete with Paperless-ngx 2.7

PaperlessMerge can easily merge documents in paperless-ng(x)

## Usage

```
Usage:
  paperlessmerge docid1 docid2 ... [flags]

Flags:
  -c, --config string     config file
  -d, --delete            deletes merged documents
  -h, --help              help for paperlessmerge
  -k, --ignoretls         do not validate tls certificate
  -p, --password string   paperless auth password
  -s, --server string     paperless base url
  -u, --username string   paperless auth username
```

### Examples

```shell
paperlessmerge 100 101 102 # merge those 3 docs and upload the new one
paperlessmerge 160 161 -d # same as above but the originals are deleted
```

## Configuration

### Location

Paperlessmerge will try to load a config-file located in `$XDG_CONFIG_HOME/paperlessmerge/config.toml`.
For example on a common linux installation the file should be located in `~/.config/paperlessmerge/config.toml`
This path is resolved by [github.com/adrg/xdg](https://github.com/adrg/xdg).
For other OSes take a look at [Default locations](https://github.com/adrg/xdg/blob/master/README.md#xdg-base-directory=)

### Sample Configuration

```toml
server = "http://127.0.0.1:8080"
username = "admin"
password = "insecure"
```
