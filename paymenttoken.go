package paypalutil

import "strings"

type (
	PaymentToken       string
	PaymentTokenName   string
	PaymentTokenPrefix string
)

const (
	OrderId          PaymentTokenName = "OrderId"
	BillingAgreement PaymentTokenName = "BillingAgreement"
	ExpressCheckout  PaymentTokenName = "ExpressCheckout"
)

const (
	BA PaymentTokenPrefix = "BA"
	EC PaymentTokenPrefix = "EC"
)

func (s PaymentToken) String() string {
	if s == "" {
		return ""
	}
	return string(s)
}

func (s PaymentToken) Bytes() []byte {
	if s == "" {
		return nil
	}
	return []byte(s)
}

func (s PaymentTokenName) String() string {
	if s == "" {
		return ""
	}
	return string(s)
}

func (s PaymentTokenName) Bytes() []byte {
	if s == "" {
		return nil
	}
	return []byte(s)
}

func (s PaymentTokenPrefix) String() string {
	if s == "" {
		return ""
	}
	return string(s)
}

func (s PaymentTokenPrefix) Bytes() []byte {
	if s == "" {
		return nil
	}
	return []byte(s)
}

func (s PaymentToken) IsValid() bool {
	if length := len(s); length < 17 || length > 20 {
		return false
	}
	before, after, found := strings.Cut(s.String(), "-")
	if found {
		switch t := PaymentTokenPrefix(before); t {
		case BA, EC:
			return isUpperNumber(after)
		default:
			return false
		}
	}
	return isUpperNumber(before)
}

func (s PaymentToken) GetDetails() (PaymentTokenName, PaymentTokenPrefix, bool) {
	if length := len(s); length < 17 || length > 20 {
		return "", "", false
	}
	before, after, found := strings.Cut(s.String(), "-")
	if found {
		if !isUpperNumber(after) {
			return "", "", false
		}
		switch t := PaymentTokenPrefix(before); t {
		case BA:
			return BillingAgreement, BA, true
		case EC:
			return ExpressCheckout, EC, true
		default:
			return "", "", false
		}
	}
	if isUpperNumber(before) {
		return OrderId, "", true
	}
	return "", "", false
}

func isUpperNumber(s string) bool {
	length := len(s)
	if length == 0 {
		return false
	}
	for i := 0; i < length; i++ {
		if 0x41 > s[i] || s[i] > 0x5A {
			if s[i] < '0' || s[i] > '9' {
				return false
			}
		}
	}
	return true
}
