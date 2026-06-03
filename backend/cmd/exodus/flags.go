package exodus

import (
	"flag"
	"fmt"
	"log"
	"os"

	"exodus/internal/config"
	"exodus/internal/constant"
)

type cliFlags struct {
	version            bool
	rescue             bool
	resetAdminPassword bool
}

func parseCLIFlags() cliFlags {
	var flags cliFlags
	flag.BoolVar(&flags.version, "version", false, "Show version information")
	flag.BoolVar(&flags.rescue, "rescue", false, "Open interactive rescue CLI and exit")
	flag.BoolVar(&flags.resetAdminPassword, "reset-admin-password", false, "Interactively reset an existing admin password, enable password auth, revoke sessions, and exit")
	flag.Parse()
	return flags
}

func runPreConfigCLI(flags cliFlags) bool {
	if flags.version {
		fmt.Println(constant.GetBuildInfo())
		return true
	}

	if flags.rescue {
		if err := runRescueCLI(); err != nil {
			log.Printf("Rescue CLI failed: %v", err)
			os.Exit(1)
		}
		return true
	}

	return false
}

func runConfiguredCLI(flags cliFlags, cfg *config.BackendConfig) bool {
	if !flags.resetAdminPassword {
		return false
	}

	if err := runAdminCredentialReset(cfg); err != nil {
		cfg.Logger.Error("Failed to reset admin credentials", "error", err)
		os.Exit(1)
	}
	return true
}

func printRescueHint() {
	fmt.Println("Hint: run `docker exec -it exodus cli` for rescue CLI.")
}
