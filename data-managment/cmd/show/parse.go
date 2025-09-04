package show

import (
	"encoding/csv"
	"io"
	"slices"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

var wantedGenes = []string{"C6orf150", "CCL5", "CXCL10", "TMEM173", "CXCL9", "CXCL11", "NFKB1", "IKBKE", "IRF3", "TREX1", "ATM", "IL6", "IL8"}

func read(reader io.Reader, patientCode string) (genesExpressions, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = '\t'

	header, err := csvReader.Read()
	if err != nil {
		if err == io.EOF {
			zap.S().Errorf("Input data is empty")
			return genesExpressions{}, err
		}
		zap.S().Errorf("Failed to read header row: %v", err)
		return genesExpressions{}, err
	}

	patientColumnIndex := -1
	for i, code := range header {
		if strings.Contains(code, patientCode) {
			patientColumnIndex = i
			break
		}
	}

	if patientColumnIndex == -1 {
		zap.S().Debugf("patient code '%s' not found in header", patientCode)
		return genesExpressions{}, errNotFound
	}

	var geneData []genePair
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			zap.S().Errorf("error reading data row: %w", err)
			return genesExpressions{}, err
		}

		geneName := record[0]
		if !slices.Contains(wantedGenes, geneName) {
			continue
		}

		expressionStr := record[patientColumnIndex]
		expression, err := strconv.ParseFloat(expressionStr, 64)
		if err != nil {
			zap.S().Errorf("could not parse expression '%s' for gene '%s': %w", expressionStr, geneName, err)
			return genesExpressions{}, err
		}

		geneData = append(geneData, genePair{
			Gene:       geneName,
			Expression: expression,
		})
	}

	result := genesExpressions{
		BCRPatientBarcode: patientCode,
		Genes:             geneData,
	}

	return result, nil
}
