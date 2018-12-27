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
    - `AcceptEnv`
    - `HostKey`
    - `ListenAddress`
    - `Port`
- The following keys may *only* have the corresponding values:
    - `AuthenticationMethods`:   `publickey`
    - `HostbasedAuthentication`: `no`
    - `IgnoreRhosts`:            `yes`
    - `PasswordAuthentication`:  `no`
    - `PermitEmptyPasswords`:    `no`
    - `PermitRootLogin`:         `no`
    - `Protocol`:                `2`
    - `PubkeyAuthentication`:    `yes`
    - `StrictModes`:             `yes`
    - `UsePrivilegeSeparation`:  `yes`
- The following keys may have multiple values, but those values must *only* be within the corresponding list:
    - `HostKey`:                  `/etc/ssh/ssh_host_ed25519_key`, `/etc/ssh/ssh_host_rsa_key`
    - `KexAlgorithms`:            `curve25519-sha256@libssh.org`, `diffie-hellman-group-exchange-sha256`
    - `Ciphers`:                  `chacha20-poly1305@openssh.com`, `aes256-gcm@openssh.com`, `aes128-gcm@openssh.com`, `aes256-ctr`, `aes192-ctr`, `aes128-ctr`
    - `MACs`:                     `hmac-sha2-512-etm@openssh.com`, `hmac-sha2-256-etm@openssh.com`, `umac-128-etm@openssh.com`, `hmac-sha2-512`, `hmac-sha2-256`, `umac-128@openssh.com`
- The following keys are the *only* ones that may have an empty value:
    - `AuthorizedKeysCommand`
    - `AuthorizedKeysCommandRunAs`
- The following keys may *only* have one of the corresponding values:
    - `AddressFamily`:             `any`, `inet`, `inet6`
    - `Compression`:               `yes`, `no`, `delayed`
    - `GatewayPorts`:              `yes`, `no`, `clientspecified`
    - `LogLevel`:                  `QUIET`, `FATAL`, `ERROR`, `INFO`, `VERBOSE`, `DEBUG`, `DEBUG1`, `DEBUG2`, `DEBUG3`
    - `PermitRootLogin`:           `yes`, `no`, `forced-commands-only`, `without-password`
    - `PermitTunnel`:              `yes`, `no`, `ethernet`, `point-to-point`
    - `Protocol`:                  `1`, `2`, `1,2`, `2,1`
    - `SyslogFacility`:            `DAEMON`, `USER`, `AUTH`, `AUTHPRIV`, `LOCAL0`, `LOCAL1`, `LOCAL2`, `LOCAL3`, `LOCAL4`, `LOCAL5`, `LOCAL6`, `LOCAL7`
- The following keys may comprise of one or more - comma-separated - of the corresponding values (the error will point out an invalid value within the provided list):
    - `Ciphers`:                   `3des-cbc`, `aes128-cbc`, `aes192-cbc`, `aes256-cbc`, `aes128-ctr`, `aes192-ctr`, `aes256-ctr`, `arcfour128`, `arcfour256`, `arcfour`, `blowfish-cbc`, `cast128-cbc`
- The following keys may *only* have an integer as a value:
    - `ClientAliveCountMax`
    - `ClientAliveInterval`
    - `KeyRegenerationInterval`
    - `LoginGraceTime`
    - `MaxAuthTries`
    - `MaxSessions`
    - `MaxStartups`
    - `Port`
    - `ServerKeyBits`
    - `X11DisplayOffset`
- The following keys may *only* have a value of `yes` or `no`:
    - `AllowAgentForwarding`
    - `AllowTcpForwarding`
    - `ChallengeResponseAuthentication`
    - `GSSAPICleanupCredentials`
    - `GSSAPIKeyExchange`
    - `GSSAPIStrictAcceptorCheck`
    - `HostbasedAuthentication`
    - `HostbasedUsesNameFromPacketOnly`
    - `IgnoreRhosts`
    - `IgnoreUserKnownHosts`
    - `KerberosAuthentication`
    - `KerberosGetAFSToken`
    - `KerberosOrLocalPasswd`
    - `KerberosTicketCleanup`
    - `KerberosUseKuserok`
    - `PasswordAuthentication`
    - `PermitEmptyPasswords`
    - `PermitUserEnvironment`
    - `PrintLastLog`
    - `PrintMotd`
    - `PubkeyAuthentication`
    - `RhostsRSAAuthentication`
    - `RSAAuthentication`
    - `ShowPatchLevel`
    - `StrictModes`
    - `TCPKeepAlive`
    - `UseDNS`
    - `UseLogin`
    - `UsePAM`
    - `UsePrivilegeSeparation`
    - `X11Forwarding`
    - `X11UseLocalhost`

Any violation of the above rules will result in the error being printed to stderr and non-zero exit code.
