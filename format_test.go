// Copyright (C) edocseal. 2025-present.
//
// Created at 2025-03-25, by liasica

package edocseal

import (
	"fmt"
	"testing"
)

func TestFormat(t *testing.T) {
	strs := [][]any{
		{
			"amount", "首期应付", "/Tx", 347.31, 757.27, 394.47, 771.39,
		},
		{
			"stagesTransfer", "企业签署日期", "/Tx", 439.47, 757.27, 476.31, 771.69,
		},
	}
	var width [7]int
	for _, str := range strs {
		for n, s := range str {
			i := len(fmt.Sprintf("%v", s)) * 2
			if i > width[n] {
				width[n] = i
			}
		}
	}
	t.Logf("width: %v", width)
	for _, str := range strs {
		// t.Logf("%-40s%-40s%-10s%-10s%-10s%-10s", str[0], str[1], str[2], str[3], str[4], str[5])
		// fmt.Printf("%20.20s %40.40s|||%10.10s\t\t%6.2f, %6.2f, %6.2f, %6.2f\n",
		// 	str[0], str[1], str[2], str[3], str[4], str[5], str[6],
		// )
		fmt.Printf("%s|%s|%s|%s|%s|%s|%s\n",
			middlePrintf(str[0].(string), width[0]),
			middlePrintf(str[1].(string), width[1]),
			middlePrintf(str[2].(string), width[2]),
			middlePrintf(fmt.Sprintf("%v", str[3]), width[3]),
			middlePrintf(fmt.Sprintf("%v", str[4]), width[4]),
			middlePrintf(fmt.Sprintf("%v", str[5]), width[5]),
			middlePrintf(fmt.Sprintf("%v", str[6]), width[6]),
		)
	}
}

func middlePrintf(s string, w int) string {
	// return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(s))/2, s))
	return fmt.Sprintf("%*s", -w, fmt.Sprintf("%*s", (w+len(s))/2, s))
	// return fmt.Sprintf(fmt.Sprintf("%%-%ds", w/2), fmt.Sprintf(fmt.Sprintf("%%%ds", w/2), s))
}
