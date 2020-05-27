package core

import (
	"encoding/xml"
	"strings"
)

var defaultExcludeCommandPrefixs []string

type ExcludeCommandXmlRoot struct {
	XMLName         xml.Name `xml:"prefix"`
	ExcludeCommands []string `xml:"exclude"`
}

func (this *ExcludeCommandXmlRoot) FillWithXml(xstr string) error {
	err := xml.Unmarshal([]byte(xstr), this)
	if err != nil {
		return err
	}

	defaultExcludeCommandPrefixs = make([]string, 0, len(this.ExcludeCommands))
	for idx := range this.ExcludeCommands {
		this.ExcludeCommands[idx] = strings.Trim(this.ExcludeCommands[idx], " ")
		defaultExcludeCommandPrefixs = append(defaultExcludeCommandPrefixs, this.ExcludeCommands[idx]+" ")
	}
	return nil
}
