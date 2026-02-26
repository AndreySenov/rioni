// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package cert

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func Check(certPath string, keyPath string, buildSelfSigned bool) error {
	if certPath == "" || keyPath == "" {
		return errors.New("HTTP TLS requires both cert-file and key-file")
	}

	_, certErr := os.Stat(certPath)
	_, keyErr := os.Stat(keyPath)

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
			return fmt.Errorf("cert-file %q does not exist but key-file %q exists; fix paths or delete the key-file to allow generation", certPath, keyPath)
		}
		return fmt.Errorf("key-file %q does not exist but cert-file %q exists; fix paths or delete the cert-file to allow generation", keyPath, certPath)
	}

	if buildSelfSigned {
		log.Println("Generating self-signed certificate")
		if err := Build(certPath, keyPath); err != nil {
			return fmt.Errorf("failed to build self-signed certificate: %w", err)
		}
	}

	return nil
}
