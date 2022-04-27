package encodingutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

func validateBytes(bytes []byte, schemabytes []byte) error {

	err := EnsureJson(&bytes, false)
	if err != nil {
		return fmt.Errorf("can't parse input: %s", err.Error())
	}

	var obj interface{}
	if err = json.Unmarshal(bytes, &obj); err != nil {
		return fmt.Errorf("can't unmarshal data: %s", err.Error())
	}

	if len(schemabytes) > 0 {
		schemaLoader := gojsonschema.NewStringLoader(string(schemabytes))
		documentLoader := gojsonschema.NewStringLoader(string(bytes))

		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			return fmt.Errorf("can't validate JSON: %s", err.Error())
		}

		if !result.Valid() {
			var report string
			for i, desc := range result.Errors() {
				if i > 0 {
					report += "; "
				}
				report += fmt.Sprintf("%s", desc)
			}
			return fmt.Errorf("invalid JSON: %s", report)
		}
	} else {
		log.Printf("WARN: checking syntax only\n")
	}

	return nil
}

func ValidateFile(path string, schema string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("can't read %s: %v", path, err)
	}

	schemabytes, err := loadSchema(schema)
	if err != nil {
		return fmt.Errorf("can't parse schema: %s", err.Error())
	}

	log.Println(fmt.Sprintf("Validating %s...", path))
	return validateBytes(bytes, schemabytes)
}

func validateSTDIN(file *os.File, schema string) error {
	var stdin []byte
	stdin, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("can't read input stream: %s", err.Error())
	}

	schemabytes, err := loadSchema(schema)
	if err != nil {
		return fmt.Errorf("can't parse schema: %s", err.Error())
	}

	log.Println("Validating stream...")
	return validateBytes(stdin, schemabytes)
}

func loadSchema(schema string) ([]byte, error) {
	var schemabytes []byte
	var err error
	if len(schema) > 0 {
		log.Printf("Loading schema %s...", schema)
		schemabytes, err = ioutil.ReadFile(schema)

		schemaIsJSON := strings.HasSuffix(schema, ".json")
		err = EnsureJson(&schemabytes, schemaIsJSON)
		if err != nil {
			return nil, fmt.Errorf("can't parse schema: %s", err.Error())
		}
	}

	return schemabytes, nil
}
