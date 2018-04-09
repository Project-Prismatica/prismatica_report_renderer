package template_functions

import (
	"bytes"
	"errors"

	"github.com/antchfx/xmlquery"
	"github.com/flosch/pongo2"
	"github.com/sirupsen/logrus"
)

const (
	xpathFilterName = "xpath"
)

func init() {
	pongo2.RegisterFilter(xpathFilterName, xpathFilter)
}

func xpathFilter(xmlToQueryValue *pongo2.Value, xpathQueryValue *pongo2.Value)(
		out *pongo2.Value, err *pongo2.Error) {

	filterError := pongo2.Error{Sender: xpathFilterName}
	logrus.WithFields(logrus.Fields{"xmlToQueryValue": xmlToQueryValue,
		"xpathQueryValue": xpathQueryValue}).Debug(
		"xpathFilter",
	)

	if ! xmlToQueryValue.IsString() {
		logrus.WithField("input", xmlToQueryValue).
			Debug("is not string type")
		filterError.OrigError = errors.
			New("input to xpath must be a string")
		err = &filterError
		return
	}

	if ! xpathQueryValue.IsString() {
		logrus.WithField("xpathQueryValue", xpathQueryValue).
			Debug("is not string type")
		filterError.OrigError = errors.
			New("xpath query parameter must be a string")
		err = &filterError
		return
	}

	xmlToQueryBytes := bytes.NewBufferString(xmlToQueryValue.String())
	parsedXml, xmlParseError := xmlquery.Parse(xmlToQueryBytes)
	if xmlParseError != nil {
		filterError.OrigError = errors.New("could parse input XML: " +
			xmlParseError.Error())
		err = &filterError
		return
	}

	nodesFound := xmlquery.Find(parsedXml, xpathQueryValue.String())
	if nodesFound == nil {
		logrus.WithField("xpathQuery", xpathQueryValue.String()).
			Warn("returned no results")
		out = pongo2.AsValue("")
		return
	}

	logrus.WithFields(logrus.Fields{"xpathQuery": xpathQueryValue.String(),
		"resultCount": len(nodesFound)}).
		Debug("returning")

	if len(nodesFound) > 1 {
		var foundInnerText []string
		for _, element := range nodesFound {
			foundInnerText = append(foundInnerText, element.InnerText())
		}

		out = pongo2.AsValue(foundInnerText)
		return
	}

	out = pongo2.AsValue(nodesFound[0].InnerText())

	return
}
