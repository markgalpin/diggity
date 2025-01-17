package slsa

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
)

const (
	provenanceType = "https://in-toto.io/Statement"
)

var ProvenanceMetadata map[string]interface{}

// Provenance adds provenance metadata to SBOM result
func Provenance(req *bom.ParserRequirements) *model.SLSA {
	provenance, err := parseProvenanceMetadata(*req.Arguments.Provenance)
	if err != nil {
		err = errors.New("provenance: " + err.Error())
		*req.Errors = append(*req.Errors, err)
		return nil
	}
	return &model.SLSA{
		Provenance: provenance,
	}
}

// Parse provenance metadata
func parseProvenanceMetadata(filename string) (map[string]interface{}, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Validate if file is valid JSON
	if !json.Valid(file) {
		return nil, nil
	}

	if err = json.Unmarshal(file, &ProvenanceMetadata); err != nil {
		return nil, err
	}

	// Validate provenance metadata
	if !validateProvenance(ProvenanceMetadata) {
		return nil, errors.New(filename + " contains invalid provenance metadata.")
	}

	return ProvenanceMetadata, nil
}

// Check if metadata is valid provenance
func validateProvenance(metadata map[string]interface{}) bool {
	if _, ok := metadata["_type"]; !ok {
		return false
	}
	if _, ok := metadata["subject"]; !ok {
		return false
	}
	if _, ok := metadata["predicateType"]; !ok {
		return false
	}
	if _, ok := metadata["predicate"]; !ok {
		return false
	}
	if !strings.Contains(metadata["_type"].(string), provenanceType) {
		return false
	}
	return true
}
