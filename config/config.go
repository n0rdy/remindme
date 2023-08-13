package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/logger"
	"n0rdy.me/remindme/utils"
	"os"
	"strconv"
)

func FetchAdminConfigs() (*common.AdminConfigs, error) {
	configsAsBytes, err := os.ReadFile(getAdminConfigsFilePath())
	if err != nil {
		return nil, err
	}

	adminConfigs := &common.AdminConfigs{}
	err = yaml.Unmarshal(configsAsBytes, adminConfigs)
	if err != nil {
		return nil, err
	}
	return adminConfigs, nil
}

func PersistAdminConfigs(adminConfigs common.AdminConfigs) error {
	configsAsBytes, err := yaml.Marshal(adminConfigs)
	if err != nil {
		return err
	}

	configsFile, err := os.Create(getAdminConfigsFilePath())
	if err != nil {
		return err
	}
	defer configsFile.Close()

	_, err = configsFile.Write(configsAsBytes)
	return err
}

func DeleteAdminConfigs() error {
	err := os.Remove(getAdminConfigsFilePath())
	if errors.Is(err, os.ErrNotExist) {
		logger.Info("admin configs file does not exist, nothing to delete")
		return nil
	}
	return err
}

func ResolveRunningServerPort() (int, error) {
	configs, err := FetchAdminConfigs()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// for backward compatibility with the previous versions when it was not possible to change the server port, so nothing was persisted
			// This should catch the case when the server has been started by the old version of the remindme app, and the user tries to stop it with the new version
			// Otherwise, the error will be thrown by the `stop` command, as it can't resolve the server port from the admin configs file
			//
			// Another reason: if the default port is requested for the `start` command,
			// it won't be persisted in the admin configs file to avoid OS-specific file access errors for the users that don't need the port binding feature (so, for the majority, I guess)
			logger.Info("server port is not persisted, falling back to default: " + strconv.Itoa(common.DefaultHttpServerPort))
			return common.DefaultHttpServerPort, nil
		}
		return 0, err
	}

	return configs.ServerPort, nil
}

func getAdminConfigsFilePath() string {
	return utils.GetOsSpecificAppDataDir() + common.AdminConfigsFileName
}
