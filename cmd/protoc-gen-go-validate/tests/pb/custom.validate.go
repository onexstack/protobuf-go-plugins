package pb

import "errors"

// This file is the hand-written half of the (validate.unexported) pair, and it
// is deliberately named so that the generator never writes it: protoc rewrites
// feat.pb.validate.go in full on every run, so anything living there is lost.
//
// The generated file keeps the compiled rules under _Validate/_ValidateAll.
// The exported pair below is what every caller actually reaches -- reqval.Apply
// type-asserts interface{ Validate() error }, and a parent message embedding
// this one asserts the same -- so the compiled rules still run, with these
// extra checks appended after them.

// Validate runs the rules compiled from feat.proto, then the check those rules
// cannot express: the in-list constrains which codes are legal, not which ids
// may accompany them.
func (x *Custom) Validate() error {
	if err := x._Validate(); err != nil {
		return err
	}
	return x.checkIds()
}

// ValidateAll is Validate returning every violation rather than the first, as
// the generated pair does. Both are written because a parent message's
// generated code dispatches on whichever of the two it finds.
func (x *Custom) ValidateAll() error {
	if err := x._ValidateAll(); err != nil {
		return err
	}
	return x.checkIds()
}

func (x *Custom) checkIds() error {
	for _, id := range x.GetIds() {
		if id == "blocked" {
			return errors.New("ids must not contain \"blocked\"")
		}
	}
	return nil
}
