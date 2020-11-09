package http

import (
	"bytes"
	"net/url"
	"sort"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : http.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-07 18:45
 * description :
 * history :
 */

func ParseUrlValues(data map[string]string) url.Values {
	if data == nil {
		return url.Values{}
	}
	values := url.Values{}
	for k, v := range data {
		values[k] = []string{v}
	}
	return values
}

func ParseQuery(query string) (map[string]string, error) {
	values, err := url.ParseQuery(query)
	if err != nil {
		return map[string]string{}, err
	}
	mp := make(map[string]string, len(values))
	for i, v := range values {
		mp[i] = v[0]
	}
	return mp, nil
}

// convert to query and this parameters has be sorted.
func SortedQuery(r url.Values) string {
	var keys []string
	for k := range r {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	i := 0
	buf := bytes.NewBuffer(nil)
	for _, k := range keys {
		//if k == "sign" || k == "sign_type" {
		//	continue
		//}
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(r[k][0])
		i++
	}
	return buf.String()
}
