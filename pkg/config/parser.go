package config

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type Config struct {
	Mapper     Mapper         `toml:"mapper"`
	Resolution Resolution     `toml:"resolution"`
	Jobs       map[string]Job `toml:"job"`
	RTtl       int            `toml:"resolution-ttl"`
}

type Mapper struct {
	Bin      string `toml:"bin"`
	FilePath string `toml:"filepath"`
}

type Resolution struct {
	Hostname string   `toml:"hostname"`
	Dns      []string `toml:"dns"`
}

type Job struct {
	FromPort uint16 `toml:"from-port"`
	ToIp     string `toml:"to-ip"`
	ToPort   uint16 `toml:"to-port"`
	CheckTtl int64  `toml:"check-ttl"`
}

type IParser interface {
	Parse() (*Config, error)
}

type Parser struct {
	Path string
}

func (t *Parser) Parse() (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(t.Path, &config); err != nil {
		return nil, errors.Wrapf(err, "toml decode failed.")
	}

	return &config, nil
}

type IParserFactory interface {
	NewParser(path string) IParser
}

type ParserFactory struct {
}

func (pf *ParserFactory) NewParser(path string) IParser {
	return &Parser{Path: path}
}
