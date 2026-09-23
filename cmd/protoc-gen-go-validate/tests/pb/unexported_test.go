package pb

import (
	"strings"
	"testing"
)

// The properties under test are exactly the ones the option exists to provide,
// so they are asserted against the real generated fixture rather than a stub:
// that the compiled rules still run after the rename, that the hand-written
// half runs after them, and that a message without the option is untouched.

func TestUnexportedKeepsCompiledRules(t *testing.T) {
	// code is not a member of the in-list, and ids carries one element so that
	// min_items is satisfied. The only rule that can reject this is the
	// compiled one, reached through the hand-written Validate.
	err := (&Custom{Code: "cc", Ids: []string{"ok"}}).Validate()
	if err == nil {
		t.Fatal("compiled in-list rule did not run through the hand-written Validate")
	}
	if !strings.Contains(err.Error(), "must be in list") {
		t.Errorf("want the compiled in-list violation, got: %v", err)
	}
}

func TestUnexportedKeepsRepeatedRules(t *testing.T) {
	// A message-level rule rather than a string rule, to catch a template edit
	// that only preserved part of the rule block.
	err := (&Custom{Code: "aa"}).Validate()
	if err == nil {
		t.Fatal("compiled min_items rule did not run through the hand-written Validate")
	}
	if !strings.Contains(err.Error(), "at least 1 item") {
		t.Errorf("want the compiled min_items violation, got: %v", err)
	}
}

func TestHandWrittenRuleRuns(t *testing.T) {
	// "aa" satisfies the in-list and the single id satisfies min_items, so the
	// compiled rules pass and the failure can only come from the hand-written
	// check. This is the case that distinguishes the option from Ignored: the
	// generated rules ran first and this still got a chance to run.
	err := (&Custom{Code: "aa", Ids: []string{"blocked"}}).Validate()
	if err == nil {
		t.Fatal("hand-written rule did not run")
	}
	if !strings.Contains(err.Error(), "blocked") {
		t.Errorf("want the hand-written violation, got: %v", err)
	}
}

func TestBothHalvesAcceptValidInput(t *testing.T) {
	if err := (&Custom{Code: "aa", Ids: []string{"ok"}}).Validate(); err != nil {
		t.Fatalf("valid input rejected: %v", err)
	}
}

func TestValidateAllIsWired(t *testing.T) {
	// ValidateAll is the second of the pair, and a parent message's generated
	// code dispatches on it. A message that kept only Validate would compile
	// here and fail there.
	if err := (&Custom{Code: "cc", Ids: []string{"ok"}}).ValidateAll(); err == nil {
		t.Fatal("ValidateAll did not run the compiled rules")
	}
}

func TestUnmarkedMessageIsUnchanged(t *testing.T) {
	// The control. Plain carries no option, so its exported pair is the
	// generated one and behaves as PGV has always produced.
	if err := (&Plain{Name: ""}).Validate(); err == nil {
		t.Fatal("Plain lost its compiled rule")
	}
	if err := (&Plain{Name: "ok"}).Validate(); err != nil {
		t.Fatalf("Plain rejected valid input: %v", err)
	}
}

func TestOneofRulesSurviveTheRename(t *testing.T) {
	// The oneof switch and the required-oneof check live inside the same
	// conditional block as the field rules. They are checked here because a
	// misplaced {{ end }} in that block surfaces as a missing rule rather than
	// as a syntax error when the message happens to have no fields.
	if err := (&WithOneof{}).Validate(); err == nil {
		t.Fatal("required-oneof rule did not survive the rename")
	}
	if err := (&WithOneof{Pick: &WithOneof_A{A: "x"}}).Validate(); err != nil {
		t.Fatalf("valid oneof rejected: %v", err)
	}
}
