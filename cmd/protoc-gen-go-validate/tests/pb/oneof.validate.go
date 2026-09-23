package pb

// WithOneof is here to prove the parts of the template that are not field
// rules survive the rename: the oneof switch and the required-oneof check are
// emitted inside the same conditional block as everything else, so a misplaced
// {{ end }} there surfaces as a silently missing rule.

// Validate runs the rules compiled from feat.proto. There is nothing to add to
// them here, which is itself a valid use of the option -- the message opts out
// so that a later change has somewhere to go without an edit to the .proto.
func (x *WithOneof) Validate() error { return x._Validate() }

// ValidateAll mirrors the generated pair, for callers that want every
// violation rather than the first.
func (x *WithOneof) ValidateAll() error { return x._ValidateAll() }
