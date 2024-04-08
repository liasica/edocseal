// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-08, by liasica

package edocseal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePdf(t *testing.T) {
	fields, err := ParsePdfFields("templates/input.pdf")
	require.NoError(t, err)
	t.Log(fields)
}
