package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	docopt "github.com/docopt/docopt-go"
)

type entry struct {
	Key    string
	Values []string
}

var (
	// Version for sshd-config
	Version string
	elogger *log.Logger
	logger  *log.Logger
)

func init() {
	elogger = log.New(os.Stderr, "", 0)
	logger = log.New(os.Stdout, "", 0)
}

func readFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	lines := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return strings.Join(lines, "\n"), err
}

func configRead(path string) (map[string]entry, error) {
	re, _ := regexp.Compile(`^([A-Za-z0-9_]+) \s*(.+)$`)

	entries := make(map[string]entry)
	text, err := readFile(path)
	if err != nil {
		return entries, err
	}

	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		line = strings.ReplaceAll(line, "\t", " ")

		params := re.FindStringSubmatch(line)
		if len(params) != 3 {
			continue
		}

		key, value := params[1], params[2]
		values := []string{value}
		if entry, ok := entries[key]; ok {
			values = append(entry.Values, values...)
		}

		entries[key] = entry{Key: key, Values: values}
	}

	if scanner.Err() != nil {
		return entries, scanner.Err()
	}

	if len(entries) == 0 {
		return entries, fmt.Errorf("no entries found in file")
	}

	return entries, nil
}

func configWrite(entries map[string]entry, filename string) error {
	lines := []string{}
	for _, entry := range entries {
		for _, value := range entry.Values {
			lines = append(lines, fmt.Sprintf("%s %s", entry.Key, value))
		}
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	sort.Strings(lines)
	_, err = file.WriteString(strings.Join(lines, "\n"))
	if err != nil {
		return err
	}

	return nil
}

func commandAdd(arguments docopt.Opts, entries map[string]entry, filename string) {
	key, _ := arguments.String("<key>")
	value, _ := arguments.String("<value>")

	values := []string{value}
	if entry, ok := entries[key]; ok {
		values = append(entry.Values, []string{value}...)
	}

	entries[key] = entry{Key: key, Values: values}
	configWrite(entries, filename)
}

func commandGet(arguments docopt.Opts, entries map[string]entry) {
	key, _ := arguments.String("<key>")
	if entry, ok := entries[key]; ok {
		for _, value := range entry.Values {
			logger.Printf("%s", value)
		}
	}
}

func inList(name string, values []string, slice []string) bool {
	exit := true
	for _, value := range values {
		found := false
		for _, s := range slice {
			if s == value {
				found = true
				break
			}
		}

		if !found {
			elogger.Printf("error: for key '%s', expected one of '%s', actual '%s'", name, slice, value)
			exit = false
		}
	}

	return exit
}

func commandLint(arguments docopt.Opts, entries map[string]entry) {
	bestPracticesValidation := map[string]string{
		"AuthenticationMethods":   "publickey",
		"HostbasedAuthentication": "no",
		"IgnoreRhosts":            "yes",
		"PasswordAuthentication":  "no",
		"PermitEmptyPasswords":    "no",
		"PermitRootLogin":         "no",
		"Protocol":                "2",
		"PubkeyAuthentication":    "yes",
		"StrictModes":             "yes",
		"UsePrivilegeSeparation":  "yes",
	}

	bestPracticesSliceValidation := map[string][]string{
		"HostKey":       []string{"/etc/ssh/ssh_host_ed25519_key", "/etc/ssh/ssh_host_rsa_key"},
		"KexAlgorithms": []string{"curve25519-sha256@libssh.org", "diffie-hellman-group-exchange-sha256"},
		"Ciphers":       []string{"chacha20-poly1305@openssh.com", "aes256-gcm@openssh.com", "aes128-gcm@openssh.com", "aes256-ctr", "aes192-ctr", "aes128-ctr"},
		"MACs":          []string{"hmac-sha2-512-etm@openssh.com", "hmac-sha2-256-etm@openssh.com", "umac-128-etm@openssh.com", "hmac-sha2-512", "hmac-sha2-256", "umac-128@openssh.com"},
	}

	commaValidation := map[string][]string{
		"Ciphers": []string{"3des-cbc", "aes128-cbc", "aes192-cbc", "aes256-cbc", "aes128-ctr", "aes192-ctr", "aes256-ctr", "arcfour128", "arcfour256", "arcfour", "blowfish-cbc", "cast128-cbc"},
	}

	emptyValidation := map[string]bool{
		"AuthorizedKeysCommand":      true,
		"AuthorizedKeysCommandRunAs": true,
	}

	integerValidation := map[string]bool{
		"ClientAliveCountMax":     true,
		"ClientAliveInterval":     true,
		"KeyRegenerationInterval": true,
		"LoginGraceTime":          true,
		"MaxAuthTries":            true,
		"MaxSessions":             true,
		"MaxStartups":             true,
		"Port":                    true,
		"ServerKeyBits":           true,
		"X11DisplayOffset":        true,
	}

	listValidation := map[string][]string{
		"AddressFamily":   []string{"any", "inet", "inet6"},
		"Compression":     []string{"yes", "no", "delayed"},
		"GatewayPorts":    []string{"yes", "no", "clientspecified"},
		"LogLevel":        []string{"QUIET", "FATAL", "ERROR", "INFO", "VERBOSE", "DEBUG", "DEBUG1", "DEBUG2", "DEBUG3"},
		"PermitRootLogin": []string{"yes", "no", "forced-commands-only", "without-password"},
		"PermitTunnel":    []string{"yes", "no", "ethernet", "point-to-point"},
		"Protocol":        []string{"1", "2", "1,2", "2,1"},
		"SyslogFacility":  []string{"DAEMON", "USER", "AUTH", "AUTHPRIV", "LOCAL0", "LOCAL1", "LOCAL2", "LOCAL3", "LOCAL4", "LOCAL5", "LOCAL6", "LOCAL7"},
	}

	multipleValuesValidation := map[string]bool{
		"AcceptEnv":     true,
		"HostKey":       true,
		"ListenAddress": true,
		"Port":          true,
	}

	stringBoolValidation := map[string]bool{
		"AllowAgentForwarding":            true,
		"AllowTcpForwarding":              true,
		"ChallengeResponseAuthentication": true,
		"GSSAPICleanupCredentials":        true,
		"GSSAPIKeyExchange":               true,
		"GSSAPIStrictAcceptorCheck":       true,
		"HostbasedAuthentication":         true,
		"HostbasedUsesNameFromPacketOnly": true,
		"IgnoreRhosts":                    true,
		"IgnoreUserKnownHosts":            true,
		"KerberosAuthentication":          true,
		"KerberosGetAFSToken":             true,
		"KerberosOrLocalPasswd":           true,
		"KerberosTicketCleanup":           true,
		"KerberosUseKuserok":              true,
		"PasswordAuthentication":          true,
		"PermitEmptyPasswords":            true,
		"PermitUserEnvironment":           true,
		"PrintLastLog":                    true,
		"PrintMotd":                       true,
		"PubkeyAuthentication":            true,
		"RhostsRSAAuthentication":         true,
		"RSAAuthentication":               true,
		"ShowPatchLevel":                  true,
		"StrictModes":                     true,
		"TCPKeepAlive":                    true,
		"UseDNS":                          true,
		"UseLogin":                        true,
		"UsePAM":                          true,
		"UsePrivilegeSeparation":          true,
		"X11Forwarding":                   true,
		"X11UseLocalhost":                 true,
	}

	exitCode := 0
	for _, e := range entries {
		if validValue, ok := bestPracticesValidation[e.Key]; ok {
			if e.Values[0] != validValue {
				elogger.Printf("error: for key '%s', expected %s, actual '%s'", e.Key, validValue, e.Values[0])
				exitCode = 1
			}
		}

		if slice, ok := bestPracticesSliceValidation[e.Key]; ok {
			for _, value := range e.Values {
				values := strings.Split(value, ",")
				if !inList(e.Key, values, slice) {
					exitCode = 1
				}
			}
		}

		if slice, ok := commaValidation[e.Key]; ok {
			for _, value := range e.Values {
				values := strings.Split(value, ",")
				if !inList(e.Key, values, slice) {
					exitCode = 1
				}
			}
		}

		if _, ok := emptyValidation[e.Key]; !ok {
			for _, value := range e.Values {
				if len(value) == 0 {
					elogger.Printf("error: for key '%s', value may not be empty", e.Key)
					exitCode = 1
				}
			}
		}

		if _, ok := integerValidation[e.Key]; ok {
			for _, value := range e.Values {
				if _, err := strconv.Atoi(value); err != nil {
					elogger.Printf("error: for key '%s', expected integer, actual '%s'", e.Key, value)
					exitCode = 1
				}
			}
		}

		if slice, ok := listValidation[e.Key]; ok {
			if !inList(e.Key, e.Values, slice) {
				exitCode = 1
			}
		}

		if _, ok := multipleValuesValidation[e.Key]; !ok {
			if len(e.Values) > 1 {
				elogger.Printf("error: for key '%s', multiple values not allowed", e.Key)
				exitCode = 1
			}
		}

		if _, ok := stringBoolValidation[e.Key]; ok {
			if !inList(e.Key, e.Values, []string{"yes", "no"}) {
				exitCode = 1
			}
		}
	}

	os.Exit(exitCode)
}

func commandSet(arguments docopt.Opts, entries map[string]entry, filename string) {
	key, _ := arguments.String("<key>")
	value, _ := arguments.String("<value>")
	entries[key] = entry{Key: key, Values: []string{value}}
	configWrite(entries, filename)
}

func commandUnset(arguments docopt.Opts, entries map[string]entry, filename string) {
	key, _ := arguments.String("<key>")
	delete(entries, key)
	configWrite(entries, filename)
}

func main() {
	usage := `sshd-config.

Usage: sshd-config <command> [<key>] [<value>] [--filename=<filename>]
       sshd-config -h | --help
       sshd-config -v | --version

Options:
  -h --help              Show this screen.
  -v --version           Show version.
  --filename=<filename>  The sshd-config to modify [default: /etc/ssh/sshd_config]

Commands:
   add        Add a value to a key
   get        Get a key's values
   help       Print this help output
   lint       Lint a config against best practices
   set        Set a value on a key
   unset      Unset all instances of a key`

	arguments, _ := docopt.ParseArgs(usage, os.Args[1:], Version)

	filename, err := arguments.String("--filename")
	if err != nil {
		elogger.Printf("error: %s", err)
	}

	entries, err := configRead(filename)
	if err != nil {
		elogger.Printf("error: %s", err)
	}

	command, _ := arguments.String("<command>")
	switch command {
	case "add":
		commandAdd(arguments, entries, filename)
		break
	case "get":
		commandGet(arguments, entries)
		break
	case "help":
		fmt.Println(usage)
		break
	case "lint":
		commandLint(arguments, entries)
		break
	case "set":
		commandSet(arguments, entries, filename)
		break
	case "unset":
		commandUnset(arguments, entries, filename)
		break
	default:
		elogger.Printf("error: %s", fmt.Errorf("%s is not a command. See 'sshd-config help'", command))
	}
}
