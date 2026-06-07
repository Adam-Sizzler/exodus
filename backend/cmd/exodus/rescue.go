package exodus

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"exodus/internal/config"
	"exodus/internal/db"
	"exodus/internal/dbutil"
	"exodus/internal/jobqueue"
	"exodus/internal/security"

	"github.com/redis/go-redis/v9"
	"golang.org/x/term"
)

type cliFlags struct {
	rescue bool
}

type rescueResources struct {
	db    *sql.DB
	redis *redis.Client
	cfg   *config.BackendConfig
}

type rescueSecretPayload struct {
	NodeCertPem  string `json:"nodeCertPem"`
	NodeKeyPem   string `json:"nodeKeyPem"`
	CaCertPem    string `json:"caCertPem"`
	JWTPublicKey string `json:"jwtPublicKey"`
}

type cliAction struct {
	Value string
	Label string
	Hint  string
}

func (a cliAction) String() string {
	return a.Label
}

func parseCLIFlags() cliFlags {
	var flags cliFlags

	for _, arg := range os.Args[1:] {
		switch strings.ToLower(strings.TrimSpace(arg)) {
		case "cli", "rescue", "--rescue", "-rescue":
			flags.rescue = true
		}
	}

	return flags
}

func runPreConfigCLI(flags cliFlags) bool {
	if !flags.rescue {
		return false
	}

	if err := runRescueCLI(); err != nil {
		fmt.Printf("❌ Rescue CLI failed: %v\n", err)
		os.Exit(1)
	}

	return true
}

func runConfiguredCLI(_ cliFlags, _ *config.BackendConfig) bool {
	return false
}

func printRescueHint() {
	fmt.Println("Hint: run `docker exec -it exodus cli` for rescue CLI.")
}

func runRescueCLI() error {
	reader := bufio.NewReader(os.Stdin)

	printCLIBox("Exodus Rescue CLI v0.4")

	printStatus("◐", "🌱 Checking database connection...")
	resources, err := openRescueResources()
	if err != nil {
		printStatus("✖", "❌ Failed to connect to database.")
		return err
	}
	defer resources.close()
	printStatus("✔", "✅ Database connected!")

	printStatus("◐", "🌱 Checking Redis connection...")
	if err := checkRescueRedis(resources); err != nil {
		printStatus("✖", "❌ Failed to connect to Redis.")
		return err
	}
	printStatus("✔", "✅ Redis connected!")

	action, err := promptAction()
	if err != nil {
		return err
	}

	switch action {
	case "reset-superadmin":
		return resetSuperadmin(resources, reader)
	case "enable-password-auth":
		return enablePasswordAuth(resources, reader)
	case "reset-certs":
		return resetCerts(resources, reader)
	case "get-secret-key-for-node":
		return getSecretKeyForNode(resources)
	case "fix-postgres-collation":
		return fixPostgresCollation(resources, reader)
	case "truncate-hwid-user-devices":
		return truncateHwidUserDevices(resources, reader)
	case "truncate-srh-table":
		return truncateSRHTable(resources, reader)
	case "exit":
		printStatus("ℹ", "👋 Exiting...")
		return nil
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}

func openRescueResources() (*rescueResources, error) {
	var cfg config.BackendConfig
	var sqldb *sql.DB

	err := runSilently(func() error {
		loadedCfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("load configuration: %w", err)
		}

		cfg = loadedCfg

		sqldb, err = db.OpenAndInitDB(&cfg)
		if err != nil {
			return fmt.Errorf("connect database: %w", err)
		}

		return nil
	})
	if err != nil {
		if sqldb != nil {
			_ = sqldb.Close()
		}

		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqldb.PingContext(ctx); err != nil {
		_ = sqldb.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	var redisClient *redis.Client

	err = runSilently(func() error {
		client, err := jobqueue.NewRedisClient(&cfg)
		if err != nil {
			return err
		}

		if client == nil {
			return fmt.Errorf("redis is not configured")
		}

		redisClient = client

		return nil
	})
	if err != nil {
		_ = sqldb.Close()
		return nil, err
	}

	return &rescueResources{
		db:    sqldb,
		redis: redisClient,
		cfg:   &cfg,
	}, nil
}

func checkRescueRedis(resources *rescueResources) error {
	if resources == nil || resources.redis == nil {
		return fmt.Errorf("redis is not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := resources.redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
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

type actionKey int

const (
	actionKeyNone actionKey = iota
	actionKeyUp
	actionKeyDown
	actionKeyEnter
	actionKeyInterrupt
)

func promptAction() (string, error) {
	actions := rescueActions()
	if len(actions) == 0 {
		return "", errors.New("no rescue actions configured")
	}

	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		return promptActionPlain(actions)
	}

	selected := len(actions) - 1

	fmt.Println()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("enable raw terminal mode: %w", err)
	}

	terminalRestored := false
	restoreTerminal := func() {
		if terminalRestored {
			return
		}

		_ = term.Restore(int(os.Stdin.Fd()), oldState)
		terminalRestored = true
	}
	defer restoreTerminal()

	renderActionPrompt(actions, selected)

	for {
		key, err := readActionKey(os.Stdin)
		if err != nil {
			restoreTerminal()
			clearPromptArea(len(actions) + 2)
			return "", fmt.Errorf("read action key: %w", err)
		}

		switch key {
		case actionKeyUp:
			selected--
			if selected < 0 {
				selected = len(actions) - 1
			}
			clearPromptArea(len(actions) + 1)
			renderActionPrompt(actions, selected)
		case actionKeyDown:
			selected++
			if selected >= len(actions) {
				selected = 0
			}
			clearPromptArea(len(actions) + 1)
			renderActionPrompt(actions, selected)
		case actionKeyEnter:
			restoreTerminal()
			clearPromptArea(len(actions) + 2)
			printSelectedAction(actions[selected].Label)
			return actions[selected].Value, nil
		case actionKeyInterrupt:
			restoreTerminal()
			clearPromptArea(len(actions) + 2)
			printStatus("ℹ", "👋 Exiting...")
			os.Exit(0)
		}
	}
}

func rescueActions() []cliAction {
	return []cliAction{
		{
			Value: "reset-superadmin",
			Label: "Reset superadmin",
			Hint:  "Fully reset superadmin",
		},
		{
			Value: "enable-password-auth",
			Label: "Enable password authentication",
			Hint:  "Enable password authentication",
		},
		{
			Value: "reset-certs",
			Label: "Reset certs",
			Hint:  "Fully reset certs",
		},
		{
			Value: "get-secret-key-for-node",
			Label: "Get SECRET_KEY for an Exodus Node",
			Hint:  "Get SECRET_KEY in cases, where you can not get from Panel",
		},
		{
			Value: "fix-postgres-collation",
			Label: "Fix Collation",
			Hint:  "Fix Collation issues for current database",
		},
		{
			Value: "truncate-hwid-user-devices",
			Label: "Clean up HWID Devices",
			Hint:  "Remove all HWID Devices from the database",
		},
		{
			Value: "truncate-srh-table",
			Label: "Clean up SRH Table",
			Hint:  "Remove all SRH data from the database",
		},
		{
			Value: "exit",
			Label: "Exit",
		},
	}
}

func promptActionPlain(actions []cliAction) (string, error) {
	fmt.Println()
	fmt.Println("Select an action:")
	for index, action := range actions {
		if action.Hint != "" {
			fmt.Printf("%d) %s (%s)\n", index+1, action.Label, action.Hint)
			continue
		}

		fmt.Printf("%d) %s\n", index+1, action.Label)
	}

	fmt.Printf("Enter number [%d]: ", len(actions))

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read action: %w", err)
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return actions[len(actions)-1].Value, nil
	}

	selected, err := strconv.Atoi(text)
	if err != nil || selected < 1 || selected > len(actions) {
		return "", fmt.Errorf("invalid action: %s", text)
	}

	selected--
	printSelectedAction(actions[selected].Label)

	return actions[selected].Value, nil
}

func renderActionPrompt(actions []cliAction, selected int) {
	fmt.Printf("%s %s\r\n", ansi("36", "❯"), ansi("36", "Select an action"))

	for index, action := range actions {
		if index == selected {
			fmt.Printf("%s\r\n", formatActiveAction(action))
			continue
		}

		fmt.Printf("%s\r\n", formatInactiveAction(action))
	}
}

func formatActiveAction(action cliAction) string {
	line := fmt.Sprintf("%s %s", ansi("32", "●"), action.Label)
	if action.Hint != "" {
		line += fmt.Sprintf(" %s", ansi("90", fmt.Sprintf("(%s)", action.Hint)))
	}

	return line
}

func formatInactiveAction(action cliAction) string {
	return fmt.Sprintf("%s %s", ansi("90", "○"), ansi("90", action.Label))
}

func readActionKey(reader io.Reader) (actionKey, error) {
	var input [1]byte
	if _, err := reader.Read(input[:]); err != nil {
		return actionKeyNone, err
	}

	switch input[0] {
	case 3:
		return actionKeyInterrupt, nil
	case '\r', '\n':
		return actionKeyEnter, nil
	case 'k', 'K':
		return actionKeyUp, nil
	case 'j', 'J':
		return actionKeyDown, nil
	case 27:
		var sequence [2]byte
		if _, err := io.ReadFull(reader, sequence[:]); err != nil {
			return actionKeyNone, err
		}

		if sequence[0] != '[' {
			return actionKeyNone, nil
		}

		switch sequence[1] {
		case 'A':
			return actionKeyUp, nil
		case 'B':
			return actionKeyDown, nil
		}
	}

	return actionKeyNone, nil
}

func resetSuperadmin(resources *rescueResources, reader *bufio.Reader) error {
	answer, err := promptConfirm(reader, "Are you sure you want to delete the superadmin?")
	if err != nil {
		return err
	}
	if !answer {
		return errors.New("aborted")
	}

	printStatus("◐", "🔄 Deleting superadmin...")

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

	printStatus("✔", fmt.Sprintf("✅ Superadmin %s deleted successfully.", username))

	return nil
}

func enablePasswordAuth(resources *rescueResources, reader *bufio.Reader) error {
	printStatus("◐", "🔄 Enabling password authentication...")

	answer, err := promptConfirm(reader, "Are you sure you want to enable password authentication?")
	if err != nil {
		return err
	}
	if !answer {
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

	printStatus("✔", "✅ Password authentication enabled successfully.")

	return nil
}

func resetCerts(resources *rescueResources, reader *bufio.Reader) error {
	answer, err := promptConfirm(
		reader,
		"Are you sure you want to delete the certs? You will need to add new certs to all nodes again.",
	)
	if err != nil {
		return err
	}
	if !answer {
		return errors.New("aborted")
	}

	printStatus("◐", "🔄 Deleting certs...")

	var keygenUUID string

	err = resources.db.QueryRow(`
		SELECT uuid
		FROM keygen
		ORDER BY created_at ASC
		LIMIT 1
	`).Scan(&keygenUUID)

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("certs not found")
	}
	if err != nil {
		return fmt.Errorf("find certs: %w", err)
	}

	if _, err := resources.db.Exec(dbutil.Rebind("DELETE FROM keygen WHERE uuid = ?"), keygenUUID); err != nil {
		return fmt.Errorf("delete certs: %w", err)
	}

	printStatus("✔", "✅ Certs deleted successfully.")
	fmt.Println(`⚠ Restart Exodus to apply changes by running "docker compose down && docker compose up -d".`)

	return nil
}

func getSecretKeyForNode(resources *rescueResources) error {
	printStatus("◐", "🔑 Getting SECRET_KEY for node...")

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
		return fmt.Errorf("keygen not found; reset certs first or restart Exodus")
	}
	if err != nil {
		return fmt.Errorf("fetch keygen data: %w", err)
	}

	if strings.TrimSpace(caCert) == "" || strings.TrimSpace(caKey) == "" {
		return fmt.Errorf("certs not found; reset certs first or restart Exodus")
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

	printStatus("✔", "✅ SECRET_KEY for node generated successfully.")
	fmt.Printf("\nSECRET_KEY=\"%s\"\n", base64.StdEncoding.EncodeToString(raw))

	return nil
}

func fixPostgresCollation(resources *rescueResources, reader *bufio.Reader) error {
	printStatus("◐", "🔄 Fixing Collation...")

	answer, err := promptConfirm(reader, "Are you sure you want to fix Collation?")
	if err != nil {
		return err
	}
	if !answer {
		return errors.New("aborted")
	}

	var dbName string
	if err := resources.db.QueryRow(`SELECT current_database()`).Scan(&dbName); err != nil {
		return fmt.Errorf("read current database name: %w", err)
	}

	printStatus("◐", fmt.Sprintf("🔄 Refreshing Collation for database: %s", dbName))

	escapedName := strings.ReplaceAll(dbName, `"`, `""`)
	if _, err := resources.db.Exec(fmt.Sprintf(`ALTER DATABASE "%s" REFRESH COLLATION VERSION`, escapedName)); err != nil {
		return fmt.Errorf("refresh collation version: %w", err)
	}

	printStatus("✔", "✅ Collation fixed successfully.")

	return nil
}

func truncateHwidUserDevices(resources *rescueResources, reader *bufio.Reader) error {
	printStatus("◐", "🔄 Cleaning up HWID Devices...")

	answer, err := promptConfirm(reader, "Are you sure you want to clean up HWID Devices?")
	if err != nil {
		return err
	}
	if !answer {
		return errors.New("aborted")
	}

	if _, err := resources.db.Exec(`TRUNCATE hwid_user_devices`); err != nil {
		return fmt.Errorf("clean up HWID Devices: %w", err)
	}

	printStatus("✔", "✅ HWID Devices cleaned up successfully.")

	return nil
}

func truncateSRHTable(resources *rescueResources, reader *bufio.Reader) error {
	printStatus("◐", "🔄 Cleaning up SRH Table...")

	answer, err := promptConfirm(reader, "Are you sure you want to clean up SRH Table?")
	if err != nil {
		return err
	}
	if !answer {
		return errors.New("aborted")
	}

	if _, err := resources.db.Exec(`TRUNCATE user_subscription_request_history RESTART IDENTITY`); err != nil {
		return fmt.Errorf("clean up SRH Table: %w", err)
	}

	printStatus("✔", "✅ SRH Table cleaned up successfully.")

	return nil
}

func promptConfirm(reader *bufio.Reader, label string) (bool, error) {
	for {
		fmt.Printf("? %s [y/N]: ", label)

		text, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("read confirmation: %w", err)
		}

		switch strings.ToLower(strings.TrimSpace(text)) {
		case "y", "yes", "true", "1":
			return true, nil
		case "", "n", "no", "false", "0":
			fmt.Println("❌ Aborted.")
			return false, nil
		default:
			fmt.Println("Please answer yes or no.")
		}
	}
}

func printCLIBox(title string) {
	width := len([]rune(title)) + 4
	line := strings.Repeat("─", width)
	empty := strings.Repeat(" ", width)

	fmt.Println()
	fmt.Printf(" ╭%s╮\n", line)
	fmt.Printf(" │%s│\n", empty)
	fmt.Printf(" │  %s  │\n", title)
	fmt.Printf(" │%s│\n", empty)
	fmt.Printf(" ╰%s╯\n", line)
	fmt.Println()
}

func printSelectedAction(label string) {
	fmt.Println()
	fmt.Printf("%s Select an action\n", ansi("32", "✔"))
	fmt.Println(label)
}

func printStatus(symbol string, message string) {
	now := time.Now().Format("3:04:05 PM")
	coloredNow := ansi("90", now)
	width := terminalWidth()

	plainLength := len([]rune(symbol)) + 1 + len([]rune(message)) + 1 + len([]rune(now))

	if plainLength >= width {
		fmt.Printf("%s %s %s\n", symbol, message, coloredNow)
		return
	}

	padding := strings.Repeat(" ", width-plainLength)
	fmt.Printf("%s %s%s%s\n", symbol, message, padding, coloredNow)
}

func terminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil && width > 0 {
		return width
	}

	columns := strings.TrimSpace(os.Getenv("COLUMNS"))
	if columns != "" {
		parsedWidth, err := strconv.Atoi(columns)
		if err == nil && parsedWidth > 0 {
			return parsedWidth
		}
	}

	return 120
}

func clearPromptArea(lines int) {
	if lines <= 0 {
		return
	}

	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return
	}

	fmt.Printf("\x1b[%dA", lines)
	fmt.Print("\x1b[J")
}

func ansi(code string, value string) string {
	return "\x1b[" + code + "m" + value + "\x1b[0m"
}

func runSilently(fn func() error) error {
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		return runWithLogSilenced(fn)
	}
	defer devNull.Close()

	stdoutFD := int(os.Stdout.Fd())
	stderrFD := int(os.Stderr.Fd())

	savedStdout, err := syscall.Dup(stdoutFD)
	if err != nil {
		return runWithLogSilenced(fn)
	}
	defer syscall.Close(savedStdout)

	savedStderr, err := syscall.Dup(stderrFD)
	if err != nil {
		return runWithLogSilenced(fn)
	}
	defer syscall.Close(savedStderr)

	oldLogWriter := log.Writer()
	log.SetOutput(io.Discard)

	_ = syscall.Dup2(int(devNull.Fd()), stdoutFD)
	_ = syscall.Dup2(int(devNull.Fd()), stderrFD)

	fnErr := fn()

	_ = syscall.Dup2(savedStdout, stdoutFD)
	_ = syscall.Dup2(savedStderr, stderrFD)

	log.SetOutput(oldLogWriter)

	return fnErr
}

func runWithLogSilenced(fn func() error) error {
	oldLogWriter := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(oldLogWriter)

	return fn()
}
