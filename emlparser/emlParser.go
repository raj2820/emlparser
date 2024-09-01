package emlparser

import (
    "fmt"
    "io/ioutil"
    "mime/multipart"
    "net/mail"
    "strings"

    "github.com/miekg/dns"
)

type EMLParser struct {
    content string
}

func NewEMLParser(content string) *EMLParser {
    return &EMLParser{content: content}
}

func NewEMLParserFromURL(url string) (*EMLParser, error) {
    // ... (Similar to HTML parser)
}

func (p *EMLParser) Parse(options *ParseOptions) (map[string]interface{}, error) {
    // Parse the EML message
    msg, err := mail.ReadMessage(strings.NewReader(p.content))
    if err != nil {
        return nil, fmt.Errorf("error parsing EML message: %w", err)
    }

    // Extract relevant data from the message headers and body
    result := make(map[string]interface{})
    for _, header := range msg.Header {
        result[header[0]] = strings.Join(header[1:], "")
    }

    // Handle attachments
    if msg.Multipart {
        result["attachments"] = make([]map[string]interface{}, 0)
        for {
            part, err := msg.Multipart.NextPart()
            if err == io.EOF {
                break
            }
            if err != nil {
                return nil, fmt.Errorf("error reading EML part: %w", err)
            }

            attachment := make(map[string]interface{})
            attachment["filename"] = part.FileName()
            attachment["content-type"] = part.Header.Get("Content-Type")
            attachment["content-disposition"] = part.Header.Get("Content-Disposition")

            // Extract attachment content
            content, err := ioutil.ReadAll(part)
            if err != nil {
                return nil, fmt.Errorf("error reading EML attachment: %w", err)
            }
            attachment["content"] = string(content)

            result["attachments"] = append(result["attachments"].([]map[string]interface{}), attachment)
        }
    }

    // Extract body content
    body, err := ioutil.ReadAll(msg.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading EML body: %w", err)
    }
    result["body"] = string(body)

    // Extract DNS records
    if options.ExtractDNSRecords {
        result["dns_records"] = p.extractDNSRecords(result)
    }

    return result, nil
}

func (p *EMLParser) extractDNSRecords(result map[string]interface{}) map[string]interface{} {
    dnsRecords := make(map[string]interface{})

    // Extract MX records
    if mxRecords, ok := result["MX"].([]string); ok {
        dnsRecords["mx"] = p.extractMXRecords(mxRecords)
    }

    // Extract TXT records
    if txtRecords, ok := result["TXT"].([]string); ok {
        dnsRecords["txt"] = txtRecords
    }

    // Extract other DNS records as needed
    // ...

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
    HeaderNames []string
    ExtractDNSRecords bool
}