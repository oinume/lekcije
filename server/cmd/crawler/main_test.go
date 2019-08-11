package main

import (
	"io"
	"os"
	"testing"
)

func Test_crawlerMain_run(t *testing.T) {
	type fields struct {
		outStream io.Writer
		errStream io.Writer
	}

	tests := map[string]struct {
		args    []string
		fields  fields
		wantErr bool
	}{
		"normal": {
			args: []string{"crawler", "-ids=1"}, // TODO: Add fetcher/fake package
			fields: fields{
				outStream: os.Stdout,
				errStream: os.Stderr,
			},
			wantErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &crawlerMain{
				outStream: test.fields.outStream,
				errStream: test.fields.errStream,
			}
			if err := m.run(test.args); (err != nil) != test.wantErr {
				t.Errorf("crawlerMain.run() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
