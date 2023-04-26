package api

import (
	"simpledice/util"

	"github.com/go-playground/validator/v10"
)

var validAsset validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if asset, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedAsset(asset)
	}
	return false
}
