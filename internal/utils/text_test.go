package utils

import (
	"strings"
	"testing"
)

func TestEscapeHTML(t *testing.T) {
	input := "A&B <C>"
	expected := "A&amp;B &lt;C&gt;"

	if got := EscapeHTML(input); got != expected {
		t.Fatalf("unexpected escaped HTML:\nexpected: %q\nactual:   %q", expected, got)
	}
}

func TestFormatResponseWithQuestionEscapesQuestionOnly(t *testing.T) {
	result := FormatResponseWithQuestion("<hi> &", "<b>bold</b>")

	if !strings.Contains(result, "Question:\n&lt;hi&gt; &amp;\n\n") {
		t.Fatalf("escaped question not found in result: %q", result)
	}
	if !strings.Contains(result, "Answer:\n<b>bold</b>") {
		t.Fatalf("raw answer not found in result: %q", result)
	}
	if strings.Contains(result, "&lt;b&gt;") {
		t.Fatalf("answer should not be escaped: %q", result)
	}
}

func TestSplitMessageShortText(t *testing.T) {
	text := "short message"
	parts := SplitMessage(text, 100)

	if len(parts) != 1 || parts[0] != text {
		t.Fatalf("unexpected split result: %#v", parts)
	}
}

func TestSplitMessageSplitsOnNewline(t *testing.T) {
	text := "one\ntwo three"
	parts := SplitMessage(text, 7)

	expected := []string{"one", "two three"}
	if len(parts) != len(expected) {
		t.Fatalf("unexpected number of parts: %#v", parts)
	}
	for i := range expected {
		if parts[i] != expected[i] {
			t.Fatalf("part %d mismatch: expected %q got %q", i, expected[i], parts[i])
		}
	}
}

func TestSplitMessageSplitsOnSpace(t *testing.T) {
	text := "alpha beta gamma"
	parts := SplitMessage(text, 10)

	expected := []string{"alpha", "beta gamma"}
	if len(parts) != len(expected) {
		t.Fatalf("unexpected number of parts: %#v", parts)
	}
	for i := range expected {
		if parts[i] != expected[i] {
			t.Fatalf("part %d mismatch: expected %q got %q", i, expected[i], parts[i])
		}
	}
}

func TestSplitMessageRespectsMaxLength(t *testing.T) {
	text := "one two three four five six seven eight nine ten"
	maxLength := 12

	parts := SplitMessage(text, maxLength)
	if len(parts) < 2 {
		t.Fatalf("expected multiple parts, got %#v", parts)
	}
	for i, part := range parts {
		if part == "" {
			t.Fatalf("part %d is empty", i)
		}
		if len(part) > maxLength {
			t.Fatalf("part %d exceeds max length: %d", i, len(part))
		}
	}
}
