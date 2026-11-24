package main

import "testing"

func TestParseYAMLDocument(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{
			name:   "block scalar with comment",
			source: "asyncapi: 3.0.0\ninfo:\n  title: Block Title # comment\n",
			want:   "Block Title",
		},
		{
			name:   "quoted scalar",
			source: "asyncapi: '3.0.0'\ninfo:\n  title: \"Quoted Title\"\n",
			want:   "Quoted Title",
		},
		{
			name:   "flow mapping",
			source: "asyncapi: 3.0.0\ninfo: {title: Flow Title, version: 1.0.0}\n",
			want:   "Flow Title",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			doc := parseYAMLDocument([]byte(test.source))
			if doc.asyncAPI == "" {
				t.Fatal("asyncapi version not parsed")
			}
			if doc.title != test.want {
				t.Errorf("title = %q, want %q", doc.title, test.want)
			}
		})
	}
}
