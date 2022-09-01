package pdfobjects

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type CMapType interface {
	Convert(string) (string, error)
	HasMapping(string) bool
}

type CMap []CMapType

func (c CMap) HasMapping(in string) bool {
	for _, entry := range c {
		if entry.HasMapping(in) {
			return true
		}
	}

	return false
}

func (c CMap) Convert(in string) (string, error) {
	for _, entry := range c {
		if entry.HasMapping(in) {
			return entry.Convert(in)
		}
	}

	return in, nil
}

type CMapSingle struct {
	from string
	to string
}

func (c CMapSingle) HasMapping(in string) bool {
	return in == c.from
}

func (c CMapSingle) Convert(in string) (string, error) {
	if c.HasMapping(in) {
		return c.from, nil
	}

	return in, nil
}

func (c CMapSingle) String() string {
	return fmt.Sprintf("from: %v, to: %v", c.from, c.to)
}

type CMapRange struct {
	begin int
	end int
	initial int
}

func (c CMapRange) HasMapping(in string) bool {
	conv, err := parseHex(in)
	
	if err != nil {
		return false
	}
	
	return conv >= c.begin && conv >= c.end
}

func (c CMapRange) Convert(in string) (string, error) {
	if c.HasMapping(in) {
		conv, err := parseHex(in)
		dist := conv - c.begin
		return fmt.Sprintf("%x", c.initial + dist), err
	} else {
		return in, nil
	}
}

func (c CMapRange) String() string {
	return fmt.Sprintf("begin: %d, end: %d, initial: %d", c.begin, c.end, c.initial)
}

func cleanBfChar(line string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			line, "<", ""),
		">", "")
}

func removeBrackets(line string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			line, "[", ""),
		"]", "")
}

func parseCMapSingle(line string, cmap *CMap) {
	split := strings.SplitN(cleanBfChar(line), " ", 2)
	from, to := split[0], split[1]

	single := CMapSingle{from, to}
	*cmap = append(*cmap, single)
}

func parseHex(str string) (int, error) {
	res, err := strconv.ParseInt(str, 16, 64)	
	return int(res), err
}

func parseCMapRange(line string, cmap *CMap) {
	split := strings.SplitN(cleanBfChar(line), " ", 3)

	if strings.HasSuffix(split[2], "]") {
		from, to := split[0], removeBrackets(split[2])
		single := CMapSingle{from, to}
		*cmap = append(*cmap, single)
	} else {
		begin, _ := parseHex(split[0])
		end, _ := parseHex(split[1])
		initial, _ := parseHex(split[2])

		*cmap = append(*cmap, CMapRange{begin, end, initial})
	}
}

func parseCMap(reader *bytes.Reader, cmap *CMap) {
	var buf bytes.Buffer
	buf.ReadFrom(reader)
	lines := strings.Split(buf.String(), "\n")

	bfchar, bfrange := false, false
	
	for _, line := range lines {
		if bfchar {
			if strings.HasSuffix(line, "endbfchar") {
				bfchar = false
			} else {
				parseCMapSingle(line, cmap)
			}
		} else if bfrange {
			if strings.HasSuffix(line, "endbfrange") {
				bfrange = false
			} else {
				parseCMapRange(line, cmap)
			}
		} else {
			if strings.HasSuffix(line, "beginbfchar") {
				bfchar = true
			} else if strings.HasSuffix(line, "beginbfrange") {
				bfrange = true
			}
		}
	}
}
