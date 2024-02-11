package util

// Supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

// isSupportedCurrency checks if the currency is supported.
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}
