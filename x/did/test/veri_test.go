package test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/tokenchain/dp-hub/x/did/ed25519"
	"testing"
)

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

type zeroReader struct{}

func (zeroReader) Read(buf []byte) (int, error) {
	for i := range buf {
		buf[i] = 0
	}
	return len(buf), nil
}

func BenchmarkKeyGeneration(b *testing.B) {
	var zero zeroReader
	for i := 0; i < b.N; i++ {
		_, pri, err := ed25519.GenerateKey(zero)
		//require.NotNil(b, err, "key generation error in here %d", i)
		require.Nil(b, err, "key gen %s", pri)
	}
}

func BenchmarkSigning(b *testing.B) {
	var zero zeroReader
	_, priv, err := ed25519.GenerateKey(zero)
	if err != nil {
		b.Fatal(err)
	}
	message := []byte(sample_msg)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ed25519.Sign(priv[:], message)
	}
}

func BenchmarkVerification(b *testing.B) {
	var zero zeroReader
	message := []byte(sample_msg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pub, priv, err := ed25519.GenerateKey(zero)
		if err != nil {
			b.Fatal(err)
		}
		signature := ed25519.Sign(priv[:], message)
		res := ed25519.Verify(pub[:], message, signature)
		fmt.Println("try ", i, " and result is ", res)
	}
}
