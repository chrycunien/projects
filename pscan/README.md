# README

## Install
```bash
go install github.com/spf13/cobra-cli@latest
```

## Usage
```bash
# First, copy the .cobra.yaml to ~/.cobra.yaml

cobra-cli init
cobra-cli add <command>
# ex: cobra-init add hosts

# specify the parent of the new command
cobra-cli add list -p hostsCmd
```