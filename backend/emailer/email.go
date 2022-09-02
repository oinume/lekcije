package emailer

import (
	"bytes"
	"fmt"
	"io"
	"net/mail"
	"strings"
	"text/template"
)

type Template struct {
	template *template.Template
	value    string
	emails   []*Email
	current  int
	inBody   bool
}

func NewTemplate(name string, value string) *Template {
	t := &Template{
		template: template.New(name),
		value:    strings.Replace(value, "\r", "", -1),
		emails:   make([]*Email, 0, 10000),
		current:  0,
	}
	return t
}

func (t *Template) Parse() error {
	_, err := t.template.Parse(t.value)
	return err
}

func (t *Template) Execute(data interface{}) error {
	var b bytes.Buffer
	if err := t.template.Execute(&b, data); err != nil {
		return err
	}

	email := NewEmail()
	t.emails = append(t.emails, email)
	defer func() {
		t.current++
	}()

	for lineNo := 1; ; lineNo++ {
		line, err := b.ReadString([]byte("\n")[0]) // TODO: bufio.Scanner?
		if err != nil {
			if err == io.EOF {
				defer func() {
					t.inBody = false
				}()
				//fmt.Printf("[%d] line = %q\n", lineNo, line)
				if err := t.parseLine(line, lineNo, email); err != nil {
					return err
				}
				break
			} else {
				return err
			}
		}

		//fmt.Printf("[%d] line = %q\n", lineNo, line)
		if err := t.parseLine(line, lineNo, email); err != nil {
			return err
		}
	}

	return nil
}

func (t *Template) parseLine(line string, lineNo int, email *Email) error {
	if t.inBody {
		fmt.Fprint(email.Body.(*bytes.Buffer), line)
		return nil
	}

	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return fmt.Errorf("line:%v: Invalid email template", lineNo)
	}

	name := strings.ToLower(strings.TrimSpace(line[:colonIndex]))
	value := strings.TrimSpace(line[colonIndex+1:])
	//fmt.Printf("name = %q, value = %q\n", name, value)
	switch name {
	case "from":
		from, err := mail.ParseAddress(value)
		if err != nil {
			return fmt.Errorf("line:%v: Parse 'from' failed: %v", lineNo, err)
		}
		email.From = from
	case "to":
		tos, err := mail.ParseAddressList(value)
		if err != nil {
			return fmt.Errorf("line:%v: Parse 'to' failed: %v", lineNo, err)
		}
		email.Tos = tos
	case "subject":
		email.Subject = value
	case "body":
		if value != "text/plain" && value != "text/html" {
			return fmt.Errorf("line:%v: Invalid body mime type: %v", lineNo, value)
		}
		email.BodyMIMEType = value
		t.inBody = true // TODO: goroutine safe
	default:
		// TODO: accept as extra header
		return fmt.Errorf("line:%v: Unknown header %q", lineNo, name)
	}
	return nil
}

type Email struct {
	From         *mail.Address
	Tos          []*mail.Address
	Subject      string
	BodyMIMEType string
	Body         io.Reader
	bodyCache    string
	customArgs   map[string]string
}

func NewEmail() *Email {
	return &Email{
		Body:       &bytes.Buffer{},
		customArgs: make(map[string]string),
	}
}

// Create Email from Template with given data.
// Return an error if
// - Parsing template fails
func NewEmailFromTemplate(t *Template, data interface{}) (*Email, error) {
	if err := t.Parse(); err != nil {
		return nil, fmt.Errorf("Parse error: %v", err)
	}
	if err := t.Execute(data); err != nil {
		return nil, err
	}
	return t.emails[0], nil
}

func NewEmailsFromTemplate(t *Template, data []interface{}) ([]*Email, error) {
	if err := t.Parse(); err != nil {
		return nil, fmt.Errorf("Parse error: %v", err)
	}
	for _, d := range data {
		if err := t.Execute(d); err != nil {
			return nil, err
		}
	}
	return t.emails, nil
}

func (e *Email) BodyString() string {
	if e.bodyCache != "" {
		return e.bodyCache
	}
	b, err := io.ReadAll(e.Body)
	if err != nil {
		return ""
	}
	e.bodyCache = string(b)
	return e.bodyCache
}

func (e *Email) SetCustomArg(key, value string) {
	e.customArgs[key] = value
}
