package seed

import (
	"context"
	"database/sql"
	"fmt"

	"exodus/internal/config"
	"exodus/internal/security"

	"github.com/google/uuid"
)

func ensureKeygen(ctx context.Context, tx *sql.Tx, _ *config.BackendConfig) error {
	fmt.Println("◐ Seeding keygen...")

	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM keygen`).Scan(&count); err != nil {
		return fmt.Errorf("count keygen rows: %w", err)
	}

	if count > 1 {
		if _, err := tx.ExecContext(ctx, `DELETE FROM keygen`); err != nil {
			return fmt.Errorf("delete old keygen rows: %w", err)
		}
		count = 0
	}

	fmt.Printf("ℹ Keygen count: %d\n", count)

	if count == 0 {
		pubKey, privKey, err := security.GenerateJWTKeypair()
		if err != nil {
			return fmt.Errorf("generate jwt keypair: %w", err)
		}
		masterCerts, err := security.GenerateMasterCerts()
		if err != nil {
			return fmt.Errorf("generate master certs: %w", err)
		}

		query := `
			INSERT INTO keygen (uuid, priv_key, pub_key, ca_cert, ca_key, client_cert, client_key)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		if _, err := tx.ExecContext(
			ctx,
			query,
			uuid.NewString(),
			privKey,
			pubKey,
			masterCerts.CACertPEM,
			masterCerts.CAKeyPEM,
			masterCerts.ClientCertPEM,
			masterCerts.ClientKeyPEM,
		); err != nil {
			return fmt.Errorf("insert keygen row: %w", err)
		}
		fmt.Println("✔ Keygen seeded")
		return nil
	}

	var (
		id         string
		pubKey     sql.NullString
		privKey    sql.NullString
		caCert     sql.NullString
		caKey      sql.NullString
		clientCert sql.NullString
		clientKey  sql.NullString
	)
	if err := tx.QueryRowContext(
		ctx,
		`SELECT uuid, pub_key, priv_key, ca_cert, ca_key, client_cert, client_key FROM keygen ORDER BY created_at ASC LIMIT 1`,
	).Scan(&id, &pubKey, &privKey, &caCert, &caKey, &clientCert, &clientKey); err != nil {
		return fmt.Errorf("read keygen row: %w", err)
	}

	needJWT := !pubKey.Valid || pubKey.String == "" || !privKey.Valid || privKey.String == ""
	needMTLS := !caCert.Valid || caCert.String == "" || !caKey.Valid || caKey.String == "" || !clientCert.Valid || clientCert.String == "" || !clientKey.Valid || clientKey.String == ""
	if !needJWT && !needMTLS {
		fmt.Println("✔ Keygen seeded")
		return nil
	}

	updateParts := make([]string, 0, 6)
	args := make([]any, 0, 7)
	idx := 1
	if needJWT {
		newPubKey, newPrivKey, err := security.GenerateJWTKeypair()
		if err != nil {
			return fmt.Errorf("regenerate jwt keypair: %w", err)
		}
		updateParts = append(updateParts, fmt.Sprintf("pub_key = $%d", idx), fmt.Sprintf("priv_key = $%d", idx+1))
		args = append(args, newPubKey, newPrivKey)
		idx += 2
	}
	if needMTLS {
		masterCerts, err := security.GenerateMasterCerts()
		if err != nil {
			return fmt.Errorf("regenerate mTLS certificates: %w", err)
		}
		updateParts = append(updateParts, fmt.Sprintf("ca_cert = $%d", idx), fmt.Sprintf("ca_key = $%d", idx+1), fmt.Sprintf("client_cert = $%d", idx+2), fmt.Sprintf("client_key = $%d", idx+3))
		args = append(args, masterCerts.CACertPEM, masterCerts.CAKeyPEM, masterCerts.ClientCertPEM, masterCerts.ClientKeyPEM)
		idx += 4
	}

	query := fmt.Sprintf("UPDATE keygen SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = $%d", joinWithComma(updateParts), idx)
	args = append(args, id)
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("update keygen row: %w", err)
	}

	fmt.Println("✔ Keygen seeded")
	return nil
}

func joinWithComma(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += ", " + parts[i]
	}
	return out
}
