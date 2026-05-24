package downloader

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"fmt"
)

func ComputeSHA256(
	path string,
) (string, error) {

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()
	hasher := sha256.New()

	_, err = io.Copy(
		hasher,
		file,
	)

	if err != nil {
		return "", err
	}

	hash := hasher.Sum(nil)

	return hex.EncodeToString(
		hash,
	), nil
}

func ValidateSHA256(
	path string,
	expected string,
) error {

	actual, err := ComputeSHA256(path)

	if err != nil {
		return err
	}

	if actual != expected {
		return fmt.Errorf(
			"sha256 mismatch\nexpected: %s\nactual: %s",
			expected,
			actual,
		)
	}

	return nil
}