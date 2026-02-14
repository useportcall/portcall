//go:build integration

package checkout_session_test

import "github.com/useportcall/portcall/libs/go/cryptox"

// noopCrypto satisfies cryptox.ICrypto without real encryption.
// Only needed because the service constructor requires it.
type noopCrypto struct{}

var _ cryptox.ICrypto = (*noopCrypto)(nil)

func (noopCrypto) Encrypt(data string) (string, error)            { return data, nil }
func (noopCrypto) Decrypt(data string) (string, error)            { return data, nil }
func (noopCrypto) CompareHash(hashed, plain string) (bool, error) { return hashed == plain, nil }
