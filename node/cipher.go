package noise

import "golang.org/x/crypto/blake2b"

// Blake2 return a 32-bytes representation for blake2 hash.
func Blake2(i []byte) []byte {
	hash, err := blake2b.New(blake2b.Size256, nil)
	if err != nil {
		return nil
	}

	hash.Write(i)
	digest := hash.Sum(nil)
	return digest
}
