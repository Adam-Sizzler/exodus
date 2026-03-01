package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db"
	"v2ray-stat/backend/security"
)

func runAdminCredentialReset(cfg *config.BackendConfig) error {
	fmt.Println("=== Emergency admin credential reset ===")

	dbDir := filepath.Dir(cfg.Paths.Database)
	if err := os.MkdirAll(dbDir, 0o775); err != nil {
		return fmt.Errorf("create database directory: %w", err)
	}

	fileDB, err := db.OpenAndInitDB(cfg.Paths.Database, "file", cfg)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer fileDB.Close()

	reader := bufio.NewReader(os.Stdin)

	username, err := promptLine(reader, "Username", "admin", false)
	if err != nil {
		return err
	}

	password, err := promptLine(reader, "New password", "", false)
	if err != nil {
		return err
	}
	confirmPassword, err := promptLine(reader, "Confirm password", "", false)
	if err != nil {
		return err
	}
	if password != confirmPassword {
		return errors.New("password confirmation does not match")
	}

	passwordHash, err := security.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	var targetUUID string
	err = fileDB.QueryRow("SELECT uuid FROM admin WHERE username = ? LIMIT 1", username).Scan(&targetUUID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("lookup admin by username: %w", err)
	}

	if targetUUID == "" {
		return fmt.Errorf("admin user %q not found; create the first account from the web login page", username)
	}

	_, err = fileDB.Exec(`
		UPDATE admin
		SET password_hash = ?, role = 'ADMIN', updated_at = CURRENT_TIMESTAMP
		WHERE uuid = ?
	`, passwordHash, targetUUID)
	if err != nil {
		return fmt.Errorf("update existing admin credentials: %w", err)
	}

	if targetUUID != "" {
		if _, err := fileDB.Exec("DELETE FROM admin_sessions WHERE admin_uuid = ?", targetUUID); err != nil {
			return fmt.Errorf("revoke existing sessions: %w", err)
		}
	}

	fmt.Printf("Admin password successfully reset.\nUsername: %s\n", username)
	return nil
}

func promptLine(reader *bufio.Reader, label, defaultValue string, allowEmpty bool) (string, error) {
	if strings.TrimSpace(defaultValue) != "" {
		fmt.Printf("%s [%s]: ", label, defaultValue)
	} else {
		fmt.Printf("%s: ", label)
	}

	text, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read %s: %w", label, err)
	}
	text = strings.TrimSpace(text)

	if text == "" && defaultValue != "" {
		text = defaultValue
	}
	if !allowEmpty && strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("%s cannot be empty", strings.ToLower(label))
	}
	return text, nil
}
