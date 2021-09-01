// This file is part of Arduino Lint.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of Arduino Lint.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

// Package schema contains code for working with JSON schema.
package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path"

	"github.com/ory/jsonschema/v3"
	"github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonreference"
)

// dataLoaderType is the signature of the function that returns the byte encoded data associated with the given file name.
type dataLoaderType func(filename string) ([]byte, error)

// Schema is the type of the compiled JSON schema object.
type Schema struct {
	Compiled   *jsonschema.Schema
	dataLoader dataLoaderType // Function to load the schema data.
}

// ValidationResult is the type of the result of the validation of the instance document against the JSON schema.
type ValidationResult struct {
	Result     *jsonschema.ValidationError
	dataLoader dataLoaderType // Function used to load the JSON schema data.
}

// Compile compiles the schema files specified by the filename arguments and returns the compiled schema.
func Compile(schemaFilename string, referencedSchemaFilenames []string, dataLoader dataLoaderType) Schema {
	compiler := jsonschema.NewCompiler()

	// Define a custom schema loader for the binary encoded schema.
	compiler.LoadURL = func(schemaFilename string) (io.ReadCloser, error) {
		schemaData, err := dataLoader(schemaFilename)
		if err != nil {
			return nil, err
		}

		return ioutil.NopCloser(bytes.NewReader(schemaData)), nil
	}

	// Load the referenced schemas.
	for _, referencedSchemaFilename := range referencedSchemaFilenames {
		if err := loadReferencedSchema(compiler, referencedSchemaFilename, dataLoader); err != nil {
			panic(err)
		}
	}

	// Compile the schema.
	compiledSchema, err := compiler.Compile(schemaFilename)
	if err != nil {
		panic(err)
	}

	return Schema{
		Compiled:   compiledSchema,
		dataLoader: dataLoader,
	}
}

// Validate validates an instance against a JSON schema and returns nil if it was success, or the
// jsonschema.ValidationError object otherwise.
func Validate(instanceInterface map[string]interface{}, schemaObject Schema) ValidationResult {
	validationError := schemaObject.Compiled.ValidateInterface(instanceInterface)
	validationResult := ValidationResult{
		Result:     nil,
		dataLoader: schemaObject.dataLoader,
	}

	if validationError != nil {
		result, ok := validationError.(*jsonschema.ValidationError)
		if !ok {
			panic(validationError)
		}
		validationResult.Result = result
	}

	if validationResult.Result == nil {
		logrus.Debug("Schema validation of instance document passed")
	} else {
		logrus.Debug("Schema validation of instance document failed:")
		logValidationError(validationResult)
		logrus.Trace("-----------------------------------------------")
	}
	return validationResult
}

// loadReferencedSchema adds a schema that is referenced by the parent schema to the compiler object.
func loadReferencedSchema(compiler *jsonschema.Compiler, schemaFilename string, dataLoader dataLoaderType) error {
	// Get the $id value from the schema to use as the `url` argument for the `compiler.AddResource()` call.
	id, err := schemaID(schemaFilename, dataLoader)
	if err != nil {
		return err
	}

	schemaData, err := dataLoader(schemaFilename)
	if err != nil {
		return err
	}

	return compiler.AddResource(id, bytes.NewReader(schemaData))
}

// schemaID returns the value of the schema's $id key.
func schemaID(schemaFilename string, dataLoader dataLoaderType) (string, error) {
	schemaInterface := unmarshalJSONFile(schemaFilename, dataLoader)

	id, ok := schemaInterface.(map[string]interface{})["$id"].(string)
	if !ok {
		return "", fmt.Errorf("Schema %s is missing an $id keyword", schemaFilename)
	}

	return id, nil
}

// unmarshalJSONFile returns the data from a JSON file.
func unmarshalJSONFile(filename string, dataLoader dataLoaderType) interface{} {
	data, err := dataLoader(filename)
	if err != nil {
		panic(err)
	}

	var dataInterface interface{}
	if err := json.Unmarshal(data, &dataInterface); err != nil {
		panic(err)
	}

	return dataInterface
}

// logValidationError logs the schema validation error data.
func logValidationError(validationError ValidationResult) {
	logrus.Trace("--------Schema validation failure cause--------")
	logrus.Tracef("Error message: %s", validationError.Result.Error())
	logrus.Tracef("Instance pointer: %v", validationError.Result.InstancePtr)
	logrus.Tracef("Schema URL: %s", validationError.Result.SchemaURL)
	logrus.Tracef("Schema pointer: %s", validationError.Result.SchemaPtr)
	logrus.Tracef("Schema pointer value: %v", validationErrorSchemaPointerValue(validationError))
	logrus.Tracef("Failure context: %v", validationError.Result.Context)
	logrus.Tracef("Failure context type: %T", validationError.Result.Context)

	// Recursively log all causes.
	for _, validationErrorCause := range validationError.Result.Causes {
		logValidationError(
			ValidationResult{
				Result:     validationErrorCause,
				dataLoader: validationError.dataLoader,
			},
		)
	}
}

// validationErrorSchemaPointerValue returns the object identified by the validation error's schema JSON pointer.
func validationErrorSchemaPointerValue(validationError ValidationResult) interface{} {
	return schemaPointerValue(validationError.Result.SchemaURL, validationError.Result.SchemaPtr, validationError.dataLoader)
}

// schemaPointerValue returns the object identified by the given JSON pointer from the schema file.
func schemaPointerValue(schemaURL, schemaPointer string, dataLoader dataLoaderType) interface{} {
	return jsonPointerValue(schemaPointer, path.Base(schemaURL), dataLoader)
}

// jsonPointerValue returns the object identified by the given JSON pointer from the JSON file.
func jsonPointerValue(jsonPointer string, fileName string, dataLoader dataLoaderType) interface{} {
	jsonReference, err := gojsonreference.NewJsonReference(jsonPointer)
	if err != nil {
		panic(err)
	}
	jsonInterface := unmarshalJSONFile(fileName, dataLoader)
	jsonPointerValue, _, err := jsonReference.GetPointer().Get(jsonInterface)
	if err != nil {
		panic(err)
	}
	return jsonPointerValue
}
