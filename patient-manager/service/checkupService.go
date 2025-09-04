package service

import (
	"PatientManager/app"
	"PatientManager/model"
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ICheckupService interface {
	Create(checkup *model.Checkup, recordUuid string) (*model.Checkup, error)
	Update(checkupUuid uuid.UUID, checkupUpdateData *model.Checkup) (*model.Checkup, error)
	GetAll(recordUuid uuid.UUID) ([]model.Checkup, error)
	Delete(checkupUuid uuid.UUID) error
	AddImagesToCheckup(checkupUuid string, files []string) (*model.Checkup, error)
	DeleteImage(imageUuid string) error
}

type CheckupService struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewChekupService() ICheckupService {
	var service ICheckupService
	app.Invoke(func(db *gorm.DB, logger *zap.SugaredLogger) {
		service = &CheckupService{
			db:     db,
			logger: logger,
		}
	})

	return service
}

func (c *CheckupService) Create(checkup *model.Checkup, recordUuid string) (*model.Checkup, error) {
	checkup.Uuid = uuid.New()
	c.logger.Infof("Creating checkup for medical record uuid: %s", recordUuid)

	var medicalRecord model.MedicalRecord
	if err := c.db.Where("uuid = ?", recordUuid).First(&medicalRecord).Error; err != nil {
		c.logger.Errorf("Error finding medical record with UUID %s: %v", recordUuid, err)
		return nil, err
	}

	checkup.MedicalRecordID = medicalRecord.ID

	rez := c.db.Create(checkup)
	if rez.Error != nil {
		c.logger.Errorf("Error creating checkup: %v", rez.Error)
		return nil, rez.Error
	}

	c.logger.Infof("Successfully created checkup with UUID: %s", checkup.Uuid)
	return checkup, nil
}

func (c *CheckupService) findByUuid(checkupUuid uuid.UUID) (*model.Checkup, error) {
	var checkup model.Checkup
	rez := c.db.
		Preload("Images").
		Where("uuid = ?", checkupUuid).
		First(&checkup)

	if rez.Error != nil {
		if rez.Error == gorm.ErrRecordNotFound {
			c.logger.Warnf("Checkup with UUID %s not found", checkupUuid)
		} else {
			c.logger.Errorf("Error finding checkup with UUID %s: %v", checkupUuid, rez.Error)
		}
		return nil, rez.Error
	}
	return &checkup, nil
}

func (c *CheckupService) Update(checkupUuid uuid.UUID, checkupUpdateData *model.Checkup) (*model.Checkup, error) {
	existingCheckup, err := c.findByUuid(checkupUuid)
	if err != nil {
		return nil, err
	}

	c.logger.Debugf("Updating checkup with UUID: %s", checkupUuid)

	existingCheckup.UpdateCheckup(checkupUpdateData)

	rez := c.db.Save(existingCheckup)
	if rez.Error != nil {
		c.logger.Errorf("Error saving updated checkup with UUID %s: %v", checkupUuid, rez.Error)
		return nil, rez.Error
	}

	c.logger.Infof("Successfully updated checkup with UUID: %s", checkupUuid)
	return existingCheckup, nil
}

func (c *CheckupService) Delete(checkupUuid uuid.UUID) error {
	c.logger.Infof("Attempting to delete checkup with UUID: %s", checkupUuid)

	rez := c.db.Where("uuid = ?", checkupUuid).Delete(&model.Checkup{})

	if rez.Error != nil {
		c.logger.Errorf("Error deleting checkup with UUID %s: %v", checkupUuid, rez.Error)
		return rez.Error
	}

	if rez.RowsAffected == 0 {
		c.logger.Warnf("No checkup found with UUID %s to delete", checkupUuid)
		return gorm.ErrRecordNotFound
	}

	c.logger.Infof("Successfully checkup with UUID: %s", checkupUuid)
	return nil
}

func (c *CheckupService) GetAll(recordUuid uuid.UUID) ([]model.Checkup, error) {
	c.logger.Infof("Fetching all checkups for medical record uuid: %s", recordUuid)

	var medicalRecord model.MedicalRecord
	if err := c.db.Where("uuid = ?", recordUuid).First(&medicalRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.logger.Warnf("Medical record with UUID %s not found", recordUuid)
		} else {
			c.logger.Errorf("Error finding medical record with UUID %s: %v", recordUuid, err)
		}
		return nil, err
	}

	var checkups []model.Checkup
	rez := c.db.Preload("MedicalRecord").
		Preload("Images").
		Where("medical_record_id = ?", medicalRecord.ID).
		Order("checkup_date desc").
		Find(&checkups)
	if rez.Error != nil {
		c.logger.Errorf("Error fetching checkups for medical record ID %d: %v", medicalRecord.ID, rez.Error)
		return nil, rez.Error
	}

	c.logger.Infof("Successfully fetched %d checkups for medical record uuid: %s", len(checkups), recordUuid)
	return checkups, nil
}

func (c *CheckupService) AddImagesToCheckup(checkupUuid string, paths []string) (*model.Checkup, error) {
	parsedUuid, err := uuid.Parse(checkupUuid)
	if err != nil {
		c.logger.Errorf("Failed to parse checkup UUID %s: %v", checkupUuid, err)
		return nil, err
	}

	var checkup model.Checkup
	if err := c.db.Where("uuid = ?", parsedUuid).First(&checkup).Error; err != nil {
		c.logger.Errorf("Checkup with UUID %s not found: %v", checkupUuid, err)
		return nil, err
	}

	for _, path := range paths {
		image := model.Image{
			Uuid:      uuid.New(),
			Path:      path,
			CheckupID: checkup.ID,
		}
		if err := c.db.Create(&image).Error; err != nil {
			c.logger.Errorf("Failed to create image record for checkup %s: %v", checkupUuid, err)
			return nil, err
		}
	}

	return c.findByUuid(parsedUuid)
}

func (c *CheckupService) DeleteImage(imageUuid string) error {
	var image model.Image
	if err := c.db.Where("uuid = ?", imageUuid).First(&image).Error; err != nil {
		return err
	}
	if err := os.Remove(image.Path); err != nil {
		c.logger.Warnf("Failed to delete image file from disk: %v", err)
	}

	return c.db.Delete(&image).Error
}
