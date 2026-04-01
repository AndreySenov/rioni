// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package cert

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/AndreySenov/rioni/internal/cfg"
)

func Check(tls cfg.Tls) error {
	_, certErr := os.Stat(tls.CertFile())
	_, keyErr := os.Stat(tls.KeyFile())

	if certErr == nil && keyErr == nil {
		return nil
	}

	if certErr != nil && !os.IsNotExist(certErr) {
		return certErr
	}
	if keyErr != nil && !os.IsNotExist(keyErr) {
		return keyErr
	}

	certMissing := os.IsNotExist(certErr)
	keyMissing := os.IsNotExist(keyErr)

	if certMissing != keyMissing {
		if certMissing {
			return fmt.Errorf("cert-file %q does not exist but key-file %q exists; fix paths or delete the key-file to allow generation", tls.CertFile(), tls.KeyFile())
		}
		return fmt.Errorf("key-file %q does not exist but cert-file %q exists; fix paths or delete the cert-file to allow generation", tls.KeyFile(), tls.CertFile())
	}

	if tls.IsBuildSelfSigned() {
		slog.Info("generating self-signed certificate",
			"cert_file", tls.CertFile(),
			"key_file", tls.KeyFile(),
		)
		if err := Build(tls); err != nil {
			return fmt.Errorf("failed to build self-signed certificate: %w", err)
		}
	}

	return nil
}
