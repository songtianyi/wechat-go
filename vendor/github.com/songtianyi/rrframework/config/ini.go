package rrconfig

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	sectionRegex = regexp.MustCompile(`^\[(.*)\]$`)
	assignRegex  = regexp.MustCompile(`^([^=]+)=(.*)$`)
)

type IniConfig struct {
	secs map[string]string
}

func LoadIniConfigFromFile(path string) (*IniConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(f)
	s := &IniConfig{
		secs: make(map[string]string),
	}
	if err := s.doParse(buf); err != nil {
		return nil, err
	}
	return s, nil
}

// Get("a.b")
// Get("a")
func (s *IniConfig) Get(key string) (string, error) {
	if v, ok := s.secs[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("key %s not exist", key)
}

func (s *IniConfig) Dump() string {
	rs := ""
	for k, v := range s.secs {
		rs += k + "=" + v + "\n"
	}
	return rs
}

func (s *IniConfig) doParse(buf *bufio.Reader) error {
	section := ""
	for {
		var (
			line string
			err  error
		)

		if line, err = buf.ReadString('\n'); err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		i := strings.IndexAny(line, ";")
		if i != -1 {
			// remove comment
			line = string(([]byte(line))[:i])
		}
		if len(line) == 0 {
			// Skip blank lines
			continue
		}
		line = strings.TrimSpace(line)

		if groups := assignRegex.FindStringSubmatch(line); groups != nil {
			key, val := groups[1], groups[2]
			key, val = strings.TrimSpace(key), strings.TrimSpace(val)
			s.secs[section+key] = val
		} else if groups := sectionRegex.FindStringSubmatch(line); groups != nil {
			name := strings.TrimSpace(groups[1])
			section = name + "."
		} else {
			return fmt.Errorf("parse error, invalid line %s", line)
		}

	}
	return nil
}
