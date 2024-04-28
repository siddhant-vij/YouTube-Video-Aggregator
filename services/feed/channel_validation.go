package feed

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	xpp "github.com/mmcdole/goxpp"
	"golang.org/x/net/html/charset"
)

type FeedType int

const (
	FeedTypeUnknown FeedType = iota
	FeedTypeAtom
	// Will handle other feed types later
)

func ValidateFeedURL(url string) bool {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return detectFeedType(resp.Body) == FeedTypeAtom
}

func detectFeedType(feed io.Reader) FeedType {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(feed)

	var firstChar byte
loop:
	for {
		ch, err := buffer.ReadByte()
		if err != nil {
			return FeedTypeUnknown
		}
		switch ch {
		case ' ', '\r', '\n', '\t':
		case 0xFE, 0xFF, 0x00, 0xEF, 0xBB, 0xBF:
		default:
			firstChar = ch
			buffer.UnreadByte()
			break loop
		}
	}

	if firstChar == '<' {
		p := xpp.NewXMLPullParser(bytes.NewReader(buffer.Bytes()), false, newReaderLabel)

		_, err := findRoot(p)
		if err != nil {
			return FeedTypeUnknown
		}

		name := strings.ToLower(p.Name)
		switch name {
		case "feed":
			return FeedTypeAtom
		default:
			return FeedTypeUnknown
		}
	}
	return FeedTypeUnknown
}

func newReaderLabel(label string, input io.Reader) (io.Reader, error) {
	conv, err := charset.NewReaderLabel(label, input)
	if err != nil {
		return nil, err
	}
	return conv, nil
}

func findRoot(p *xpp.XMLPullParser) (event xpp.XMLEventType, err error) {
	for {
		event, err = p.Next()
		if err != nil {
			return event, err
		}
		if event == xpp.StartTag {
			break
		}

		if event == xpp.EndDocument {
			return event, fmt.Errorf("failed to find root node before document end")
		}
	}
	return
}
