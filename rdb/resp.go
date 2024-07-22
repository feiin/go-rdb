package rdb

import (
	"bufio"
	"errors"
	"fmt"
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

type RedisResponse struct {
	conn *Conn
}

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

func (res *RedisResponse) WriteOK() error {
	_, err := res.conn.bw.WriteString("+OK\r\n")
	if err != nil {
		return err
	}
	return res.conn.bw.Flush()
}

func (res *RedisResponse) WriteError(err error) error {
	_, rErr := res.conn.bw.WriteString(fmt.Sprintf("%c%s\r\n", Errors, err.Error()))
	if rErr != nil {
		return rErr
	}
	return res.conn.bw.Flush()
}

func (res *RedisResponse) WriteString(str string) error {
	_, err := res.conn.bw.WriteString(fmt.Sprintf("%c%s\r\n", Strings, str))
	if err != nil {
		return err
	}
	return res.conn.bw.Flush()
}

func (res *RedisResponse) WriteIntegers(ret int) error {
	_, err := res.conn.bw.WriteString(fmt.Sprintf("%c%d\r\n", Integers, ret))
	if err != nil {
		return err
	}
	return res.conn.bw.Flush()
}

func (res *RedisResponse) WriteBulkStrings(val interface{}) error {

	strVal, ok := val.(string)
	if !ok {
		_, err := res.conn.bw.WriteString(fmt.Sprintf("%c-1\r\n", BulkStrings))
		if err != nil {
			return err
		}
	} else {
		_, err := res.conn.bw.WriteString(fmt.Sprintf("%c%d\r\n%s\r\n", BulkStrings, len(strVal), strVal))
		if err != nil {
			return err
		}
	}
	return res.conn.bw.Flush()
}
