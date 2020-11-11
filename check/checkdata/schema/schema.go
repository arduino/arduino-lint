// This file is part of arduino-check.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-check.
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
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"path/filepath"

	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/ory/jsonschema/v3"
	"github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonreference"
)

// Compile compiles the schema files specified by the filename arguments and returns the compiled schema.
func Compile(schemaFilename string, referencedSchemaFilenames []string, schemasPath *paths.Path) *jsonschema.Schema {
	compiler := jsonschema.NewCompiler()

	// Load the referenced schemas.
	for _, referencedSchemaFilename := range referencedSchemaFilenames {
		if err := loadReferencedSchema(compiler, referencedSchemaFilename, schemasPath); err != nil {
			panic(err)
		}
	}

	// Compile the schema.
	schemaPath := schemasPath.Join(schemaFilename)
	schemaURI := pathURI(schemaPath)
	compiledSchema, err := compiler.Compile(schemaURI)
	if err != nil {
		panic(err)
	}

	return compiledSchema
}

// Validate validates an instance against a JSON schema and returns nil if it was success, or the
// jsonschema.ValidationError object otherwise.
func Validate(instanceObject *properties.Map, schemaObject *jsonschema.Schema, schemasPath *paths.Path) *jsonschema.ValidationError {
	// Convert the instance data from the native properties.Map type to the interface type required by the schema
	// validation package.
	instanceObjectMap := instanceObject.AsMap()
	instanceInterface := make(map[string]interface{}, len(instanceObjectMap))
	for k, v := range instanceObjectMap {
		instanceInterface[k] = v
	}

	validationError := schemaObject.ValidateInterface(instanceInterface)
	result, _ := validationError.(*jsonschema.ValidationError)
	if result == nil {
		logrus.Debug("Schema validation of instance document passed")

	} else {
		logrus.Debug("Schema validation of instance document failed:")
		logValidationError(result, schemasPath)
		logrus.Trace("-----------------------------------------------")
	}
	return result
}

// loadReferencedSchema adds a schema that is referenced by the parent schema to the compiler object.
func loadReferencedSchema(compiler *jsonschema.Compiler, schemaFilename string, schemasPath *paths.Path) error {
	schemaPath := schemasPath.Join(schemaFilename)
	schemaFile, err := schemaPath.Open()
	if err != nil {
		return err
	}
	defer schemaFile.Close()

	// Get the $id value from the schema to use as the `url` argument for the `compiler.AddResource()` call.
	id, err := schemaID(schemaFilename, schemasPath)
	if err != nil {
		return err
	}

	return compiler.AddResource(id, schemaFile)
}

// schemaID returns the value of the schema's $id key.
func schemaID(schemaFilename string, schemasPath *paths.Path) (string, error) {
	schemaPath := schemasPath.Join(schemaFilename)
	schemaInterface := unmarshalJSONFile(schemaPath)

	id, ok := schemaInterface.(map[string]interface{})["$id"].(string)
	if !ok {
		return "", fmt.Errorf("Schema %s is missing an $id keyword", schemaPath)
	}

	return id, nil
}

// unmarshalJSONFile returns the data from a JSON file.
func unmarshalJSONFile(filePath *paths.Path) interface{} {
	fileBuffer, err := filePath.ReadFile()
	if err != nil {
		panic(err)
	}

	var dataInterface interface{}
	if err := json.Unmarshal(fileBuffer, &dataInterface); err != nil {
		panic(err)
	}

	return dataInterface
}

// compile compiles the parent schema and returns the resulting jsonschema.Schema object.
func compile(compiler *jsonschema.Compiler, schemaFilename string, schemasPath *paths.Path) (*jsonschema.Schema, error) {
	schemaPath := schemasPath.Join(schemaFilename)
	schemaURI := pathURI(schemaPath)
	return compiler.Compile(schemaURI)
}

// pathURI returns the URI representation of the path argument.
func pathURI(path *paths.Path) string {
	absolutePath, err := path.Abs()
	if err != nil {
		panic(err)
	}
	uriFriendlyPath := filepath.ToSlash(absolutePath.String())
	// In order to be valid, the path in the URI must start with `/`, but Windows paths do not.
	if uriFriendlyPath[0] != '/' {
		uriFriendlyPath = "/" + uriFriendlyPath
	}
	pathURI := url.URL{
		Scheme: "file",
		Path:   uriFriendlyPath,
	}

	return pathURI.String()
}

// logValidationError logs the schema validation error data
func logValidationError(validationError *jsonschema.ValidationError, schemasPath *paths.Path) {
	logrus.Trace("--------Schema validation failure cause--------")
	logrus.Tracef("Error message: %s", validationError.Error())
	logrus.Tracef("Instance pointer: %v", validationError.InstancePtr)
	logrus.Tracef("Schema URL: %s", validationError.SchemaURL)
	logrus.Tracef("Schema pointer: %s", validationError.SchemaPtr)
	logrus.Tracef("Schema pointer value: %v", validationErrorSchemaPointerValue(validationError, schemasPath))
	logrus.Tracef("Failure context: %v", validationError.Context)
	logrus.Tracef("Failure context type: %T", validationError.Context)

	// Recursively log all causes.
	for _, validationErrorCause := range validationError.Causes {
		logValidationError(validationErrorCause, schemasPath)
	}
}

// validationErrorSchemaPointerValue returns the object identified by the validation error's schema JSON pointer.
func validationErrorSchemaPointerValue(validationError *jsonschema.ValidationError, schemasPath *paths.Path) interface{} {
	return schemaPointerValue(validationError.SchemaURL, validationError.SchemaPtr, schemasPath)
}

// schemaPointerValue returns the object identified by the given JSON pointer from the schema file.
func schemaPointerValue(schemaURL, schemaPointer string, schemasPath *paths.Path) interface{} {
	schemaPath := schemasPath.Join(path.Base(schemaURL))
	return jsonPointerValue(schemaPointer, schemaPath)
}

// jsonPointerValue returns the object identified by the given JSON pointer from the JSON file.
func jsonPointerValue(jsonPointer string, filePath *paths.Path) interface{} {
	jsonReference, err := gojsonreference.NewJsonReference(jsonPointer)
	if err != nil {
		panic(err)
	}
	jsonInterface := unmarshalJSONFile(filePath)
	jsonPointerValue, _, err := jsonReference.GetPointer().Get(jsonInterface)
	if err != nil {
		panic(err)
	}
	return jsonPointerValue
}
