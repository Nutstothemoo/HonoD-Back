package qrcode

import (
	"context"
	"fmt"
	"ginapp/internal/infrastructure/s3"
	"log"
	"os"
	"sync"
)

type QRCodeService struct {
	repo s3.S3Repository
}

func NewQRCodeService(repo s3.S3Repository) *QRCodeService {
	return &QRCodeService{repo: repo}
}

func (s *QRCodeService) GenerateAndSaveQRCodeOnDisk(data string) error {
	qr := QRCode{Data: data}
	png, err := qr.Generate()
	if err != nil {
		return err
	}

	err = save("qrcode.png", png)
	if err != nil {
		return err
	}

	fmt.Println("QR code generated and saved as 'qrcode.png'")
	return nil
}

func (s *QRCodeService) GenerateAndUploadQRCode(ctx context.Context, data, s3BucketName, s3ObjectKeyPrefix string) error {
	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	encryptedData, err := EncryptData(data, encryptionKey)
	if err != nil {
		return err
	}

	qr := QRCode{Data: encryptedData}
	png, err := qr.Generate()
	if err != nil {
		return err
	}

	err = s.repo.Upload(ctx, s3BucketName, s3ObjectKeyPrefix, png)
	if err != nil {
		return err
	}

	fmt.Printf("QR code uploaded to S3 bucket %s with key %s\n", s3BucketName, s3ObjectKeyPrefix)
	return nil
}

func (s *QRCodeService) GenerateAndUploadMultipleQRCodes(data map[string]string, s3BucketName string) error {
	var wg sync.WaitGroup

	for dataValue, s3ObjectKeyPrefix := range data {
		wg.Add(1)
		go func(dataValue, s3ObjectKeyPrefix string) {
			defer wg.Done()

			qr := QRCode{Data: dataValue}
			png, err := qr.Generate()
			if err != nil {
				log.Printf("Error generating QR code for %s: %v\n", dataValue, err)
				return
			}

			err = s.repo.Upload(context.TODO(), s3BucketName, s3ObjectKeyPrefix+".png", png)
			if err != nil {
				log.Printf("Error uploading QR code for %s to S3: %v\n", dataValue, err)
				return
			}

			fmt.Printf("QR code for %s uploaded to S3 bucket %s with key %s\n", dataValue, s3BucketName, s3ObjectKeyPrefix)
		}(dataValue, s3ObjectKeyPrefix)
	}

	wg.Wait()
	return nil
}

func save(filename string, pngData []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(pngData)
	if err != nil {
		return err
	}

	return nil
}
