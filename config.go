package main

import "strings"

type DatabaseLine struct {
	Alias    string
	DBType   string
	Hostname string
	Port     string
	Username string
	Password string
}

type Config struct {
	NATSURL            string `envconfig:"NATS_URL" required:"true"`
	NATSToken          string `envconfig:"NATS_TOKEN" required:"false"`
	Databases          string `envconfig:"DATABASES" required:"true"` // alias:dbtype:hostname:port:username:password separated by semicolon
	NATSMetricsSubject string `envconfig:"NATS_METRICS_SUBJECT" required:"true" default:"svc.metrics"`
	MetricsIdent       string `envconfig:"METRICS_IDENT" required:"true" default:"storage_service"`
}

func (c *Config) DatabasesMap() map[string]DatabaseLine {
	databaseMap := map[string]DatabaseLine{}

	for _, line := range strings.Split(c.Databases, ";") {
		parts := strings.Split(line, ":")
		databaseLine := DatabaseLine{
			Alias:    parts[0],
			DBType:   parts[1],
			Hostname: parts[2],
			Port:     parts[3],
			Username: parts[4],
			Password: parts[5],
		}
		databaseMap[databaseLine.Alias+":"+databaseLine.DBType] = databaseLine
	}

	return databaseMap
}
