package show

import (
	"data-managment/util/bucket"
	"data-managment/util/repo"
	"errors"

	"go.uber.org/zap"
)

var errNotFound error = errors.New("patient code not found")

type genesExpressions struct {
	BCRPatientBarcode string     `json:"bcr_patient_barcode"`
	Genes             []genePair `json:"genes"`
}

type genePair struct {
	Gene       string  `json:"gene"`
	Expression float64 `json:"expression"`
}

func connectData(patientCode string) error {
	patientData, err := repo.RepoCli.Get(patientCode)
	if err != nil {
		return err
	}

	zap.S().Debugf("Found patient (code = %s) %+v ", patientCode, patientData)

	files, err := bucket.Bucket.GetFiles()
	if err != nil {
		zap.S().Errorf("Failed to retrieve files err = %v", err)
		return err
	}
	zap.S().Debugf("Reading through %d files", len(files))

	var geneData genesExpressions
	for _, file := range files {
		rez, err := read(file, patientCode)
		defer file.Close()
		if err != nil {
			if err == errNotFound {
				continue
			}
			zap.S().Errorf("Failed to read file, err = %v", err)
			return err
		}
		geneData = rez
	}

	zap.S().Debugf("Found data for %s, %+v ", patientCode, geneData)

	if err := draw(patientData, geneData); err != nil {
		return err
	}

	return nil
}
