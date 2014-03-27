// Copyright 2014 InsyncHQ, Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.
package ses

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

type Addresses []string

type Destination struct {
	BccAddresses Addresses
	CcAddresses  Addresses
	ToAddresses  Addresses
}

type Content struct {
	Data    string
	Charset string
}

type Body struct {
	Html *Content
	Text *Content
}

type Message struct {
	Subject *Content
	Body    *Body
}

type Email struct {
	Destination      *Destination
	Message          *Message
	ReplyToAddresses Addresses
	ReturnPath       string
	Source           string
}

func (e *Email) UrlValues() url.Values {
	data := make(url.Values)
	for _, p := range toParams(e, nil) {
		data.Add(p[0], p[1])
	}
	return data
}

func toParams(o interface{}, p []string) [][]string {
	var params [][]string
	var s reflect.Value
	var t reflect.Type
	s = reflect.ValueOf(o)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	t = s.Type()
	var prefix string
	if len(p) > 0 {
		prefix = strings.Join(p, ".") + "."
	} else {
		prefix = ""
	}
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if f.Kind() == reflect.Ptr {
			f = f.Elem()
			if !f.IsValid() {
				continue
			}
		}
		if f.Kind() == reflect.String {
			if f.String() == "" {
				continue
			}
		}
		n := t.Field(i).Name
		switch f.Type() {
		case reflect.TypeOf(Addresses{}):
			for j := 0; j < f.Len(); j++ {
				params = append(params, []string{
					fmt.Sprintf("%s%s.members.%d", prefix, n, j+1),
					f.Index(j).String(),
				})
			}
		case reflect.TypeOf(Content{}):
			cs := f.FieldByName("Charset")
			if cs.String() == "" {
				cs.SetString("UTF-8")
			}
		case reflect.TypeOf(""):
			params = append(params, []string{
				fmt.Sprintf("%s%s", prefix, n),
				f.String(),
			})
		}
		if f.Kind() == reflect.Struct {
			np := append(p, n)
			params = append(params, toParams(f.Interface(), np)...)
		}
	}
	return params
}
