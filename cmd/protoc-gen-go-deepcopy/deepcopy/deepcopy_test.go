// SPDX-FileCopyrightText: Copyright 2021 The protobuf-tools Authors
// SPDX-License-Identifier: BSD-3-Clause

package deepcopy

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// generate runs the plugin over a single synthetic file and returns the
// generated source. Building the descriptor here rather than from testdata
// keeps the test independent of a protoc toolchain.
func generate(t *testing.T, messages ...string) string {
	t.Helper()

	descriptors := make([]*descriptorpb.DescriptorProto, 0, len(messages))
	for _, name := range messages {
		descriptors = append(descriptors, &descriptorpb.DescriptorProto{Name: proto.String(name)})
	}

	const filename = "test/v1/test.proto"
	file := &descriptorpb.FileDescriptorProto{
		Name:    proto.String(filename),
		Package: proto.String("test.v1"),
		Syntax:  proto.String("proto3"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("example.com/test/v1;v1"),
		},
		MessageType: descriptors,
	}

	gen, err := protogen.Options{}.New(&pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{filename},
		ProtoFile:      []*descriptorpb.FileDescriptorProto{file},
	})
	if err != nil {
		t.Fatalf("protogen.Options.New: %v", err)
	}

	if len(gen.Files) != 1 {
		t.Fatalf("got %d files, want 1", len(gen.Files))
	}
	// Deliberately not asserting on the return: GenerateFile returns nil for a
	// file that declares no messages, which TestNoMessagesNoFile covers.
	GenerateFile(gen, gen.Files[0])

	for _, f := range gen.Response().GetFile() {
		if strings.HasSuffix(f.GetName(), FileNameSuffix) {
			return f.GetContent()
		}
	}

	// No output at all is a legitimate result -- see TestNoMessagesNoFile --
	// so this reports it as empty rather than failing.
	return ""
}

// TestDeepCopyIntoDoesNotAssignTheStruct is the regression test for the defect
// this fork exists to fix.
//
// The upstream template ended DeepCopyInto in `*out = *p`. That is a struct
// assignment, and a generated message embeds protoimpl.MessageState, which
// holds a sync.Mutex -- so the copy and the original shared one mutex, and
// `go vet` reported every generated method as a lock copy. Asserting on the
// text is deliberate: the alternative, compiling and racing the output, cannot
// be done from inside a code generator, and the textual form is exactly what
// regressed.
func TestDeepCopyIntoDoesNotAssignTheStruct(t *testing.T) {
	t.Parallel()

	got := generate(t, "Widget")

	if strings.Contains(got, "*out = *p") {
		t.Errorf("generated DeepCopyInto assigns the struct, which copies the embedded "+
			"protoimpl.MessageState and its sync.Mutex:\n%s", got)
	}
	if !strings.Contains(got, "proto.Reset(out)") {
		t.Errorf("generated DeepCopyInto does not reset the destination:\n%s", got)
	}
	if !strings.Contains(got, "proto.Merge(out, in)") {
		t.Errorf("generated DeepCopyInto does not merge through the reflection API:\n%s", got)
	}
}

// TestGeneratedShape pins the parts consumers depend on. The three method names
// and the file suffix are referenced from outside this package -- onex's
// generated zz_generated.* files sit beside this output -- so renaming any of
// them is a breaking change rather than an internal refactor.
func TestGeneratedShape(t *testing.T) {
	t.Parallel()

	got := generate(t, "Widget")

	for _, want := range []string{
		"func (in *Widget) DeepCopyInto(out *Widget) {",
		"func (in *Widget) DeepCopy() *Widget {",
		"func (in *Widget) DeepCopyInterface() interface{} {",
		"package v1",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("generated output is missing %q:\n%s", want, got)
		}
	}
}

// TestNoMessagesNoFile covers the early return: a file that declares no
// messages produces no output at all, rather than an empty package.
func TestNoMessagesNoFile(t *testing.T) {
	t.Parallel()

	if got := generate(t); got != "" {
		t.Errorf("expected no generated file for a message-less input, got:\n%s", got)
	}
}
