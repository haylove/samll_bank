//  *@createTime    2022/3/21 2:53
//  *@author        hay&object
//  *@version       v1.0.0

package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/haylove/small_bank/util"
)

// validCurrency is the validator of currency
var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
