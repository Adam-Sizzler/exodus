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
	case "truncate-hwid-user-devices":
		return truncateHwidUserDevices(resources, reader)
	case "truncate-srh-table":
		return truncateSRHTable(resources, reader)
	case "truncate-users-usage-table":
		return truncateUsersUsageTable(resources, reader)
	case "delete-users-usage-by-date-range":
		return deleteUsersUsageByDateRange(resources, reader)
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

	return promptSelect(actions, len(actions)-1)
}

func promptSelect(actions []cliAction, initialIndex int) (string, error) {
	if len(actions) == 0 {
		return "", errors.New("no options configured")
	}

	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		return promptSelectPlain(actions, initialIndex)
	}

	selected := initialIndex
	if selected < 0 || selected >= len(actions) {
		selected = len(actions) - 1
	}

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
			Value: "truncate-users-usage-table",
			Label: "Clean up Users Usage Table",
			Hint:  "Remove all users traffic statistics data from the database",
		},
		{
			Value: "delete-users-usage-by-date-range",
			Label: "Delete Users Usage by date range",
			Hint:  "Remove traffic statistics for a period (day-month-year); choose single or batched",
		},
		{
			Value: "exit",
			Label: "Exit",
		},
	}
}

func promptSelectPlain(actions []cliAction, initialIndex int) (string, error) {
	if initialIndex < 0 || initialIndex >= len(actions) {
		initialIndex = len(actions) - 1
	}

	fmt.Println()
	fmt.Println("Select an action:")
	for index, action := range actions {
		if action.Hint != "" {
			fmt.Printf("%d) %s (%s)\n", index+1, action.Label, action.Hint)
			continue
		}

		fmt.Printf("%d) %s\n", index+1, action.Label)
	}

	fmt.Printf("Enter number [%d]: ", initialIndex+1)

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read action: %w", err)
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return actions[initialIndex].Value, nil
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

	err = resources.db.QueryRow(`
		SELECT uuid, username
		FROM admin
		ORDER BY created_at ASC
		LIMIT 1
	`).Scan(&adminUUID, &username)

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("superadmin not found")
	}
	if err != nil {
		return fmt.Errorf("find superadmin: %w", err)
	}

	if _, err := resources.db.Exec("DELETE FROM admin WHERE uuid = $1", adminUUID); err != nil {
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

	if _, err := resources.db.Exec("DELETE FROM keygen WHERE uuid = $1", keygenUUID); err != nil {
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

func truncateUsersUsageTable(resources *rescueResources, reader *bufio.Reader) error {
	printStatus("◐", "🔄 Cleaning up Users Usage Table...")

	answer, err := promptConfirm(reader, "Are you sure you want to clean up Users Usage Table?")
	if err != nil {
		return err
	}
	if !answer {
		return errors.New("aborted")
	}

	if _, err := resources.db.Exec(`TRUNCATE nodes_user_usage_history RESTART IDENTITY`); err != nil {
		return fmt.Errorf("clean up Users Usage Table: %w", err)
	}
	if _, err := resources.db.Exec(`VACUUM nodes_user_usage_history`); err != nil {
		return fmt.Errorf("vacuum Users Usage Table: %w", err)
	}
	if _, err := resources.db.Exec(`REINDEX TABLE nodes_user_usage_history`); err != nil {
		return fmt.Errorf("reindex Users Usage Table: %w", err)
	}

	printStatus("✔", "✅ Users Usage Table cleaned up successfully.")

	return nil
}

const dateInputLayout = "02-01-2006" // strict DD-MM-YYYY

func promptStrictDate(reader *bufio.Reader, label string, example string) (time.Time, error) {
	fmt.Printf(
		"Enter the %s date in strict format day-month-year (DD-MM-YYYY), e.g. %s: ",
		label, example,
	)

	text, err := reader.ReadString('\n')
	if err != nil {
		return time.Time{}, fmt.Errorf("read %s date: %w", label, err)
	}

	parsed, err := time.Parse(dateInputLayout, strings.TrimSpace(text))
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"invalid %s date: expected strict format DD-MM-YYYY, e.g. %s", label, example,
		)
	}

	return parsed, nil
}

func promptBatchSize(reader *bufio.Reader) (int, error) {
	fmt.Print("Batch size (rows per DELETE) [50000]: ")

	text, err := reader.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("read batch size: %w", err)
	}

	text = strings.TrimSpace(text)
	if text == "" {
		text = "50000"
	}

	batchSize, err := strconv.Atoi(text)
	if err != nil || batchSize <= 0 {
		return 0, fmt.Errorf("invalid batch size: expected a positive integer")
	}

	return batchSize, nil
}

func formatThousands(n int64) string {
	s := strconv.FormatInt(n, 10)

	neg := strings.HasPrefix(s, "-")
	if neg {
		s = s[1:]
	}

	var out []byte
	for i, c := range []byte(s) {
		if i != 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, c)
	}

	if neg {
		return "-" + string(out)
	}

	return string(out)
}

func renderDeleteProgress(current, total int64, startedAt time.Time, lastBatchMs int64) {
	ratio := 1.0
	if total > 0 {
		ratio = float64(current) / float64(total)
		if ratio > 1 {
			ratio = 1
		}
	}

	const width = 28
	filled := int(ratio*width + 0.5)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	pct := fmt.Sprintf("%5.1f%%", ratio*100)

	elapsedSec := time.Since(startedAt).Seconds()
	etaStr := "—"
	if current > 0 && total > current {
		etaStr = fmt.Sprintf("%.1f", (elapsedSec/float64(current))*float64(total-current))
	}

	line := fmt.Sprintf(
		"  [%s] %s  %s/%s  | %.1fs elapsed | ETA %ss | last %dms | do NOT close",
		bar, pct, formatThousands(current), formatThousands(total), elapsedSec, etaStr, lastBatchMs,
	)

	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Print("\r\x1b[K")
		fmt.Print(line)
	} else {
		fmt.Println(strings.TrimSpace(line))
	}
}

func runSingleUsageDelete(resources *rescueResources, startStr, endStr string) (int64, error) {
	printStatus("◐", "🔄 Deleting records... (do NOT close this window)")

	result, err := resources.db.Exec(`
		DELETE FROM nodes_user_usage_history
		WHERE created_at >= $1::date
		  AND created_at <= $2::date
	`, startStr, endStr)
	if err != nil {
		return 0, fmt.Errorf("delete records: %w", err)
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("read rows affected: %w", err)
	}

	printStatus("✔", fmt.Sprintf("✅ Deleted %s record(s).", formatThousands(deleted)))

	return deleted, nil
}

func runBatchedUsageDelete(
	resources *rescueResources,
	startStr, endStr string,
	batchSize int,
	totalToDelete int64,
	startedAt time.Time,
) (int64, error) {
	var totalDeleted int64
	var timings []int64

	printStatus("◐", "🔄 Deleting records in batches... (do NOT close this window)")

	for {
		batchStart := time.Now()

		result, err := resources.db.Exec(`
			DELETE FROM nodes_user_usage_history
			WHERE ctid IN (
				SELECT ctid
				FROM nodes_user_usage_history
				WHERE created_at >= $1::date
				  AND created_at <= $2::date
				LIMIT $3
			)
		`, startStr, endStr, batchSize)
		if err != nil {
			if term.IsTerminal(int(os.Stdout.Fd())) {
				fmt.Println()
			}
			return totalDeleted, fmt.Errorf("delete records: %w", err)
		}

		batchMs := time.Since(batchStart).Milliseconds()

		deleted, err := result.RowsAffected()
		if err != nil {
			if term.IsTerminal(int(os.Stdout.Fd())) {
				fmt.Println()
			}
			return totalDeleted, fmt.Errorf("read rows affected: %w", err)
		}

		if deleted == 0 {
			break
		}

		totalDeleted += deleted
		timings = append(timings, batchMs)

		renderDeleteProgress(totalDeleted, totalToDelete, startedAt, batchMs)
	}

	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println()
	}

	if len(timings) > 0 {
		first := timings[0]
		last := timings[len(timings)-1]
		min, max, sum := timings[0], timings[0], int64(0)
		for _, t := range timings {
			if t < min {
				min = t
			}
			if t > max {
				max = t
			}
			sum += t
		}
		avg := sum / int64(len(timings))

		printInfoBox([]string{
			"Batched delete summary",
			fmt.Sprintf("batches:            %d", len(timings)),
			fmt.Sprintf("batch size:         %s", formatThousands(int64(batchSize))),
			fmt.Sprintf("deleted total:      %s", formatThousands(totalDeleted)),
			fmt.Sprintf("first / avg / last: %dms / %dms / %dms", first, avg, last),
			fmt.Sprintf("min / max:          %dms / %dms", min, max),
		})
	}

	return totalDeleted, nil
}

func deleteUsersUsageByDateRange(resources *rescueResources, reader *bufio.Reader) error {
	fmt.Println(
		"This will permanently delete users traffic statistics " +
			"(nodes_user_usage_history) for the selected date range.",
	)

	method, err := promptSelect([]cliAction{
		{
			Value: "single",
			Label: "Single query (fast)",
			Hint:  "One DELETE — fastest overall, but holds one longer lock",
		},
		{
			Value: "batched",
			Label: "Batched (low-lock + progress bar)",
			Hint:  "Many small DELETEs — shorter locks, live progress, slower overall",
		},
	}, 0)
	if err != nil {
		return err
	}

	startDate, err := promptStrictDate(reader, "START", "01-01-2024")
	if err != nil {
		return err
	}

	endDate, err := promptStrictDate(reader, "END", "31-12-2024")
	if err != nil {
		return err
	}

	if endDate.Before(startDate) {
		return fmt.Errorf("END date can not be earlier than START date")
	}

	batchSize := 0
	if method == "batched" {
		batchSize, err = promptBatchSize(reader)
		if err != nil {
			return err
		}
	}

	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")

	printStatus("◐", "🔍 Counting affected rows...")

	var rowsToDelete int64
	if err := resources.db.QueryRow(`
		SELECT COUNT(*)
		FROM nodes_user_usage_history
		WHERE created_at >= $1::date
		  AND created_at <= $2::date
	`, startStr, endStr).Scan(&rowsToDelete); err != nil {
		return fmt.Errorf("count rows: %w", err)
	}

	if rowsToDelete == 0 {
		printStatus("ℹ", fmt.Sprintf(
			"ℹ️ No records found between %s and %s (inclusive). Nothing to delete.", startStr, endStr,
		))
		return nil
	}

	methodLine := "in a single query."
	if method == "batched" {
		methodLine = fmt.Sprintf("in batches of %s.", formatThousands(int64(batchSize)))
	}

	printInfoBox([]string{
		fmt.Sprintf("About to delete %s record(s)", formatThousands(rowsToDelete)),
		fmt.Sprintf("from %s to %s (inclusive)", startStr, endStr),
		"from table \"nodes_user_usage_history\" " + methodLine,
	})

	fmt.Println(
		"⚠ Do NOT close this window until the operation is finished.\n" +
			"⚠ A final VACUUM runs at the end to reclaim space.\n" +
			"⚠ Interrupting the operation may leave the table bloated.",
	)

	answer, err := promptConfirm(reader, fmt.Sprintf(
		"Are you sure you want to permanently delete these %s record(s)?", formatThousands(rowsToDelete),
	))
	if err != nil {
		return err
	}
	if !answer {
		return errors.New("aborted")
	}

	startedAt := time.Now()

	var totalDeleted int64
	if method == "batched" {
		totalDeleted, err = runBatchedUsageDelete(resources, startStr, endStr, batchSize, rowsToDelete, startedAt)
	} else {
		totalDeleted, err = runSingleUsageDelete(resources, startStr, endStr)
	}
	if err != nil {
		return err
	}

	printStatus("◐", "🧹 Reclaiming space (VACUUM)... (do NOT close this window)")
	if _, err := resources.db.Exec(`VACUUM nodes_user_usage_history`); err != nil {
		printStatus("⚠", fmt.Sprintf("⚠️ Final VACUUM failed (table left as-is): %v", err))
	}

	elapsedSec := time.Since(startedAt).Seconds()
	printStatus("✔", fmt.Sprintf(
		"✅ Done in %.1fs. Removed %s record(s) from %s to %s.",
		elapsedSec, formatThousands(totalDeleted), startStr, endStr,
	))

	return nil
}

func printInfoBox(lines []string) {
	width := 0
	for _, line := range lines {
		if l := len([]rune(line)); l > width {
			width = l
		}
	}
	width += 4

	border := strings.Repeat("─", width)

	fmt.Println()
	fmt.Printf(" ╭%s╮\n", border)
	for _, line := range lines {
		padding := strings.Repeat(" ", width-2-len([]rune(line)))
		fmt.Printf(" │  %s%s │\n", line, padding)
	}
	fmt.Printf(" ╰%s╯\n", border)
	fmt.Println()
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
