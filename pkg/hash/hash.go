package hash

import (
	"crypto"
	"encoding/base64"
	"fmt"
)

// Output formatting wrapper for hash
func Hash(script []byte, hashfns[]crypto.Hash) (map[crypto.Hash]string) {
  hashes := hash(script, hashfns)
  output := map[crypto.Hash]string{}
  for k, v := range hashes {
    output[k] = base64.RawStdEncoding.EncodeToString(v)
  }
  return output
}

// Takes a plaintext, and an array of hashfunctions, and returns the pt hashed with all the hashes
func hash(script []byte, hashfns[]crypto.Hash) (map[crypto.Hash][]byte) {
  hashes := map[crypto.Hash][]byte{}
  for _, fn := range hashfns {
    if !fn.Available() {
      fmt.Printf("%v is not available", fn)
    }
    h := fn.New()
    h.Write(script)
    hashes[fn] = h.Sum(nil)
  }
  return hashes
}
