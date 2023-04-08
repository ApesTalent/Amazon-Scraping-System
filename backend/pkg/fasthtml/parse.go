package fasthtml

import (
	"strings"
)

func findNodeStartEnd(data string, pos int) (start int, end int) {
	start, end = -1, -1

	for i := pos; i > 0; i-- {
		if data[i] == '<' {
			start = i
			break
		}
	}

	var deep int
	for i := pos; i < len(data)-1; i++ {
		if (data[i] == '<' && data[i+1] == '/') || (data[i] == '/' && data[i+1] == '>') {
			deep--
			if deep == -1 {
				for ie := i; ie < len(data)-1; ie++ {
					if data[ie] == '>' {
						end = ie + 1
						break
					}
				}
				break
			}
			continue
		}

		if data[i] == '<' && data[i+1] != '/' {
			if data[i+1] == 'b' && data[i+2] == 'r' && data[i+3] == '>' {
				continue
			} else {
				deep++
				continue
			}

		}
	}

	return start, end
}

// Gets the first node with the passed attribute
func GetTagWithParams(data, name, val string) string {
	idx := strings.Index(data, name+"=\""+val+"\"")
	if idx == -1 {
		return ""
	}
	start, end := findNodeStartEnd(data, idx)

	if end < 0 || start < 0 {
		return ""
	}

	return data[start:end]
}

func GetInner(data string) string {
	var start, end = -1, -1
	for i := 0; i < len(data)-1; i++ {
		if data[i] == '>' {
			start = i + 1
			break
		}
	}

	for i := len(data) - 1; i > 1; i-- {
		if data[i] == '<' {
			end = i
			break
		}
	}

	return data[start:end]
}

func FindNodes(data, part string) []string {
	res := make([]string, 0, 16)

	i := 0
	for i < len(data) {
		idx := strings.Index(data[i:], part)
		if idx == -1 {
			break
		}

		// apply offset
		idx += i

		start, end := findNodeStartEnd(data, idx)
		if end < 0 || start < 0 {
			break
		}

		res = append(res, data[start:end])
		i = idx + len(part)
	}
	return res
}

func GetAttr(data, attr string) string {
	res := strings.Split(data, " "+attr+"=\"")
	if len(res) < 2 {
		return ""
	}

	qIdx := strings.Index(res[1], "\"")
	if qIdx == -1 {
		return ""
	}

	return res[1][:qIdx]
}
