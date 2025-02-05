// SPDX-License-Identifier: AGPL-3.0-only
// Provenance-includes-location: https://github.com/thanos-io/thanos/blob/master/pkg/block/metadata/hash.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: The Thanos Authors.

package metadata

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-kit/log"
	"github.com/pkg/errors"

	"github.com/grafana/dskit/runutil"
)

// HashFunc indicates what type of hash it is.
type HashFunc string

const (
	// SHA256Func shows that SHA256 has been used to generate the hash.
	SHA256Func HashFunc = "SHA256"
	// NoneFunc shows that hashes should not be added. Used internally.
	NoneFunc HashFunc = ""
)

// ObjectHash stores the hash of an object in the object storage.
type ObjectHash struct {
	Func  HashFunc `json:"hashFunc"`
	Value string   `json:"value"`
}

// Equal returns true if two hashes are equal.
func (oh *ObjectHash) Equal(other *ObjectHash) bool {
	return oh.Value == other.Value
}

// CalculateHash calculates the hash of the given type.
func CalculateHash(p string, hf HashFunc, logger log.Logger) (ObjectHash, error) {
	switch hf {
	case SHA256Func:
		f, err := os.Open(filepath.Clean(p))
		if err != nil {
			return ObjectHash{}, errors.Wrap(err, "opening file")
		}
		defer runutil.CloseWithLogOnErr(logger, f, "closing %s", p)

		h := sha256.New()

		if _, err := io.Copy(h, f); err != nil {
			return ObjectHash{}, errors.Wrap(err, "copying")
		}

		return ObjectHash{
			Func:  SHA256Func,
			Value: hex.EncodeToString(h.Sum(nil)),
		}, nil
	}
	return ObjectHash{}, fmt.Errorf("hash function %v is not supported", hf)

}
