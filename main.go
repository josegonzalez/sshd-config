package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	docopt "github.com/docopt/docopt-go"
)

type entry struct {
	Key    string
	Values []string
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

func commandGet(arguments docopt.Opts, entries map[string]entry) {
	key, _ := arguments.String("<key>")
	if entry, ok := entries[key]; ok {
		for _, value := range entry.Values {
			println(value)
		}
	}
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
   unset      Unset all instances of a key`

	arguments, _ := docopt.ParseArgs(usage, os.Args[1:], "0.1.0")

	filename, err := arguments.String("<filename>")
	if err != nil {
		log.Printf("error: %s", err)
	}

	entries, err := configRead(filename)
	if err != nil {
		log.Printf("error: %s", err)
	}

	command, _ := arguments.String("<command>")
	switch command {
	case "add":
		commandAdd(arguments, entries, filename)
		break
	case "get":
		commandGet(arguments, entries)
		break
	case "set":
		commandSet(arguments, entries, filename)
		break
	case "unset":
		commandUnset(arguments, entries, filename)
		break
	default:
		log.Printf("error: %s", fmt.Errorf("%s is not a command. See 'ini-writer help'", command))
	}
}
