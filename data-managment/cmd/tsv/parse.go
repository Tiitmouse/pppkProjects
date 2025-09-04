package tsv

import (
	"data-managment/util/repo"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path"
	"strconv"

	"go.uber.org/zap"
)

const separator = '\t'

func action(filePath string) error {

	if path.Ext(filePath) != ".tsv" {
		zap.S().Errorf("File is not tsv")
		return errors.New("file is not tsv")
	}

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0755)
	if err != nil {
		zap.S().Errorf("Failed to open file %s", filePath)
		return err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = separator

	data, err := parseClinicalSurvivalData(csvReader)
	if err != nil {
		return err
	}
	zap.S().Infof("parsed %d entries", len(data))

	interfaceSlice := make([]any, len(data))
	for i, v := range data {
		interfaceSlice[i] = v
	}

	err = repo.RepoCli.AddPatients(interfaceSlice)
	if err != nil {
		zap.S().Errorf("Error saving patient data, err = %v", err)
		return err
	}

	zap.S().Infof("Data successfully uploaded")

	return nil
}

func parseClinicalSurvivalData(reader *csv.Reader) ([]repo.PatientData, error) {
	var arr []repo.PatientData

	_, err := reader.Read()
	if err != nil {
		zap.S().Errorf("failed to read header err = %v", err)
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			zap.S().Debugf("Finished parsing")
			break
		}
		if err != nil {
			zap.S().Errorf("failed to read row err = %v", err)
			return nil, err
		}

		arr = append(arr, repo.PatientData{
			BCRPatientBarcode: record[1],
			ClinicalStage:     record[7],
			OS:                convertToBool(record[25]),
			DSS:               convertToBool(record[27]),
		})
	}

	return arr, nil
}

func convertToBool(str string) bool {
	br, _ := strconv.Atoi(str)
	return br == 1
}
