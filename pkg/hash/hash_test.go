package hash

import (
	"crypto"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"testing"
)

var (
	input   = "test-data"
	hash256 = "sha256-oYYABCL+q4VzKcaE6f6RQSsaXbCEEAs3qYz8lbYqqGc="
	hash384 = "sha384-vnqLl9XO0CW4TZyma7tvp0uml21xQDQmSDV3NngFxQLt/Njwq8uKuaO6VOHiigwM"
	hash512 = "sha512-uFFwBg3iPS4GYj4luRB4uIknI3tsfV7Qd7mJjLh2txcd+lfxPMA1Wso1mTuFAhjPC2nRFncuI81dW5a+MEjp/A=="
)

func TestSingleSHA512(t *testing.T) {
	output := Hash([]byte(input), []crypto.Hash{crypto.SHA512})
	res_hash := output[crypto.SHA512]

	if hash512 != res_hash {
		t.Fatalf("Expected %v to hash to %v, but got %v", input, hash512, res_hash)
	}
}

func TestSingleSHA256(t *testing.T) {
	output := Hash([]byte(input), []crypto.Hash{crypto.SHA256})
	res_hash := output[crypto.SHA256]

	if hash256 != res_hash {
		t.Fatalf("Expected %v to hash to %v, but got %v", input, hash256, res_hash)
	}
}

func TestSingleSHA384(t *testing.T) {
	output := Hash([]byte(input), []crypto.Hash{crypto.SHA384})
	res_hash := output[crypto.SHA384]

	if hash384 != res_hash {
		t.Fatalf("Expected %v to hash to %v, but got %v", input, hash384, res_hash)
	}
}

func TestMultipleSHA(t *testing.T) {
	//TODO to be sure that we get no weirdness with the data, this should probably be data from an actual JS file
	exp_hashes := map[crypto.Hash]string{
		crypto.SHA256: "sha256-oYYABCL+q4VzKcaE6f6RQSsaXbCEEAs3qYz8lbYqqGc=",
		crypto.SHA384: "sha384-vnqLl9XO0CW4TZyma7tvp0uml21xQDQmSDV3NngFxQLt/Njwq8uKuaO6VOHiigwM",
		crypto.SHA512: "sha512-uFFwBg3iPS4GYj4luRB4uIknI3tsfV7Qd7mJjLh2txcd+lfxPMA1Wso1mTuFAhjPC2nRFncuI81dW5a+MEjp/A==",
	}
	output := Hash([]byte(input), []crypto.Hash{crypto.SHA256, crypto.SHA384, crypto.SHA512})

	for k, v := range output {
		expected := exp_hashes[k]
		if expected != v {
			t.Fatalf("Expected %v to hash as %v with %v, but got %v", input, expected, k, v)
		}
	}
}
