package exodus

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"exodus/internal/config"
	"exodus/internal/constant"
	"exodus/internal/db"
	"exodus/internal/dbutil"
	"exodus/internal/jobqueue"
	"exodus/internal/security"

	"github.com/redis/go-redis/v9"
)

type rescueResources struct {
	db    *sql.DB
	redis *redis.Client
}

type rescueSecretPayload struct {
	NodeCertPem  string `json:"nodeCertPem"`
	NodeKeyPem   string `json:"nodeKeyPem"`
	CaCertPem    string `json:"caCertPem"`
	JWTPublicKey string `json:"jwtPublicKey"`
}

func runRescueCLI() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Exodus Rescue CLI v0.4 ===")
	fmt.Println("Checking database connection...")
	resources, err := openRescueResources()
	if err != nil {
		return err
	}
	defer resources.close()
	fmt.Println("Database connected!")
	fmt.Println("Redis connected!")

	for {
		fmt.Println()
		fmt.Println("Select an action:")
		fmt.Println("1. Reset superadmin")
		fmt.Println("2. Enable password authentication")
		fmt.Println("3. Reset certs")
		fmt.Println("4. Get SECRET_KEY for an Exodus Node")
		fmt.Println("5. Fix Collation")
		fmt.Println("6. Clean up HWID Devices")
		fmt.Println("7. Clean up SRH Table")
		fmt.Println("8. Show version")
		fmt.Println("0. Exit")

		choice, err := promptLine(reader, "Select action", "0", false)
		if err != nil {
			return err
		}

		switch strings.ToLower(strings.TrimSpace(choice)) {
		case "1", "reset-superadmin":
			return resetSuperadmin(resources, reader)
		case "2", "enable-password-auth":
			return enablePasswordAuth(resources, reader)
		case "3", "reset-certs":
			return resetCerts(resources, reader)
		case "4", "get-secret-key-for-node":
			return getSecretKeyForNode(resources)
		case "5", "fix-postgres-collation":
			return fixPostgresCollation(resources, reader)
		case "6", "truncate-hwid-user-devices":
			return truncateHwidUserDevices(resources, reader)
		case "7", "truncate-srh-table":
			return truncateSRHTable(resources, reader)
		case "8", "version":
			fmt.Println(constant.GetBuildInfo())
		case "0", "q", "quit", "exit":
			fmt.Println("Exiting...")
			return nil
		default:
			fmt.Printf("Unknown action: %s\n", choice)
		}
	}
}

func openRescueResources() (*rescueResources, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	sqldb, err := db.OpenAndInitDB(&cfg)
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqldb.PingContext(ctx); err != nil {
		_ = sqldb.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	fmt.Println("Checking Redis connection...")
	redisClient, err := jobqueue.NewRedisClient(&cfg)
	if err != nil {
		_ = sqldb.Close()
		return nil, err
	}
	if redisClient == nil {
		_ = sqldb.Close()
		return nil, fmt.Errorf("redis is not configured")
	}

	return &rescueResources{
		db:    sqldb,
		redis: redisClient,
	}, nil
}

func (r *rescueResources) close() {
	if r == nil {
		return
	}
	if r.redis != nil {
		_ = r.redis.Close()
	}
	if r.db != nil {
		_ = r.db.Close()
	}
}

func resetSuperadmin(resources *rescueResources, reader *bufio.Reader) error {
	ok, err := promptConfirm(reader, "Are you sure you want to delete the superadmin?")
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("aborted")
	}

	var (
		adminUUID string
		username  string
	)
	err = resources.db.QueryRow(dbutil.Rebind(`
		SELECT uuid, username
		FROM admin
		ORDER BY created_at ASC
		LIMIT 1
	`)).Scan(&adminUUID, &username)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("superadmin not found")
	}
	if err != nil {
		return fmt.Errorf("find superadmin: %w", err)
	}

	if _, err := resources.db.Exec(dbutil.Rebind("DELETE FROM admin WHERE uuid = ?"), adminUUID); err != nil {
		return fmt.Errorf("delete superadmin: %w", err)
	}

	fmt.Printf("Superadmin %s deleted successfully.\n", username)
	return nil
}

func enablePasswordAuth(resources *rescueResources, reader *bufio.Reader) error {
	ok, err := promptConfirm(reader, "Are you sure you want to enable password authentication?")
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("aborted")
	}

	if _, err := resources.db.Exec(`
		UPDATE exodus_settings
		SET password_settings = jsonb_set(
			COALESCE(password_settings, '{}'::jsonb),
			'{enabled}',
			'true'::jsonb,
			true
		)
		WHERE id = 1
	`); err != nil {
		return fmt.Errorf("enable password authentication: %w", err)
	}

	fmt.Println("Password authentication enabled successfully.")
	return nil
}

func resetCerts(resources *rescueResources, reader *bufio.Reader) error {
	ok, err := promptConfirm(reader, "Are you sure you want to delete the certs? You will need to add new certs to all nodes again.")
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("aborted")
	}

	var keygenUUID string
	err = resources.db.QueryRow(`SELECT uuid FROM keygen ORDER BY created_at ASC LIMIT 1`).Scan(&keygenUUID)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("certs not found")
	}
	if err != nil {
		return fmt.Errorf("find certs: %w", err)
	}

	if _, err := resources.db.Exec(dbutil.Rebind("DELETE FROM keygen WHERE uuid = ?"), keygenUUID); err != nil {
		return fmt.Errorf("delete certs: %w", err)
	}

	fmt.Println("Certs deleted successfully.")
	fmt.Println(`Restart Exodus to apply changes by running "docker compose down && docker compose up -d".`)
	return nil
}

func getSecretKeyForNode(resources *rescueResources) error {
	var (
		pubKey string
		caCert string
		caKey  string
	)
	err := resources.db.QueryRow(`
		SELECT pub_key, ca_cert, ca_key
		FROM keygen
		ORDER BY created_at ASC
		LIMIT 1
	`).Scan(&pubKey, &caCert, &caKey)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("keygen not found; restart Exodus or reset certs first")
	}
	if err != nil {
		return fmt.Errorf("fetch keygen data: %w", err)
	}
	if strings.TrimSpace(caCert) == "" || strings.TrimSpace(caKey) == "" {
		return fmt.Errorf("certs not found; restart Exodus or reset certs first")
	}

	nodeCert, err := security.GenerateNodeCert(caCert, caKey)
	if err != nil {
		return fmt.Errorf("generate node certificate: %w", err)
	}

	raw, err := json.Marshal(rescueSecretPayload{
		NodeCertPem:  nodeCert.NodeCertPEM,
		NodeKeyPem:   nodeCert.NodeKeyPEM,
		CaCertPem:    caCert,
		JWTPublicKey: pubKey,
	})
	if err != nil {
		return fmt.Errorf("encode secret payload: %w", err)
	}

	fmt.Println("SECRET_KEY for node generated successfully.")
	fmt.Printf("\nSECRET_KEY=\"%s\"\n", base64.StdEncoding.EncodeToString(raw))
	return nil
}

func fixPostgresCollation(resources *rescueResources, reader *bufio.Reader) error {
	ok, err := promptConfirm(reader, "Are you sure you want to fix Collation?")
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("aborted")
	}

	var dbName string
	if err := resources.db.QueryRow(`SELECT current_database()`).Scan(&dbName); err != nil {
		return fmt.Errorf("read current database name: %w", err)
	}

	escapedName := strings.ReplaceAll(dbName, `"`, `""`)
	if _, err := resources.db.Exec(fmt.Sprintf(`ALTER DATABASE "%s" REFRESH COLLATION VERSION`, escapedName)); err != nil {
		return fmt.Errorf("refresh collation version: %w", err)
	}

	fmt.Println("Collation fixed successfully.")
	return nil
}

func truncateHwidUserDevices(resources *rescueResources, reader *bufio.Reader) error {
	ok, err := promptConfirm(reader, "Are you sure you want to clean up HWID Devices?")
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("aborted")
	}

	if _, err := resources.db.Exec(`TRUNCATE hwid_user_devices`); err != nil {
		return fmt.Errorf("clean up HWID Devices: %w", err)
	}

	fmt.Println("HWID Devices cleaned up successfully.")
	return nil
}

func truncateSRHTable(resources *rescueResources, reader *bufio.Reader) error {
	ok, err := promptConfirm(reader, "Are you sure you want to clean up SRH Table?")
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("aborted")
	}

	if _, err := resources.db.Exec(`TRUNCATE user_subscription_request_history RESTART IDENTITY`); err != nil {
		return fmt.Errorf("clean up SRH Table: %w", err)
	}

	fmt.Println("SRH Table cleaned up successfully.")
	return nil
}

func promptConfirm(reader *bufio.Reader, label string) (bool, error) {
	for {
		fmt.Printf("%s [y/N]: ", label)
		text, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("read confirmation: %w", err)
		}

		switch strings.ToLower(strings.TrimSpace(text)) {
		case "y", "yes", "true", "1":
			return true, nil
		case "", "n", "no", "false", "0":
			return false, nil
		default:
			fmt.Println("Please answer yes or no.")
		}
	}
}
