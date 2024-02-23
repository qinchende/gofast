package iox

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type (
	textReadOptions struct {
		keepSpace     bool
		withoutBlanks bool
		omitPrefix    string
	}

	TextReadOption func(*textReadOptions)
)

// The first returned reader needs to be read first, because the content
// read from it will be written to the underlying buffer of the second reader.
func DupReadCloser(reader io.ReadCloser) (io.ReadCloser, io.ReadCloser) {
	var buf bytes.Buffer
	tee := io.TeeReader(reader, &buf)
	return ioutil.NopCloser(tee), ioutil.NopCloser(&buf)
}

func KeepSpace() TextReadOption {
	return func(o *textReadOptions) {
		o.keepSpace = true
	}
}

// ReadBytes reads exactly the bytes with the length of len(buf)
func ReadBytes(reader io.Reader, buf []byte) error {
	var got int

	for got < len(buf) {
		n, err := reader.Read(buf[got:])
		if err != nil {
			return err
		}

		got += n
	}

	return nil
}

func ReadText(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

func ReadTextLines(filename string, opts ...TextReadOption) ([]string, error) {
	var readOpts textReadOptions
	for _, opt := range opts {
		opt(&readOpts)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !readOpts.keepSpace {
			line = strings.TrimSpace(line)
		}
		if readOpts.withoutBlanks && len(line) == 0 {
			continue
		}
		if len(readOpts.omitPrefix) > 0 && strings.HasPrefix(line, readOpts.omitPrefix) {
			continue
		}

		lines = append(lines, line)
	}

	return lines, scanner.Err()
}

func WithoutBlank() TextReadOption {
	return func(o *textReadOptions) {
		o.withoutBlanks = true
	}
}

func OmitWithPrefix(prefix string) TextReadOption {
	return func(o *textReadOptions) {
		o.omitPrefix = prefix
	}
}
