package shared

import (
	pgs "github.com/lyft/protoc-gen-star/v2"

	"github.com/onexstack/protobuf-go-plugins/validate"
)

// Disabled returns true if validations are disabled for msg
func Disabled(msg pgs.Message) (disabled bool, err error) {
	_, err = msg.Extension(validate.E_Disabled, &disabled)
	return
}

// Ignore returns true if validations aren't to be generated for msg
func Ignored(msg pgs.Message) (ignored bool, err error) {
	_, err = msg.Extension(validate.E_Ignored, &ignored)
	return
}

// Unexported returns true if the rule body should be emitted under the
// underscore-prefixed names instead of the exported ones.
//
// It is the validate-side counterpart of defaults' (defaults.unexported): the
// generated file keeps everything the rules compile to, but the two entry
// points are renamed _Validate and _ValidateAll so that a hand-written file can
// supply Validate and ValidateAll, call through to the generated body, and add
// checks of its own. protoc rewrites the whole generated file on every run, so
// without this the hand-written half would be lost.
//
// It differs from Ignored, which emits nothing at all for the message and
// therefore leaves the hand-written method with no compiled rules to call.
func Unexported(msg pgs.Message) (unexported bool, err error) {
	_, err = msg.Extension(validate.E_Unexported, &unexported)
	return
}

// RequiredOneOf returns true if the oneof field requires a field to be set
func RequiredOneOf(oo pgs.OneOf) (required bool, err error) {
	_, err = oo.Extension(validate.E_Required, &required)
	return
}
