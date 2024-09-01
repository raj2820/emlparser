package main

import (
    "fmt"

    "github.com/raj2820/emlparser"
)

func main() {
    // Sample EML content
    emlContent := `
    From: John Doe <johndoe@example.com>
    To: Jane Smith <janesmith@example.com>
    Subject: Test Email

    This is a test email.
    `

    // Create a new EML parser
    parser := emlparser.NewEMLParser(emlContent)

    // Set parsing options (optional)
    options := &emlparser.ParseOptions{
        HeaderNames:        []string{"From", "To", "Subject"},
        ExtractDNSRecords: true,
    }

    // Parse the EML message
    result, err := parser.Parse(options)
    if err != nil {
        fmt.Println("Error parsing EML:", err)
        return
    }

    // Access the parsed data
    fmt.Println("Parsed EML data:")
    fmt.Println("  From:", result["From"])
    fmt.Println("  To:", result["To"])
    fmt.Println("  Subject:", result["Subject"])
    fmt.Println("  Body:", result["body"])

    // Access DNS records (if extracted)
    if dnsRecords, ok := result["dns_records"].(map[string]interface{}); ok {
        fmt.Println("  DNS Records:")
        for key, value := range dnsRecords {
            fmt.Printf("    %s: %v\n", key, value)
        }
    }
}