# README

## Install
```bash
go install github.com/spf13/cobra-cli@latest
```

## Add Commands
```bash
# First, copy the .cobra.yaml to ~/.cobra.yaml

cobra-cli init
cobra-cli add <command>
# ex: cobra-init add hosts

# specify the parent of the new command
cobra-cli add list -p hostsCmd
```

## Viper
config.yaml
```yaml
hosts-file: newFile.hosts
```

```bash
PSCAN_HOSTS_FILE=newFile.hosts ./pscan hosts add host01 host02
PSCAN_HOSTS_FILE=newFile.hosts ./pscan hosts list

./pscan hosts list --config config.yaml
```