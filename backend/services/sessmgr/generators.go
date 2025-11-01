// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sessmgr

import (
	"encoding/base64"
	"encoding/binary"
	"math/rand/v2"
)

// generateSessionId creates a Base64 URL–encoded string from 32 bytes of
// non-cryptographic random data. Intended for demos and tests only.
//
// ⚠️ Not secure! Use crypto/rand for production systems.
func generateSessionId() string {
	id := make([]byte, 32)
	binary.LittleEndian.PutUint64(id[0*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[1*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[2*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[3*8:], rand.Uint64())
	return base64.RawURLEncoding.EncodeToString(id)
}

// Example secure version (production):
// func generateSecureSessionId() string {
// 	id := make([]byte, 32)
// 	if _, err := crypto/rand.Read(id); err != nil {
// 		panic(err)
// 	}
// 	return base64.RawURLEncoding.EncodeToString(id)
// }
