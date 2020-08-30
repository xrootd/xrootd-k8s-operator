package template

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/pkg/errors"
)

func createTempTextFile(contents string) (f *os.File, err error) {
	f, err = ioutil.TempFile("", "template")
	if err == nil {
		_, err = f.WriteString(contents)
	}
	return
}

func cleanupTempTextFile(file *os.File, t *testing.T) {
	if err := os.Remove(file.Name()); err != nil {
		t.Error(errors.WithMessage(err, "Deleting temporary file failed!"))
	}
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
	defer cleanupTempTextFile(file, t)
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
	defer cleanupTempTextFile(file, t)
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
	defer cleanupTempTextFile(file, t)
	result, err := ApplyTemplate(file.Name(), data{Count: 2})
	if err != nil {
		t.Errorf("error in applying template: %v", err)
	}
	t.Log(result)
	if len(result) == 0 {
		t.Error("empty result")
	}
}

func TestApplyTemplateWithInvalidPath(t *testing.T) {
	_, err := ApplyTemplate("wrong/path", nil)
	if err == nil {
		t.Errorf("should face error since path is wrong")
	}
}

func TestApplyTemplateNilData(t *testing.T) {
	type data struct {
		name string
	}
	const contents = `
		hello {{ .Name }}
	`
	file, err := createTempTextFile(contents)
	if err != nil {
		t.Errorf("error in writing contents: %v", err)
		return
	}
	defer cleanupTempTextFile(file, t)
	_, err = ApplyTemplate(file.Name(), data{name: "string"})
	if err == nil {
		t.Errorf("should face error executing with invalid template content")
	}
}
