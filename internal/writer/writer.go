// Package writer is the package that contains the writer functionality.
package writer

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/aiven/go-api-schemas/internal/pkg/types"
	"github.com/aiven/go-api-schemas/internal/pkg/util"
)

// logger is a pointer to the logger.
var logger *util.Logger

// flags is a pointer to the flags.
var flags *pflag.FlagSet

// result is the result of the diff process.
var result types.DiffResult

// write is a function that writes map[string]types.UserConfigSchema to a file.
func write(filename string, schema map[string]types.UserConfigSchema) error {
	// indentSpaces is the number of spaces to use for indentation.
	const indentSpaces = 2

	logger.Info.Printf("writing %s", filename)

	outputDir, err := flags.GetString("output-dir")
	if err != nil {
		return err
	}

	f, err := os.Create(strings.Join([]string{outputDir, filename}, string(os.PathSeparator)))
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err = f.Close()
	}(f)

	e := yaml.NewEncoder(f)

	defer func(e *yaml.Encoder) {
		err = e.Close()
	}(e)

	e.SetIndent(indentSpaces)

	if err = e.Encode(schema); err != nil {
		return err
	}

	return err
}

// writeServiceTypes writes the service types to a file.
func writeServiceTypes() error {
	return write(util.ServiceTypesFilename, result[types.KeyServiceTypes])
}

// writeIntegrationTypes writes the integration types to a file.
func writeIntegrationTypes() error {
	return write(util.IntegrationTypesFilename, result[types.KeyIntegrationTypes])
}

// writeIntegrationEndpointTypes writes the integration endpoint types to a file.
func writeIntegrationEndpointTypes() error {
	return write(util.IntegrationEndpointTypesFilename, result[types.KeyIntegrationEndpointTypes])
}

// setup sets up the writer.
func setup(l *util.Logger, f *pflag.FlagSet, r types.DiffResult) {
	logger = l
	flags = f
	result = r
}

// Run runs the writer.
func Run(logger *util.Logger, flags *pflag.FlagSet, result types.DiffResult) error {
	setup(logger, flags, result)

	if err := writeServiceTypes(); err != nil {
		return err
	}

	if err := writeIntegrationTypes(); err != nil {
		return err
	}

	if err := writeIntegrationEndpointTypes(); err != nil {
		return err
	}

	return nil
}
