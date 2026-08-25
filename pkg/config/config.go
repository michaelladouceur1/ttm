package config

import (
	_ "embed"
	"ttm/pkg/models"
	"ttm/pkg/paths"

	"github.com/michaelladouceur1/gonfig"
)

type Config struct {
	DaysToDisplay     int                `yaml:"daysToDisplay"`
	MaxTasksToDisplay int                `yaml:"maxTasksToDisplay"`
	AddFlags          ConfigDefaultFlags `yaml:"addFlags"`
	ListFlags         ConfigListFlags    `yaml:"listFlags"`
	Logging           LoggingConfig      `yaml:"logging"`
	Storage           StorageConfig      `yaml:"storage"`
}

type ConfigDefaultFlags struct {
	Priority string `yaml:"priority"`
	Status   string `yaml:"status"`
}

type ConfigListFlags struct {
	Priority []string `yaml:"priority"`
	Status   []string `yaml:"status"`
}

type StorageConfig struct {
	Type       string           `yaml:"type"`
	GoogleDocs GoogleDocsConfig `yaml:"googleDocs"`
	Postgres   PostgresConfig   `yaml:"postgres"`
}

type LoggingConfig struct {
	Theme string `yaml:"theme"`
}

type GoogleDocsConfig struct {
	DocumentID      string `yaml:"documentId"`
	CredentialsFile string `yaml:"credentialsFile"`
}

type PostgresConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	PasswordEnv string `yaml:"passwordEnv"`
	DBName      string `yaml:"dbname"`
	SSLMode     string `yaml:"sslmode"`
}

func NewConfig() (*gonfig.Gonfig[Config], error) {
	cfg := &Config{
		DaysToDisplay:     7,
		MaxTasksToDisplay: 25,
		Storage: StorageConfig{
			Type: "sqlite",
			Postgres: PostgresConfig{
				Host:    "localhost",
				Port:    5432,
				DBName:  "ttmdb",
				SSLMode: "disable",
			},
		},
		Logging: LoggingConfig{
			Theme: "minimal",
		},
		AddFlags: ConfigDefaultFlags{
			Priority: string(models.PriorityHigh),
			Status:   string(models.StatusOpen),
		},
		ListFlags: ConfigListFlags{
			Priority: []string{},
			Status:   []string{string(models.StatusOpen)},
		},
	}

	opts := gonfig.GonfigFileOptions{
		Type:           gonfig.YAML,
		RootDir:        paths.GetTTMDirectory(),
		Name:           "config",
		Watch:          true,
		ValidationMode: gonfig.VMRevert,
	}

	config, err := gonfig.NewGonfig(cfg, opts)
	if err != nil {
		return nil, err
	}

	return config, nil
}
