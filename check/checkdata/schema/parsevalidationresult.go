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
	"fmt"
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
		matchedSchemaPointer := validationErrorSchemaPointerMatch(schemaPointerRegexp, validationError, schemasPath)
		if matchedSchemaPointer != "" {
			logrus.Tracef("Matched!")
			if validationErrorSchemaPointerValueMatch(schemaPointerValueRegexp, validationError.SchemaURL, matchedSchemaPointer, schemasPath) {
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

// validationErrorSchemaPointerMatch matches the JSON schema pointer related to the validation failure against a regular expression.
func validationErrorSchemaPointerMatch(
	schemaPointerRegexp *regexp.Regexp,
	validationError *jsonschema.ValidationError,
	schemasPath *paths.Path,
) string {
	logrus.Tracef("Checking schema pointer: %s match with regexp: %s", validationError.SchemaPtr, schemaPointerRegexp)
	if schemaPointerRegexp.MatchString(validationError.SchemaPtr) {
		return validationError.SchemaPtr
	}

	// The schema validator does not provide full pointer past logic inversion keywords to the lowest level keywords related to the validation error cause.
	// Therefore the sub-keywords must be checked for matches in order to be able to interpret the exact cause of the failure.
	if regexp.MustCompile("(/not)|(/oneOf)$").MatchString(validationError.SchemaPtr) {
		return validationErrorSchemaSubPointerMatch(schemaPointerRegexp, validationError.SchemaPtr, validationErrorSchemaPointerValue(validationError, schemasPath))
	}

	return ""
}

// validationErrorSchemaSubPointerMatch recursively checks JSON pointers of all keywords under the parent pointer for match against a regular expression.
// The matching JSON pointer is returned.
func validationErrorSchemaSubPointerMatch(schemaPointerRegexp *regexp.Regexp, parentPointer string, pointerValueObject interface{}) string {
	// Recurse through iterable objects.
	switch assertedObject := pointerValueObject.(type) {
	case []interface{}:
		for index, element := range assertedObject {
			// Append index to JSON pointer and check for match.
			matchingPointer := validationErrorSchemaSubPointerMatch(schemaPointerRegexp, fmt.Sprintf("%s/%d", parentPointer, index), element)
			if matchingPointer != "" {
				return matchingPointer
			}
		}
	case map[string]interface{}:
		for key := range assertedObject {
			// Append key to JSON pointer and check for match.
			matchingPointer := validationErrorSchemaSubPointerMatch(schemaPointerRegexp, parentPointer+"/"+key, assertedObject[key])
			if matchingPointer != "" {
				return matchingPointer
			}
			// TODO: Follow references. For now, the schema code must be written so that the problematic keywords are after the reference.
		}
	}

	// pointerValueObject is not further iterable. Check for match against the parent JSON pointer.
	logrus.Tracef("Checking schema pointer: %s match with regexp: %s", parentPointer, schemaPointerRegexp)
	if schemaPointerRegexp.MatchString(parentPointer) {
		return parentPointer
	}
	return ""
}

// validationErrorSchemaPointerValueMatch marshalls the data in the schema at the given JSON pointer and returns whether
// it matches against the given regular expression.
func validationErrorSchemaPointerValueMatch(
	schemaPointerValueRegexp *regexp.Regexp,
	schemaURL,
	schemaPointer string,
	schemasPath *paths.Path,
) bool {
	marshalledSchemaPointerValue, err := json.Marshal(schemaPointerValue(schemaURL, schemaPointer, schemasPath))
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
