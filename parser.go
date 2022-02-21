package kensaku

import "strings"

func Parse(str string) {
	ops := make([]Operator, 0)
	flat := strings.Builder{}

	for i := 0; i < len(str); i++ {
		c := str[i]

		if c == '(' {
			op, ok, newi := readOperator(str, i+1)
			if ok {
				ops = append(ops, op)
				i = newi
			}
		} else if c == '"' {
			quot, ok, newi := readQuoted(str, i+1)
			if ok {
				ops = append(ops, &TextOperator{
					text:  quot,
					exact: true,
				})
				i = newi
			}
		} else {
			flat.WriteByte(c)
		}
	}

	if flat.Len() > 0 {
		ops = append(ops, &TextOperator{
			text: flat.String(),
		})
	}
}

func readOperator(str string, start int) (Operator, bool, int) {
	progress := strings.Builder{}
	colonidx := -1

	var field string

	for i := start; i < len(str); i++ {
		c := str[i]

		if field == "" && c == ':' {
			field = progress.String()
			colonidx = i
			progress.Reset()
		} else if c == ')' {
			return &TextOperator{
				baseOperator: baseOperator{field: field},
				text:         progress.String(),
			}, true, i
		} else {
			if c == '"' && i == colonidx+1 {
				quot, ok, newi := readQuoted(str, i+1)
				if ok {
					return &TextOperator{
						baseOperator: baseOperator{field: field},
						text:         progress.String(),
					}, true, i
				}
			}
			progress.WriteByte(c)
		}
	}

	return nil, false, start
}

func readQuoted(str string, start int) (string, bool, int) {
	progress := strings.Builder{}

	for i := start; i < len(str); i++ {
		c := str[i]

		if c == '\\' && i < len(str)-1 && str[i+1] == '"' {
			progress.WriteByte('"')
			i++
			continue
		} else if c == '"' {
			return progress.String(), true, i
		}

		progress.WriteByte(c)
	}

	return "", false, start
}
