package exodus

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"exodus/internal/config"
	"exodus/internal/constant"
)

func runRescueCLI() error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("=== Exodus Rescue CLI ===")
		fmt.Println("1. Reset admin password")
		fmt.Println("2. Show version")
		fmt.Println("0. Exit")

		choice, err := promptLine(reader, "Select action", "1", false)
		if err != nil {
			return err
		}

		switch strings.ToLower(strings.TrimSpace(choice)) {
		case "1", "reset", "reset-admin", "reset-admin-password":
			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("load configuration: %w", err)
			}
			return runAdminCredentialReset(&cfg)
		case "2", "version":
			fmt.Println(constant.GetBuildInfo())
		case "0", "q", "quit", "exit":
			return nil
		default:
			fmt.Printf("Unknown action: %s\n", choice)
		}
	}
}
