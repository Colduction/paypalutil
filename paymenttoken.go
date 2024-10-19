package paypalutil

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	PaymentTokenInvalidFormatError string
	PaymentTokenInvalidPrefixError string
	PaymentTokenInvalidSuffixError string
	PaymentTokenInvalidSizeError   string
)

func (s PaymentTokenInvalidFormatError) Error() string {
	if s == "" {
		return "paypalutil: payment token cannot be empty"
	}
	return fmt.Sprintf("paypalutil: payment token %s has an invalid format", strconv.Quote(string(s)))
}

func (s PaymentTokenInvalidPrefixError) Error() string {
	if s == "" {
		return "paypalutil: payment token prefix cannot be empty"
	}
	return fmt.Sprintf("paypalutil: payment token %s has an invalid prefix", strconv.Quote(string(s)))
}

func (s PaymentTokenInvalidSuffixError) Error() string {
	if s == "" {
		return "paypalutil: payment token suffix cannot be empty"
	}
	return fmt.Sprintf("paypalutil: payment token %s has an invalid suffix", strconv.Quote(string(s)))
}

func (s PaymentTokenInvalidSizeError) Error() string {
	if s == "" {
		return "paypalutil: payment token cannot be empty"
	}
	return fmt.Sprintf("paypalutil: payment token %s has an invalid size (length: %d)", strconv.Quote(string(s)), len(s))
}

type paymentTokenDetails struct {
	Type   string `json:"type,omitempty"`
	Prefix string `json:"prefix,omitempty"`
}

type PaymentTokenDetailsProvider interface {
	GetType() string
	GetTypeBytes() []byte
	GetPrefix() string
	GetPrefixBytes() []byte
	IsZero() bool
}

func (m paymentTokenDetails) GetType() string {
	return m.Type
}

func (m paymentTokenDetails) GetTypeBytes() []byte {
	return []byte(m.Type)
}

func (m paymentTokenDetails) GetPrefix() string {
	return m.Prefix
}

func (m paymentTokenDetails) GetPrefixBytes() []byte {
	return []byte(m.Prefix)
}

func (m paymentTokenDetails) IsZero() bool {
	return paymentTokenDetails{} == m
}

var (
	billingAgreement paymentTokenDetails = paymentTokenDetails{"Billing Agreement", "BA"}
	expressCheckout  paymentTokenDetails = paymentTokenDetails{"Express Checkout", "EC"}
	orderID          paymentTokenDetails = paymentTokenDetails{"Order ID", ""}
)

var lookupPTPrefix = map[string]paymentTokenDetails{
	billingAgreement.Prefix: billingAgreement,
	expressCheckout.Prefix:  expressCheckout,
	orderID.Prefix:          orderID,
}

type PaymentToken string

func (p PaymentToken) IsValidFormat() error {
	return isValidFormat(string(p))
}

func (p PaymentToken) String() string {
	return string(p)
}

func (p PaymentToken) Bytes() []byte {
	return []byte(p)
}

func (p PaymentToken) IsZero() bool {
	return p == ""
}

type paymentToken struct {
	Details paymentTokenDetails `json:"details,omitempty"`
	Token   PaymentToken        `json:"token,omitempty"`
}

type PaymentTokenProvider interface {
	GetDetails() PaymentTokenDetailsProvider
	GetToken() string
	GetTokenBytes() []byte
	IsZero() bool
}

func NewPaymentToken(token string) (PaymentTokenProvider, error) {
	p, err := newPaymentToken(token)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func newPaymentToken(token string) (paymentToken, error) {
	if length := len(token); length < 17 || length > 20 {
		return paymentToken{}, PaymentTokenInvalidSizeError(token)
	}
	before, after, exists := strings.Cut(token, "-")
	if exists {
		pttype, exists := lookupPTPrefix[before]
		if !exists {
			return paymentToken{}, PaymentTokenInvalidPrefixError(token)
		}
		if !isUpperNumber(after) {
			return paymentToken{}, PaymentTokenInvalidSuffixError(token)
		}
		return paymentToken{Details: pttype, Token: PaymentToken(token)}, nil
	}
	if !isUpperNumber(before) {
		return paymentToken{}, PaymentTokenInvalidFormatError(token)
	}
	return paymentToken{Details: orderID, Token: PaymentToken(token)}, nil
}

func (m paymentToken) GetToken() string {
	return string(m.Token)
}

func (m paymentToken) GetTokenBytes() []byte {
	return []byte(m.Token)
}

func (m paymentToken) GetDetails() PaymentTokenDetailsProvider {
	return m.Details
}

func (m paymentToken) IsZero() bool {
	return paymentToken{} == m
}

func isValidFormat(token string) error {
	if length := len(token); length < 17 || length > 20 {
		return PaymentTokenInvalidSizeError(token)
	}
	prefix, suffix := getPTPrefixSuffix(token)
	if prefix != "" {
		if _, exists := lookupPTPrefix[prefix]; !exists {
			return PaymentTokenInvalidPrefixError(token)
		}
		if !isUpperNumber(suffix) {
			return PaymentTokenInvalidSuffixError(token)
		}
		return nil
	}
	if !isUpperNumber(token) {
		return PaymentTokenInvalidFormatError(token)
	}
	return nil
}

func getPTPrefixSuffix(token string) (prefix string, suffix string) {
	before, after, exists := strings.Cut(token, "-")
	if exists {
		return before, after
	}
	return "", ""
}
