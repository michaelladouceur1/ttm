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
	ListFlags         ConfigDefaultFlags `yaml:"listFlags"`
}

type ConfigDefaultFlags struct {
	Category string `yaml:"category"`
	Priority string `yaml:"priority"`
	Status   string `yaml:"status"`
}

func NewConfig() (*gonfig.Gonfig[Config], error) {
	cfg := &Config{
		DaysToDisplay:     7,
		MaxTasksToDisplay: 25,
		AddFlags: ConfigDefaultFlags{
			Category: string(models.CategoryTask),
			Priority: string(models.PriorityHigh),
			Status:   string(models.StatusOpen),
		},
		ListFlags: ConfigDefaultFlags{
			Category: "",
			Priority: "",
			Status:   string(models.StatusOpen),
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
