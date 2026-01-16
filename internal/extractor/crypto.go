package extractor

import (
	"encoding/hex"
	"strings"
)

type CryptoValidator struct{}

func NewCryptoValidator() *CryptoValidator {
	return &CryptoValidator{}
}

func (v *CryptoValidator) ValidateBitcoin(addr string) bool {

	if len(addr) < 26 || len(addr) > 90 {
		return false
	}

	if !strings.HasPrefix(addr, "1") &&
		!strings.HasPrefix(addr, "3") &&
		!strings.HasPrefix(addr, "bc1") {
		return false
	}

	validChars := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	for _, c := range addr {
		if !strings.ContainsRune(validChars, c) && c != '1' {
			return false
		}
	}

	return true
}

func (v *CryptoValidator) ValidateEthereum(addr string) bool {
	if !strings.HasPrefix(addr, "0x") {
		return false
	}

	if len(addr) != 42 {
		return false
	}

	_, err := hex.DecodeString(addr[2:])
	return err == nil
}

func (v *CryptoValidator) ValidateMonero(addr string) bool {

	if len(addr) != 95 {
		return false
	}

	if !strings.HasPrefix(addr, "4") && !strings.HasPrefix(addr, "8") {
		return false
	}

	validChars := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	for _, c := range addr {
		if !strings.ContainsRune(validChars, c) {
			return false
		}
	}

	return true
}

func IsBitcoinAddress(s string) bool {
	return BitcoinPattern.MatchString(s)
}

func IsEthereumAddress(s string) bool {
	return EthereumPattern.MatchString(s)
}

func IsMoneroAddress(s string) bool {
	return MoneroPattern.MatchString(s)
}
