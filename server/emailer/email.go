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
	from     *mail.Address
	tos      []*mail.Address
	subject  string
	mimeType string
	body     io.Reader
	inBody   bool
}

func NewTemplate(name string, value string) *Template {
	t := &Template{
		template: template.New(name),
		value:    strings.Replace(value, "\r", "", -1),
		body:     &bytes.Buffer{},
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

	for lineNo := 1; ; lineNo++ {
		line, err := b.ReadString([]byte("\n")[0]) // TODO: bufio.Scanner?
		if err != nil {
			if err == io.EOF {
				fmt.Printf("[%d] line = %q\n", lineNo, line)
				if err := t.parseLine(line, lineNo); err != nil {
					return err
				}
				break
			} else {
				return err
			}
		}

		fmt.Printf("[%d] line = %q\n", lineNo, line)
		if err := t.parseLine(line, lineNo); err != nil {
			return err
		}
	}

	return nil
}

func (t *Template) parseLine(line string, lineNo int) error {
	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		if t.inBody {
			fmt.Fprint(t.body.(*bytes.Buffer), line)
			return nil
		} else {
			return fmt.Errorf("Line:%v: Invalid email template", lineNo)
		}
	}

	name := strings.ToLower(strings.TrimSpace(line[:colonIndex]))
	value := strings.TrimSpace(line[colonIndex+1:])
	fmt.Printf("name = %q, value = %q\n", name, value)
	switch name {
	case "from":
		from, err := mail.ParseAddress(value)
		if err != nil {
			return fmt.Errorf("Line:%v: Parse 'from' failed: %v", lineNo, err)
		}
		t.from = from
	case "to":
		tos, err := mail.ParseAddressList(value)
		if err != nil {
			return fmt.Errorf("Line:%v: Parse 'to' failed: %v", lineNo, err)
		}
		t.tos = tos
	case "subject":
		t.subject = value
	case "body":
		if value != "text/plain" && value != "text/html" {
			return fmt.Errorf("Line:%v: Invalid body mime type: %v", lineNo, value)
		}
		t.mimeType = value
		t.inBody = true
	default:
		// TODO: accept as extra header
		return fmt.Errorf("Line:%v: Unknown header %q", lineNo, name)
	}
	return nil
}

type Email struct {
	From    *mail.Address
	To      *mail.Address
	Subject string
	Body    string
}

func NewEmail() *Email {
	return &Email{}
}

func NewEmailFromTemplate(t *Template /* TODO: pass values */) (*Email, error) {
	if err := t.Parse(); err != nil {
		return nil, fmt.Errorf("Parse error: %v", err)
	}

	e := &Email{} // TODO: Set fields
	return e, nil
}

func NewEmailsFromTemplate(t *Template /* TODO: pass values */) ([]*Email, error) {

	return nil, nil
}
