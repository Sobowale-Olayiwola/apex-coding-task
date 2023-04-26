package util

const (
	sats = "sats"
)

// IsSupportedAsset returns true if the asset is supported
func IsSupportedAsset(currency string) bool {
	switch currency {
	case sats:
		return true
	}
	return false
}
