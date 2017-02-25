package email

import (
	"fmt"
	"text/template"
)

type Recipient struct {
	Address string
	Name string
}

func NewRecipient(address, name string) *Recipient {
	return &Recipient{
		Address: address,
		Name: name,
	}
}

type Template struct {
	template *template.Template
	value    string
	from     *Recipient
	to       *Recipient
}

func NewTemplate(name string, value string) *Template {
	t := &Template{
		template: template.New(name),
		value: value,
	}
	//t.template.Parse()
	return t
}

func (t *Template) Parse() error {
	return nil
}

type Email struct {
	From    *Recipient
	To      *Recipient
	Subject string
	Body    string
}

func NewEmail() *Email {
	return &Email{}
}

func NewEmailFromTemplate(t *Template /* TODO: pass values */) (*Email, error) {
	if err := t.Parse(); err != nil {
		return fmt.Errorf("Parse error: %v", err)
	}

	e := &Email{} // TODO: Set fields
	return e
}
