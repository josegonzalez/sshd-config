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
Usage: sshd-config <command> [<key>] [<value>] [--filename=<filename>]
       sshd-config -h | --help
       sshd-config --version

Options:
  -h --help              Show this screen.
  --version              Show version.
  --filename=<filename>  The sshd-config to modify [default: /etc/ssh/sshd_config]

Commands:
   add        Add a value to a key
   get        Get a key's values
   lint       Lint a config against best practices
   set        Set a value on a key
   unset      Unset all instances of a key
```

### linting

Linting runs against the following rules:

- Multiple values only allowed for:
    - AcceptEnv
    - HostKey
    - ListenAddress
    - Port
- The following keys may *only* have the corresponding values:
		- AuthenticationMethods:   publickey
		- HostbasedAuthentication: no
		- IgnoreRhosts:            yes
		- PasswordAuthentication:  no
		- PermitEmptyPasswords:    no
		- PermitRootLogin:         no
		- Protocol:                2
		- PubkeyAuthentication:    yes
		- StrictModes:             yes
		- UsePrivilegeSeparation:  yes

Any violation of the above rules will result in the error being printed to stderr and non-zero exit code.
