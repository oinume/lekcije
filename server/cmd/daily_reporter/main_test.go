package main

import (
	"io"
	"os"
	"testing"
)

func Test_dailyReporterMain_run(t *testing.T) {
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
			args: []string{"daily_reporter", "-target-date=2019-01-01"},
			fields: fields{
				outStream: os.Stdout,
				errStream: os.Stderr,
			},
			wantErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &dailyReporterMain{
				outStream: test.fields.outStream,
				errStream: test.fields.errStream,
			}
			if err := m.run(test.args); (err != nil) != test.wantErr {
				t.Errorf("dailyReporterMain.run() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
