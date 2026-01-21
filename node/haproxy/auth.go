package haproxy

import (
	"encoding/csv"
	"fmt"
	"os"
	"v2ray-stat/node/config"
)

func UpdateUsersCSV(cfg *config.NodeConfig, usersToAdd map[string]string, usersToRemove []string) error {
	path := cfg.Paths.HAProxyAuth

	var records [][]string

	if _, err := os.Stat(path); err == nil {
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		reader := csv.NewReader(file)
		records, err = reader.ReadAll()
		file.Close()
		if err != nil {
			return err
		}
	}

	userMap := make(map[string][]string)
	for _, row := range records {
		if len(row) >= 3 {
			userMap[row[1]] = row // Индекс 1 - это username
		}
	}

	for _, user := range usersToRemove {
		delete(userMap, user)
	}

	for user, uuid := range usersToAdd {
		userMap[user] = []string{"1", user, uuid}
	}

	var newRecords [][]string
	for _, row := range userMap {
		newRecords = append(newRecords, row)
	}

	tmpPath := path + ".tmp"
	file, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	w := csv.NewWriter(file)
	if err := w.WriteAll(newRecords); err != nil {
		file.Close()
		return err
	}
	file.Close()

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to move temp csv file: %v", err)
	}

	return nil
}
