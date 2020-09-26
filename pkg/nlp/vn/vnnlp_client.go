package vn

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// NLPClient provdes communication functions with a running NLPServer
type NLPClient struct {
	Address    string
	Port       string
	Host       string
	Annotators string
}

// Payload is the body content for POST request
type Payload struct {
	Text string `json:"text"`
}

// NewNLPClient generates a new NLPClient instance
func NewNLPClient(host string, port string, annotators string) *NLPClient {
	var address = strings.Join([]string{`http://`, host, `:`, port}, ``)
	var c = NLPClient{
		Host:       host,
		Port:       port,
		Address:    address,
		Annotators: annotators,
	}
	return &c
}

// DefaultClient creates a NLPClient with default credentials for a local server
func DefaultClient() *NLPClient {
	var address = strings.Join([]string{`http://`, defaultHost, `:`, defaultPort}, ``)
	var c = NLPClient{
		Host:       defaultHost,
		Port:       defaultPort,
		Address:    address,
		Annotators: defaultAnnotators,
	}
	return &c
}

// Ping sends a GET msg to the server to see if it is alive
func (c *NLPClient) Ping() bool {
	var address = c.Address
	if address == "" {
		return false
	}

	resp, err := http.Get(address)
	if err != nil {
		return false
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}

// getAnnotators gets the annotators registered with the Server
func (c *NLPClient) getAnnotators() (string, error) {
	resp, err := http.Get(strings.Join([]string{c.Address, `/annotators`}, ``))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Errorf("Unable to read server response: %v", err)
	}

	// process annotators
	var a string
	a = string(bodyBytes)
	a = strings.ReplaceAll(a, `"`, ``)
	a = strings.ReplaceAll(a, `[`, ``)
	a = strings.ReplaceAll(a, `]`, ``)
	return a, nil
}

// Annotate sends a request to the server to ask for an annotation of a string and returns the response
func (c *NLPClient) annotate(textInput string, annotators string) (*ServerResponse, error) {
	if annotators == "default" {
		annotators = c.Annotators
	}

	// Construct URL
	urlStr := strings.Join([]string{c.Address, `/handle`}, ``)
	u, _ := url.Parse(urlStr)
	query, _ := url.ParseQuery(u.RawQuery)
	query.Add(`text`, textInput)
	query.Add(`props`, annotators)
	u.RawQuery = query.Encode()

	// Set content type
	contentType := `application/json`

	// Prepare payload
	payload := Payload{Text: textInput}
	body, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	// Send request
	resp, err := http.Post(u.String(), contentType, bytes.NewBuffer(body))

	if err != nil {
		return nil, &Error{Type: ErrorTypeServerError, Err: err}
	}

	defer resp.Body.Close()
	var sr ServerResponse

	switch resp.StatusCode {
	case http.StatusOK:
		err = json.NewDecoder(resp.Body).Decode(&sr)
		if err != nil {
			return nil, errors.Errorf("Unable to read server response: %v", err)
		}
		if sr.Error != "" {
			return nil, &Error{Type: ErrorTypeServerError, Msg: sr.Error}
		}
		return &sr, nil
	case http.StatusRequestURITooLong:
		return nil, &Error{Type: ErrorTypeTextTooLong, Msg: resp.Status}
	default:
		return nil, &Error{Type: ErrorTypeRequestFailed, Msg: resp.Status}
	}
}

func (c *NLPClient) customAnnotate(text string, annotators string) (*[]Sentence, error) {
	result, err := c.annotate(text, annotators)
	if err != nil {
		return nil, err
	}
	return &result.Sentences, nil
}

// Tokenize return tokens from input string, otherwise empty string
func (c *NLPClient) Tokenize(text string) (*[]Sentence, error) {
	return c.customAnnotate(text, "wseg")
}

// PosTag returns POS tags from the input string, otherwise empty string
func (c *NLPClient) PosTag(text string) (*[]Sentence, error) {
	return c.customAnnotate(text, "wseg,pos")
}

// Ner returns NER - Named Entity Recognition from the input string, otherwise empty string
func (c *NLPClient) Ner(text string) (*[]Sentence, error) {
	return c.customAnnotate(text, "wseg,pos,ner")
}

// DepParse returns parsed dependencies from input string, otherwise empty string
func (c *NLPClient) DepParse(text string) (*[]Sentence, error) {
	return c.customAnnotate(text, "wseg,pos,ner,parse")
}

// DetectLanguage returns the detected language of the input string
func (c *NLPClient) DetectLanguage(text string) (string, error) {
	result, err := c.annotate(text, "lang")
	if err != nil {
		return "", err
	}
	return result.Language, nil
}
