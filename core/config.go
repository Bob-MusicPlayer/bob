package core

import (
	"encoding/json"
	"github.com/Bob-MusicPlayer/bob/model"
	"io/ioutil"
)

type ConfigManager struct {
	env    *Environment
	Config *model.Config
}

func NewConfigManager(env *Environment) *ConfigManager {
	return &ConfigManager{
		env: env,
	}
}

func (cm *ConfigManager) ReadConfig() error {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}

	var config model.Config

	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	cm.Config = &config

	return nil
}

func (cm *ConfigManager) GetPlayerBySource(source string) *model.Player {
	for _, player := range cm.env.ConfigManager.Config.Player {
		if player.Source == source {
			return player
		}
	}

	return nil
}
