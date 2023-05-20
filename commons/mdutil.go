package commons

import (
	"bytes"
	"html/template"
	"strings"
)

type MdUtil struct {
}

func NewMdUtil() *MdUtil {
	return &MdUtil{}
}

type ImdUtil interface {
	RenderedToHTML(html string, payload map[string]interface{}) (string, error)
}

func (mdu *MdUtil) RenderedToHTML(html string, payload map[string]interface{}) (string, error) {
	tmpl, err := template.New("status").Parse(html)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, payload)
	if err != nil {
		return "", err
	}

	renderedHTML := buf.String()

	chars := []string{"\n", "\t"}
	for _, char := range chars {
		renderedHTML = strings.Replace(renderedHTML, char, "", -1)
	}

	return renderedHTML, nil
}
