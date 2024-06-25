package config

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"ttm/pkg/models"
	"ttm/pkg/paths"
)

//go:embed default/config.json
var defaultConfig string

type Config struct {
	AddFlags          ConfigDefaultFlags `json:"addFlags"`
	ListFlags         ConfigDefaultFlags `json:"listFlags"`
	OutPath           string             `json:"outPath"`
	DaysToDisplay     int                `json:"daysToDisplay"`
	MaxTasksToDisplay int                `json:"maxTasksToDisplay"`
}

type ConfigDefaultFlags struct {
	Category string `json:"category"`
	Priority string `json:"priority"`
	Status   string `json:"status"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Init() error {
	var err error

	if os.MkdirAll(paths.GetTaskStorePath(), os.ModePerm); err != nil {
		return err
	}

	if _, err = os.Stat(paths.GetConfigPath()); err != nil {
		if err := os.WriteFile(paths.GetConfigPath(), []byte(defaultConfig), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) Load() error {
	file, err := os.Open(paths.GetConfigPath())
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(c); err != nil {
		return err
	}

	return nil
}

func (c *Config) Update() error {
	configJson, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if err := os.WriteFile(paths.GetConfigPath(), []byte(configJson), 0644); err != nil {
		return err
	}

	return nil
}

func (c *Config) UpdateAddFlags(flag interface{}, val string) error {
	var err error
	switch flag {
	default:
		return errors.New("invalid flag")
	case "category":
		err = models.Category(val).Validate()
		if err != nil {
			fmt.Println("Error updating config: ", err)
			return err
		}
		c.AddFlags.Category = val
	case "priority":
		err = models.Priority(val).Validate()
		if err != nil {
			fmt.Println("Error updating config: ", err)
			return err
		}
		c.AddFlags.Priority = val
	case "status":
		err = models.Status(val).Validate()
		if err != nil {
			fmt.Println("Error updating config: ", err)
			return err
		}
		c.AddFlags.Status = val
	}

	return nil
}
