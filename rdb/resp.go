package rdb

import (
	"bufio"
	"errors"
	"strconv"
)

// https://redis.io/docs/latest/develop/reference/protocol-spec/
var (
	// RESP2
	Strings = '+'
	// RESP2
	Errors = '-'
	// RESP2
	Integers = ':'
	// RESP2
	BulkStrings = '$'
	// RESP2
	Arrays = '*'

	// RESP3
	Nulls = '_'
	// RESP3
	Booleans = '#'
	// RESP3
	Doubles = ','
	// RESP3
	BigNumbers = '('
	// RESP3
	BulkErrors = '!'
	// RESP3
	VerbatimStrings = '='
	// RESP3
	Maps = '%'
	// RESP3
	Sets = '~'
	// RESP3
	Pushes = '>'

	errorUnknownProtocol = errors.New("unknown protocol")
	errorUnknownCommand  = errors.New("unknown command")
)

// *<参数数量>CRLF
// $<参数1的字节长度>CRLF
// <参数1的数据>CRLF
// ...
// $<参数N的字节长度>CRLF
// <参数N的数据>CRLF

// 解析redis协议
func parseRESP(rd *bufio.Reader) ([]string, error) {
	line, err := rd.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if len(line) < 3 {
		return nil, errorUnknownCommand
	}

	switch line[0] {
	default:
		return nil, errorUnknownProtocol
	case '*':
		cnt, err := strconv.Atoi(line[1 : len(line)-2])
		if err != nil {
			return nil, err
		}

		payload := make([]string, 0, cnt)
		for ; cnt > 0; cnt-- {
			line, err := parseRESPLine(rd)
			if err != nil {
				return nil, err
			}
			payload = append(payload, line)
		}
		return payload, nil

	}
}

func parseRESPLine(rd *bufio.Reader) (string, error) {
	line, err := rd.ReadString('\n')
	if err != nil {
		return "", err
	}
	if len(line) < 3 {
		return "", errorUnknownCommand
	}

	switch line[0] {
	default:
		return "", errorUnknownProtocol
	// Strings, Errors, Integers
	case '+', '-', ':':
		return line[1 : len(line)-2], nil
	case '$':
		// $<length>\r\n<data>\r\n
		length, err := strconv.Atoi(line[1 : len(line)-2])
		if err != nil {
			return "", err
		}

		if length == 0 {
			return "", nil
		}

		buf := make([]byte, length+2)

		n, err := rd.Read(buf)
		if err != nil {
			return "", err
		}

		if n != length+2 {
			return "", errorUnknownCommand
		}

		return string(buf[:length]), nil

	}
}
