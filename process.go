package main

import (
	"fmt"
	"strconv"
	"strings"
)

func processRequest(cr *clientRequest) string {
	res := make([]string, 0, cr.Limit)
	for i := 1; i <= cr.Limit; i++ {
		switch {
		case (i%cr.Int1 == 0 && i%cr.Int2 == 0):
			res = append(res, fmt.Sprintf("%s%s", cr.String1, cr.String2))
		case i%cr.Int1 == 0:
			res = append(res, cr.String1)
		case i%cr.Int2 == 0:
			res = append(res, cr.String2)
		default:
			res = append(res, strconv.Itoa(i))
		}
	}
	return strings.Join(res, ",")
}

func processRequestBis(cr *clientRequest) string {
	res := make([]string, 0, cr.Limit)
	for i := 1; i <= cr.Limit; i++ {
		switch {
		case (i%cr.Int1 == 0 && i%cr.Int2 == 0):
			res = append(res, fmt.Sprintf("%s%s", cr.String1, cr.String2))
			res = append(res, cr.String1)
			res = append(res, cr.String2)
		case i%cr.Int1 == 0:
			res = append(res, cr.String1)
		case i%cr.Int2 == 0:
			res = append(res, cr.String2)
		default:
			res = append(res, strconv.Itoa(i))
		}
	}
	return strings.Join(res, ",")
}
