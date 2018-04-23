package main

import (
	"strings"
	"bytes"
)

type fieldList []Field

type FieldArgParseError struct {
	err string
}

func (e *FieldArgParseError) Error() string {
	return e.err
}

func (s *fieldList) String() string{
	var buffer bytes.Buffer
	for _, item := range *s {
		buffer.WriteString("Element Title: ")
		buffer.WriteString(item.Header)
		buffer.WriteString(", Text: ")
		buffer.WriteString(item.Text)
		buffer.WriteString(", Short?: ")
		if item.Short {
			buffer.WriteString("yes\n")
		} else {
			buffer.WriteString("no\n")
		}
	}
	return buffer.String()
}

func (s *fieldList) Set(val string) error {
	arr := strings.Split(val, ",")
	for _, elem := range arr {
		values := strings.Split(elem, ":")
		if len(values) != 3 {
			return &FieldArgParseError{"Split produced unexpected number of parts (delimiter ':')"}
		}
		var field Field
		field.Text = values[1]
		field.Header = values[0]
		field.Short = values[2] == "yes"
		*s = append(*s, field)
	}
	return nil
}