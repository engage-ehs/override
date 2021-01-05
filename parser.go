// Package override creates a generic parser for command-line arguments to set fields in Go data
// structure. This is particularly useful when dealing with large configuration objects, with
// command-line values override default settings.
package override

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// Scan overrides values in the config object by reading command line arguments formatted as
// key=value pairs:
//
//  aws_s3bucket=my-bucket aws_region=eu-west-1
//
// Fields are matched in config structure by name after normalization of arguments (underscore `_`
// and dash `-` are dropped, all characters are lower-cased).
//
// When it is desired that a value cannot be set from the inputs, the tag `canset:"no"` should be
// added to the field (cf example).
func Scan(input []string, config interface{}) error {
	settings := make(map[string]string)

	for _, key := range input {
		i := strings.IndexRune(key, '=')
		if i == -1 {
			return fmt.Errorf("invalid argument in %s", key)
		}

		settings[key[:i]] = key[i+1:]
	}

	vcfg := reflect.Indirect(reflect.ValueOf(config))
	tycfg := vcfg.Type()

	for k, v := range settings {
		field, ok := tycfg.FieldByNameFunc(match(k))
		if !ok {
			return fmt.Errorf("no such field %s", k)
		}
		if val, canset := field.Tag.Lookup("canset"); canset && val == "no" {
			return cannotSetViolation(k)
		}

		vfield := vcfg.FieldByName(field.Name)
		if !vfield.CanSet() {
			return errors.New("cannot assign to field")
		}
		vfield.Set(reflect.ValueOf(v))
	}

	return nil
}

type cannotSetViolation string

func (err cannotSetViolation) Error() string {
	return fmt.Sprintf("field %s cannot be configured", string(err))
}

func match(key string) func(string) bool {
	return func(field string) bool {
		return norm(key) == norm(field)
	}
}

func norm(in string) string {
	var out strings.Builder

	for _, r := range in {
		switch {
		case r == '_', r == '-':
			// strip
		case unicode.IsUpper(r):
			out.WriteRune(unicode.ToLower(r))
		default:
			out.WriteRune(r)
		}
	}

	return out.String()
}
