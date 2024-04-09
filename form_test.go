// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-09, by liasica

package edocseal

import (
	"os"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func TestForm(t *testing.T) {
	var form Form
	b, _ := os.ReadFile("./runtime/input-s/x.json")
	err := jsoniter.Unmarshal(b, &form)
	require.NoError(t, err)

	for _, field := range form.Acroform.Fields {
		m := form.Objects[1]["obj:"+field.Annotation.Object]
		mb, _ := jsoniter.Marshal(m)
		var data FormFieldObject
		_ = jsoniter.Unmarshal(mb, &data)
		t.Logf("%s (%s): [%.2f, %.2f, %.2f, %.2f] T: %s",
			field.Fullname,
			field.Alternativename,
			data.Value.Rect[0],
			data.Value.Rect[1],
			data.Value.Rect[2],
			data.Value.Rect[3],
			data.Value.FT,
		)
	}
}
