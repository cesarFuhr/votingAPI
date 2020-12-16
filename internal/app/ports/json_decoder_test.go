package ports

import (
	"net/http"
	"strings"
	"testing"
)

type testStruct struct {
	property string
}

func TestDecodeJSONBody(t *testing.T) {
	t.Run("Should return err for empty body", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodPost, "whatever", nil)
		want := malformedRequest{
			status: http.StatusBadRequest,
			msg:    "Invalid: Empty body",
		}

		dst := testStruct{}
		got := decodeJSONBody(r, dst, false)

		if got.Error() != want.Error() {
			t.Errorf("want %v, got %v", want, got)
		}
	})
	t.Run("Should return err for empty body", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodPost, "whatever", strings.NewReader(""))
		want := malformedRequest{
			status: http.StatusBadRequest,
			msg:    "Invalid: Empty body",
		}

		dst := testStruct{}
		got := decodeJSONBody(r, dst, false)

		if got.Error() != want.Error() {
			t.Errorf("want %v, got %v", want, got)
		}
	})
	t.Run("Should not return an err for empty body in ignore is true", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodPost, "whatever", nil)
		var want error = nil

		dst := testStruct{}
		got := decodeJSONBody(r, dst, true)

		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	})
	t.Run("Should return an error for a invalid json body", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodPost, "whatever", strings.NewReader("{\"test\": }"))
		wantedMsgPrefix := "Request body contains invalid JSON"

		dst := testStruct{}
		got := decodeJSONBody(r, dst, false)

		assertStringPrefix(t, got.Error(), wantedMsgPrefix)
	})
}

func assertStringPrefix(t *testing.T, got, wantedPrefix string) {
	t.Helper()
	if !strings.HasPrefix(got, wantedPrefix) {
		t.Errorf("wanted %q, but got %q", wantedPrefix, got)
	}
}
