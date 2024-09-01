package emlparser

import (
    "testing"
)

func TestEMLParserWithDNS(t *testing.T) {
    // Create a sample EML content with DNS-related headers
    emlContent := `
    From: John Doe <johndoe@example.com>
    To: Jane Smith <janesmith@example.com>
    Subject: Test Email with DNS Records

    MX: 10 mail.example.com
    TXT: "v=spf1 mx -all"

    This is a test email.
    `

    parser := NewEMLParser(emlContent)

    // Parse the EML message with DNS extraction enabled
    options := &emlparser.ParseOptions{
        HeaderNames:        []string{"From", "To", "Subject"},
        ExtractDNSRecords: true,
    }
    result, err := parser.Parse(options)
    if err != nil {
        t.Fatalf("error parsing EML: %v", err)
    }

    // Assert the parsed data
    expected := map[string]interface{}{
        "From":    "John Doe <johndoe@example.com>",
        "To":      "Jane Smith <janesmith@example.com>",
        "Subject": "Test Email with DNS Records",
        "body":    "This is a test email.",
        "dns_records": map[string]interface{}{
            "mx": []map[string]interface{}{
                {
                    "priority": 10,
                    "host":     "mail.example.com",
                },
            },
            "txt": []string{
                "v=spf1 mx -all",
            },
        },
    }
    if !equalMaps(result, expected) {
        t.Errorf("expected result: %v, got: %v", expected, result)
    }
}

// ... (Rest of the test cases)