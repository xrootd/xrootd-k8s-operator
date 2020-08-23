package template

import (
	"io/ioutil"
	"os"
	"testing"
)

func createTempTextFile(contents string) (f *os.File, err error) {
	f, err = ioutil.TempFile("", "template")
	if err == nil {
		_, err = f.WriteString(contents)
	}
	return
}

func TestApplyTemplateWithNoData(t *testing.T) {
	const contents = `
		hello world
	`
	file, err := createTempTextFile(contents)
	if err != nil {
		t.Errorf("error in writing contents: %v", err)
		return
	}
	defer os.Remove(file.Name())
	result, err := ApplyTemplate(file.Name(), nil)
	if err != nil {
		t.Errorf("error in applying template: %v", err)
	}
	t.Log(result)
	if len(result) == 0 {
		t.Error("empty result")
	}
}

func TestApplyTemplateWithData(t *testing.T) {
	type data struct {
		Name string
	}
	const contents = `
		hello {{ .Name }}
	`
	file, err := createTempTextFile(contents)
	if err != nil {
		t.Errorf("error in writing contents: %v", err)
		return
	}
	defer os.Remove(file.Name())
	result, err := ApplyTemplate(file.Name(), data{Name: "shivansh"})
	if err != nil {
		t.Errorf("error in applying template: %v", err)
	}
	t.Log(result)
	if len(result) == 0 {
		t.Error("empty result")
	}
}

func TestApplyTemplateWithTemplateFunction(t *testing.T) {
	type data struct {
		Count int
	}
	const contents = `
		array - {{ Iterate .Count }}
	`
	file, err := createTempTextFile(contents)
	if err != nil {
		t.Errorf("error in writing contents: %v", err)
		return
	}
	defer os.Remove(file.Name())
	result, err := ApplyTemplate(file.Name(), data{Count: 2})
	if err != nil {
		t.Errorf("error in applying template: %v", err)
	}
	t.Log(result)
	if len(result) == 0 {
		t.Error("empty result")
	}
}