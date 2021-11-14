package main

import (
	"testing"
)

func TestMakeQueryToFilename(t *testing.T) {
	query := "language:c# location:italy"

	want := "language-csharp--location-italy"

	actual := MakeQueryToFilename(query)
	if want != actual {
		t.Fatalf(`MakeQueryToFilename(query) = %q, want match for %#q`, actual, want)
	}
}
