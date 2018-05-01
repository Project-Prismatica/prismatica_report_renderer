package cmd

import (
	"errors"
	"fmt"
	"github.com/Project-Prismatica/prismatica_report_renderer/templating_engine"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	//"github.com/Project-Prismatica/prismatica_report_renderer"
)

type LocalRuntimeParameters struct {
	inputTemplate string
	outputFile *os.File
	inputTemplateFileName, outputFileName string
	renderContextVariables map[string]string
}

const (
	tupleSeparator = "="
)

var (

localRenderCommand = &cobra.Command{
	Use: "render",
	Short: "render a template",
	Args: validateLocalRenderArguments,
	Run: doLocalRender,
}

localRenderRuntimeParameters = LocalRuntimeParameters{}

tupleSplitRegex, _ = regexp.Compile(tupleSeparator)
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	RootCmd.AddCommand(localRenderCommand)

	localRenderCommand.PersistentFlags().
		StringVarP(&localRenderRuntimeParameters.inputTemplateFileName,
			"source", "s", "input.tpl",
			"the input template file name")

	localRenderCommand.PersistentFlags().
		StringVarP(&localRenderRuntimeParameters.outputFileName,
		"output", "o", "-",
		"the filename to which to write the rendered file, - for stdout")
}

func parseTupelsToMap(args []string) (parsed map[string]string, err error) {
	parsed = make(map[string]string)

	for _, element := range args {
		splitTuple := tupleSplitRegex.Split(element, 2)

		value := strings.Join(splitTuple[1:], tupleSeparator)

		if inMap := parsed[splitTuple[0]]; inMap != "" {

			logrus.WithFields(logrus.Fields{"keyOverwritten": splitTuple[0],
				"oldValue": inMap, "newValue": value}).
				Warn("overwriting value")
		}

		if len(splitTuple) >= 1 {
			parsed[splitTuple[0]] = value

		} else {
			err = errors.New(fmt.Sprintf(
				"key/value separator '%s' not found for element '%s'",
				tupleSeparator, element,
			))
			break
		}
	}

	return
}

func getOutputFile(fileName string)(file *os.File, err error) {

	if fileName == "-" {
		file = os.Stdout
		return
	}

	file, outputFileOpenError := os.OpenFile(fileName,
		os.O_RDWR | os.O_CREATE, 0600)
	if outputFileOpenError != nil {
		logrus.WithFields(logrus.Fields{"error": outputFileOpenError,
			"fileName": fileName}).
			Error("could not open output for writing")
		err = errors.New("must supply output file name")
		return
	}

	return
}

func validateLocalRenderArguments(cmd *cobra.Command, args []string) (
		err error) {

	parsedMap, parseError := parseTupelsToMap(args)
	if parseError != nil {
		err = parseError
		return
	}
	localRenderRuntimeParameters.renderContextVariables = parsedMap


	sourceFileName := cmd.PersistentFlags().Lookup("source").Value.
		String()
	rawSourceTemplate, sourceFileReadError := ioutil.ReadFile(sourceFileName)
	if sourceFileReadError != nil {
		logrus.WithFields(logrus.Fields{"fileName": sourceFileName,
			"error": sourceFileReadError}).
			Error("could not open source file")
		err = errors.New("must supply source file name")
		return
	}
	localRenderRuntimeParameters.inputTemplate = string(rawSourceTemplate)

	outputFileName := cmd.PersistentFlags().Lookup("output").Value.
		String()
	outputFile, outputFileOpenError := getOutputFile(outputFileName)
	if outputFileOpenError != nil {
		err = outputFileOpenError
		return
	}
	localRenderRuntimeParameters.outputFile = outputFile

	return
}

func doLocalRender(cmd *cobra.Command, args []string) {
	defer localRenderRuntimeParameters.outputFile.Close()

	logrus.WithField("renderVariables",
		localRenderRuntimeParameters.renderContextVariables).
		Debug("rendering")

	renderEngine := templating_engine.NewTemplateEngine()
	template, templateCompilationError := templating_engine.
		NewTemplate(localRenderRuntimeParameters.inputTemplate)
	if templateCompilationError != nil {
		logrus.WithField("error", templateCompilationError).
			Fatal("could not interpret input template")
		return
	}

	output, err := renderEngine.Render(template,
		localRenderRuntimeParameters.renderContextVariables)
	if err != nil {
		logrus.WithField("error", err).
			Fatal("could not render template")
	}

	localRenderRuntimeParameters.outputFile.WriteString(output)
}
