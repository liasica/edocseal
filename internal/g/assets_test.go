// Copyright (C) edocseal. 2025-present.
//
// Created at 2025-06-25, by liasica

package g

import (
	"io"
	"testing"
)

func TestGetFontSong(t *testing.T) {
	fontName, reader := GetFontSong()
	if fontName != "Song" {
		t.Errorf("expected font name 'Song', got '%s'", fontName)
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read font data: %v", err)
	}

	if len(data) == 0 {
		t.Error("font data is empty")
	}
}
