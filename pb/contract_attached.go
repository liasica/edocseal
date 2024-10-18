// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-08, by liasica

package pb

func (c *ContractFormField_Checkbox) String() string {
	if c.Checkbox {
		return "On"
	}
	return "Off"
}
