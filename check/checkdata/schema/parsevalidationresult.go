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

package schema

import (
	"encoding/json"
	"regexp"

	"github.com/arduino/go-paths-helper"
	"github.com/ory/jsonschema/v3"
	"github.com/sirupsen/logrus"
)

// RequiredPropertyMissing returns whether the given required property is missing from the document.
func RequiredPropertyMissing(propertyName string, validationResult *jsonschema.ValidationError, schemasPath *paths.Path) bool {
	return ValidationErrorMatch("#", "/required$", "", "^#/"+propertyName+"$", validationResult, schemasPath)
}

// PropertyPatternMismatch returns whether the given property did not match the regular expression defined in the JSON schema.
func PropertyPatternMismatch(propertyName string, validationResult *jsonschema.ValidationError, schemasPath *paths.Path) bool {
	return ValidationErrorMatch("#/"+propertyName, "/pattern$", "", "", validationResult, schemasPath)
}

// ValidationErrorMatch returns whether the given query matches against the JSON schema validation error.
// See: https://godoc.org/github.com/ory/jsonschema#ValidationError
func ValidationErrorMatch(
	instancePointerQuery,
	schemaPointerQuery,
	schemaPointerValueQuery,
	failureContextQuery string,
	validationResult *jsonschema.ValidationError,
	schemasPath *paths.Path,
) bool {
	if validationResult == nil {
		// No error, so nothing to match
		logrus.Trace("Schema validation passed. No match is possible.")
		return false
	}

	instancePointerRegexp := regexp.MustCompile(instancePointerQuery)
	schemaPointerRegexp := regexp.MustCompile(schemaPointerQuery)
	schemaPointerValueRegexp := regexp.MustCompile(schemaPointerValueQuery)
	failureContextRegexp := regexp.MustCompile(failureContextQuery)

	return validationErrorMatch(
		instancePointerRegexp,
		schemaPointerRegexp,
		schemaPointerValueRegexp,
		failureContextRegexp,
		validationResult,
		schemasPath)
}

func validationErrorMatch(
	instancePointerRegexp,
	schemaPointerRegexp,
	schemaPointerValueRegexp,
	failureContextRegexp *regexp.Regexp,
	validationError *jsonschema.ValidationError,
	schemasPath *paths.Path,
) bool {
	logrus.Trace("--------Checking schema validation failure match--------")
	logrus.Tracef("Checking instance pointer: %s match with regexp: %s", validationError.InstancePtr, instancePointerRegexp)
	if instancePointerRegexp.MatchString(validationError.InstancePtr) {
		logrus.Tracef("Matched!")
		logrus.Tracef("Checking schema pointer: %s match with regexp: %s", validationError.SchemaPtr, schemaPointerRegexp)
		if schemaPointerRegexp.MatchString(validationError.SchemaPtr) {
			logrus.Tracef("Matched!")
			if validationErrorSchemaPointerValueMatch(schemaPointerValueRegexp, validationError, schemasPath) {
				logrus.Tracef("Matched!")
				logrus.Tracef("Checking failure context: %v match with regexp: %s", validationError.Context, failureContextRegexp)
				if validationErrorContextMatch(failureContextRegexp, validationError) {
					logrus.Tracef("Matched!")
					return true
				}
			}
		}
	}

	// Recursively check all causes for a match.
	for _, validationErrorCause := range validationError.Causes {
		if validationErrorMatch(
			instancePointerRegexp,
			schemaPointerRegexp,
			schemaPointerValueRegexp,
			failureContextRegexp,
			validationErrorCause,
			schemasPath,
		) {
			return true
		}
	}

	return false
}

// validationErrorSchemaPointerValueMatch marshalls the data in the schema at the given JSON pointer and returns whether
// it matches against the given regular expression.
func validationErrorSchemaPointerValueMatch(
	schemaPointerValueRegexp *regexp.Regexp,
	validationError *jsonschema.ValidationError,
	schemasPath *paths.Path,
) bool {
	marshalledSchemaPointerValue, err := json.Marshal(schemaPointerValue(validationError, schemasPath))
	logrus.Tracef("Checking schema pointer value: %s match with regexp: %s", marshalledSchemaPointerValue, schemaPointerValueRegexp)
	if err != nil {
		panic(err)
	}
	return schemaPointerValueRegexp.Match(marshalledSchemaPointerValue)
}

// validationErrorContextMatch parses the validation error context data and returns whether it matches against the given
// regular expression.
func validationErrorContextMatch(failureContextRegexp *regexp.Regexp, validationError *jsonschema.ValidationError) bool {
	// This was added in the github.com/ory/jsonschema fork of github.com/santhosh-tekuri/jsonschema
	// It currently only provides context about the `required` keyword.
	switch contextObject := validationError.Context.(type) {
	case nil:
		return failureContextRegexp.MatchString("")
	case *jsonschema.ValidationErrorContextRequired:
		return validationErrorContextRequiredMatch(failureContextRegexp, contextObject)
	default:
		logrus.Errorf("Unhandled validation error context type: %T", validationError.Context)
		return failureContextRegexp.MatchString("")
	}
}

// validationErrorContextRequiredMatch returns whether any of the JSON pointers of missing required properties match
// against the given regular expression.
func validationErrorContextRequiredMatch(
	failureContextRegexp *regexp.Regexp,
	contextObject *jsonschema.ValidationErrorContextRequired,
) bool {
	// See: https://godoc.org/github.com/ory/jsonschema#ValidationErrorContextRequired
	for _, requiredPropertyPointer := range contextObject.Missing {
		if failureContextRegexp.MatchString(requiredPropertyPointer) {
			return true
		}
	}
	return false
}
