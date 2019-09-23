package util

import (
	log "MyBlog/pkg/logging"
	"github.com/astaxie/beego/validation"
)

type validAction func(v *validation.Validation)
type validCallback func()

func Valid(va validAction, vc validCallback) {
	valid := validation.Validation{}
	va(&valid)
	if !valid.HasErrors() {
		vc()
	} else {
		for _, err := range valid.Errors {
			log.Error("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
}
