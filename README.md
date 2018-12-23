# sshd-config

A tool for manipulating an `sshd_config` file

## requirements

golang 1.11+

## building

```shell
make build
```

## usage

```
Usage: sshd-config <command> <filename> <key> [<value>]
       sshd-config -h | --help
       sshd-config --version

Options:
  -h --help            Show this screen.
  --version            Show version.

Commands:
   add        Add a value to a key
   get        Get a key's values
   set        Set a value on a key
   unset      Unset all instances of a key
```