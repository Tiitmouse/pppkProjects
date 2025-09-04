package repo

import (
	"context"
	"data-managment/util/env"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var RepoCli *MongoRepo

type MongoRepo struct {
	dbLock   sync.Mutex
	mongoDb  *mongo.Client
	patients *mongo.Collection
}

// TODO: change names
const database = "gene-data"
const collection = "genes"

func Setup() error {
	zap.S().Debugf("Setting up mongodb")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(env.MongoDb_Conn_String))
	if err != nil {
		return err
	}

	patients := client.Database(database).Collection(collection)

	RepoCli = &MongoRepo{
		dbLock:   sync.Mutex{},
		mongoDb:  client,
		patients: patients,
	}

	zap.S().Debugf("Mongo setup compleat")
	return nil
}

func (m *MongoRepo) Disconnect() error {
	if err := m.mongoDb.Disconnect(context.Background()); err != nil {
		return err
	}
	return nil
}

func (m *MongoRepo) AddPatients(items []any) error {
	_, err := m.patients.InsertMany(context.Background(), items)
	if err != nil {
		return err
	}

	return nil
}

type PatientData struct {
	BCRPatientBarcode string `json:"bcr_patient_barcode" bson:"bcr_patient_barcode,omitempty"`
	DSS               bool   `json:"dss" bson:"dss,omitempty"`
	OS                bool   `json:"os" bson:"os,omitempty"`
	ClinicalStage     string `json:"clinical_stage" bson:"clinical_stage,omitempty"`
}

func (m *MongoRepo) Get(code string) (PatientData, error) {
	m.dbLock.Lock()
	defer m.dbLock.Unlock()

	var patient PatientData
	filter := bson.M{"bcr_patient_barcode": code}

	err := m.patients.FindOne(context.Background(), filter).Decode(&patient)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			zap.S().Errorf("no patient found with code %s", code)
			return PatientData{}, err
		}
		return PatientData{}, err
	}

	return patient, nil
}

func (m *MongoRepo) GetAll() ([]PatientData, error) {
	cursor, err := m.patients.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var result []PatientData
	if err = cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (m *MongoRepo) Clear() error {
	rez, err := m.patients.DeleteMany(context.Background(), bson.D{})
	zap.S().Infof("Deleted %s items", rez.DeletedCount)
	if err != nil {
		zap.S().Errorf("Failed to delete items, err = %v", err)
		return err
	}

	return nil
}
