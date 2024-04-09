// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-09, by liasica

package edocseal

type Form struct {
	Acroform struct {
		Fields []struct {
			Alternativename string `json:"alternativename,omitempty"`
			Annotation      struct {
				Annotationflags int    `json:"annotationflags,omitempty"`
				Appearancestate string `json:"appearancestate,omitempty"`
				Object          string `json:"object,omitempty"`
			} `json:"annotation,omitempty"`
			Choices       []interface{} `json:"choices,omitempty"`
			Defaultvalue  interface{}   `json:"defaultvalue,omitempty"`
			Fieldflags    int           `json:"fieldflags,omitempty"`
			Fieldtype     string        `json:"fieldtype,omitempty"`
			Fullname      string        `json:"fullname,omitempty"`
			Ischeckbox    bool          `json:"ischeckbox,omitempty"`
			Ischoice      bool          `json:"ischoice,omitempty"`
			Isradiobutton bool          `json:"isradiobutton,omitempty"`
			Istext        bool          `json:"istext,omitempty"`
			Mappingname   string        `json:"mappingname,omitempty"`
			Object        string        `json:"object,omitempty"`
			Pageposfrom1  int           `json:"pageposfrom1,omitempty"`
			Parent        interface{}   `json:"parent,omitempty"`
			Partialname   string        `json:"partialname,omitempty"`
			Quadding      int           `json:"quadding,omitempty"`
			Value         interface{}   `json:"value,omitempty"`
		} `json:"fields,omitempty"`
		Hasacroform     bool `json:"hasacroform,omitempty"`
		Needappearances bool `json:"needappearances,omitempty"`
	} `json:"acroform,omitempty"`
	Objects []map[string]any `json:"qpdf,omitempty"`
}

type FormFieldObjects []map[string]FormFieldObject

type FormFieldObject struct {
	Value struct {
		FT   string    `json:"/FT,omitempty"`
		P    string    `json:"/P,omitempty"`
		Rect []float64 `json:"/Rect,omitempty"`
	} `json:"value,omitempty"`
}
