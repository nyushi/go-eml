package eml

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// Message represents email
type Message struct {
	RawHeaders     map[string][]string
	RawBody        string
	DecodedHeaders map[string][]string
	DecodedBody    *string
	Parts          []*Message
}

func (m *Message) transferEncoding() string {
	values, ok := m.RawHeaders["Content-Transfer-Encoding"]
	if !ok {
		return ""
	}
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

var encodings = map[string]encoding.Encoding{
	"iso-2022-jp": japanese.ISO2022JP,
	"euc-jp":      japanese.EUCJP,
	"shift_jis":   japanese.ShiftJIS,
}

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	enc, ok := encodings[strings.ToLower(charset)]
	if !ok {
		return input, nil
	}
	reader := transform.NewReader(input, enc.NewDecoder())
	return reader, nil
}

func parseMessage(body []byte, header map[string][]string) (*Message, error) {
	dhdr, err := decodeHeaders(header)
	if err != nil {
		return nil, fmt.Errorf("failed to decode headers: %s", err)
	}
	msg := &Message{
		RawHeaders: map[string][]string(header),
		RawBody:    string(body),

		DecodedHeaders: dhdr,
	}
	ctype, params, err := mime.ParseMediaType(msg.RawHeaders["Content-Type"][0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse Content-Type header: %s", err)
	}
	charset := params["charset"]
	te := msg.transferEncoding()

	if strings.HasPrefix(ctype, "text/") {
		decBody, err := decodeBody(charset, te, bytes.NewBufferString(msg.RawBody))
		if err != nil {
			return nil, fmt.Errorf("failed to decode body: %s", err)
		}
		msg.DecodedBody = &decBody
	} else if strings.HasPrefix(ctype, "multipart/") {
		msg.Parts = []*Message{}
		mr := multipart.NewReader(bytes.NewBuffer(body), params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("error at NextPart in multipart parsing: %s", err)
			}
			partBody, err := ioutil.ReadAll(p)
			if err != nil {
				return nil, fmt.Errorf("error at ReadAll in multipart parsing: %s", err)
			}
			partMsg, err := parseMessage(partBody, p.Header)
			if err != nil {
				return nil, fmt.Errorf("failed to parse multipart: %s", err)
			}
			msg.Parts = append(msg.Parts, partMsg)
		}
	}
	return msg, nil
}

// Parse read data from io.Reader and returns
func Parse(r io.Reader) (*Message, error) {
	tp := textproto.NewReader(bufio.NewReader(r))
	hdr, err := tp.ReadMIMEHeader()
	if err != nil {
		return nil, fmt.Errorf("failed to read mime header: %s", err)
	}
	b, err := ioutil.ReadAll(tp.R)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %s", err)
	}
	return parseMessage(b, hdr)
}

func decodeBody(charset, transferEncoding string, body io.Reader) (string, error) {
	var (
		r   io.Reader
		err error
	)
	r, err = charsetReader(charset, body)
	if err != nil {
		return "", fmt.Errorf("failed to get charsetReader: %s", err)
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("failed to convert charset: %s", err)
	}
	te := strings.ToLower(transferEncoding)
	if te == "base64" {
		data, err := base64.StdEncoding.DecodeString(string(b))
		if err != nil {
			return "", fmt.Errorf("failed to decode b64: %s", err)
		}
		b = data
	} else if te == "quoted-printable" {
		b, err = ioutil.ReadAll(quotedprintable.NewReader(bytes.NewBuffer(b)))
		if err != nil {
			return "", fmt.Errorf("failed to decode quoted-printable: %s", err)
		}
	}
	return string(b), nil
}

func decodeHeaders(hdr map[string][]string) (map[string][]string, error) {
	dec := &mime.WordDecoder{
		CharsetReader: charsetReader,
	}
	decoded := make(map[string][]string, len(hdr))
	for k, values := range hdr {
		decoded[k] = make([]string, len(values))
		for i, v := range values {
			d, err := dec.DecodeHeader(v)
			if err != nil {
				return nil, fmt.Errorf("failed to decode %s: %s: %s", k, v, err)
			}
			decoded[k][i] = d
		}
	}
	return decoded, nil
}
