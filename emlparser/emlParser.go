package emlparser

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
)

type EMLParser struct {
	content string
}

func NewEMLParser(content string) *EMLParser {
	return &EMLParser{content: content}
}

func NewEMLParserFromURL(url string) (*EMLParser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching EML content from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Type") != "message/rfc822" {
		return nil, fmt.Errorf("invalid content type: expected message/rfc822, got %s", resp.Header.Get("Content-Type"))
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading EML content from URL: %w", err)
	}

	return NewEMLParser(string(content)), nil
}

func (p *EMLParser) Parse(options *ParseOptions) (map[string]interface{}, error) {
	msg, err := mail.ReadMessage(strings.NewReader(p.content))
	if err != nil {
		return nil, fmt.Errorf("error parsing EML message: %w", err)
	}

	result := make(map[string]interface{})

	for k, v := range msg.Header {
		result[k] = strings.Join(v, " ")
	}

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("error parsing media type: %w", err)
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(msg.Body, params["boundary"])
		var attachments []map[string]interface{}

		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("error reading MIME part: %w", err)
			}

			attachment := map[string]interface{}{
				"filename":          part.FileName(),
				"content-type":      part.Header.Get("Content-Type"),
				"content-disposition": part.Header.Get("Content-Disposition"),
			}

			content, err := ioutil.ReadAll(part)
			if err != nil {
				return nil, fmt.Errorf("error reading attachment content: %w", err)
			}
			attachment["content"] = string(content)

			attachments = append(attachments, attachment)
		}
		result["attachments"] = attachments
	} else {
		body, err := ioutil.ReadAll(msg.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading EML body: %w", err)
		}
		result["body"] = string(body)
	}

	if options.ExtractDNSRecords {
		result["dns_records"] = p.extractDNSRecords(result)
	}

	return result, nil
}

func (p *EMLParser) extractDNSRecords(result map[string]interface{}) map[string]interface{} {
	dnsRecords := make(map[string]interface{})

	if mxRecords, ok := result["MX"].([]string); ok {
		dnsRecords["mx"] = p.extractMXRecords(mxRecords)
	}

	if txtRecords, ok := result["TXT"].([]string); ok {
		dnsRecords["txt"] = txtRecords
	}

	return dnsRecords
}

func (p *EMLParser) extractMXRecords(mxRecords []string) []map[string]interface{} {
	var mxRecordsResult []map[string]interface{}
	for _, mxRecord := range mxRecords {
		parts := strings.Split(mxRecord, " ")
		if len(parts) == 2 {
			priority, err := strconv.Atoi(parts[0])
			if err != nil {
				continue
			}
			mxRecordsResult = append(mxRecordsResult, map[string]interface{}{
				"priority": priority,
				"host":     parts[1],
			})
		}
	}
	return mxRecordsResult
}

type ParseOptions struct {
	HeaderNames       []string
	ExtractDNSRecords bool
}
