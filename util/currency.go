//  *@createTime    2022/3/21 2:49
//  *@author        hay&object
//  *@version       v1.0.0

package util

//contains all supported currency
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	RMB = "YUAN"
)

//IsSupportedCurrency return currency is supported or not
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, RMB:
		return true
	}
	return false
}
