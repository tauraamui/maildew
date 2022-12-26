package configdef

import (
	"gopkg.in/dealancer/validate.v2"
)

type Values struct {
	Debug   bool   `json:"debug"`
	RootKey string `json:"root_key"`
}

func (v Values) RunValidate() error {
	return v.runValidate()
}

func (v Values) runValidate() error {
	// const validationErrorHeader = "validation failed: %w"
	// defaultPersistLocToDot(v.Cameras)
	// if hasDupCameraTitles(v.Cameras) {
	// 	return xerror.Errorf(validationErrorHeader, xerror.New("camera titles must be unique"))
	// }
	return validate.Validate(&v)
}
