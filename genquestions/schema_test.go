package genquestions

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.yaml.in/yaml/v4"
)

const testData = "../tests/test-data"

func TestGenerateGolden(t *testing.T) {
	out := filepath.Join(t.TempDir(), "questions.yml")
	err := Generate(Options{
		ValuesFilename:    filepath.Join(testData, "values.yaml"),
		SchemaFilename:    filepath.Join(testData, "values.schema.json"),
		QuestionsFilename: out,
		Header:            DefaultHeader,
	})
	if err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(out)
	want, err := os.ReadFile(filepath.Join(testData, "questions.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Errorf("generated file differs from %s/questions.yml", testData)
	}
}

func TestBuild(t *testing.T) {
	values := func(s string) *Node { return mustValues(t, s) }
	tests := []struct {
		name, schema, values, want string
	}{
		{
			name:   "missing label",
			schema: `{"properties":{"a":{"x-rancher-question":{"order":1,"group":"g"}}}}`,
			values: "a: 1",
			want:   "a: missing label",
		},
		{
			name:   "missing group",
			schema: `{"properties":{"a":{"x-rancher-question":{"order":1,"label":"A"}}}}`,
			values: "a: 1",
			want:   "a: missing group",
		},
		{
			name:   "unknown parent",
			schema: `{"properties":{"a":{"x-rancher-question":{"order":1,"label":"A","subquestion_of":"b"}}}}`,
			values: "a: 1",
			want:   "a: unknown subquestion_of b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := build([]byte(tt.schema), values(tt.values), "")
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("got error %v, want %q", err, tt.want)
			}
		})
	}
}

func TestBuildHeaderAndTypes(t *testing.T) {
	schema := `{"properties":{
		"n":{"x-rancher-question":{"order":2,"label":"N","group":"g"}},
		"e":{"enum":["a",""],"x-rancher-question":{"order":1,"label":"E","group":"g"}}}}`
	got, err := build([]byte(schema), mustValues(t, "n: 3\ne: a"), "first\nsecond")
	if err != nil {
		t.Fatal(err)
	}
	want := `# first
# second
questions:
  - variable: e
    label: E
    type: enum
    default: a
    options:
      a: a
      "": None
    group: g
  - variable: n
    label: N
    type: int
    default: 3
    group: g
`
	if string(got) != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

type Node = yaml.Node

func mustValues(t *testing.T, s string) *Node {
	t.Helper()
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(s), &doc); err != nil {
		t.Fatal(err)
	}
	return doc.Content[0]
}
