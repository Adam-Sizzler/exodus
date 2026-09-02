package cmd

import (
	"bufio"
	"fmt"

	"filippo.io/age"
)

func generateEncryptionKeys(_ *rescueResources, _ *bufio.Reader) error {
	methodActions := []cliAction{
		{
			Value: "age1",
			Label: "age1 (X25519)",
			Hint:  "Native X25519 — classical security",
		},
		{
			Value: "age1pq1",
			Label: "age1pq1 (hybrid post-quantum)",
			Hint:  "X25519 + ML-KEM-768 — post-quantum resistant",
		},
	}

	method, err := promptSelect(methodActions, 0)
	if err != nil {
		return err
	}

	printStatus("◐", fmt.Sprintf("🔑 Generating %s key pair...", method))

	var (
		identityStr  string
		recipientStr string
	)

	switch method {
	case "age1":
		id, err := age.GenerateX25519Identity()
		if err != nil {
			return fmt.Errorf("failed to generate X25519 key pair: %w", err)
		}
		identityStr = id.String()
		recipientStr = id.Recipient().String()
	case "age1pq1":
		id, err := age.GenerateHybridIdentity()
		if err != nil {
			return fmt.Errorf("failed to generate hybrid post-quantum key pair: %w", err)
		}
		identityStr = id.String()
		recipientStr = id.Recipient().String()
	default:
		return fmt.Errorf("unsupported encryption method: %s", method)
	}

	printStatus("✔", fmt.Sprintf("✅ %s key pair generated successfully.", method))
	fmt.Printf("\nPUBLIC KEY (recipient) — put this into the response rule \"encryption.key\":\n%s\n", recipientStr)
	fmt.Printf("\nPRIVATE KEY (identity) — keep it secret, the client uses it to decrypt the response:\n%s\n\n", identityStr)

	return nil
}
