package utils

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/skip2/go-qrcode"
)
type BucketBasics struct {
	S3Client *s3.Client
}

func save(filename string, data []byte) error {
	
	file, err := os.Create(filename)
	if err != nil {
			return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
			return err
	}

	return nil
}

func NewS3Client() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
			return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}

func (basics BucketBasics) GenerateAndSaveQRCodeOnDisk( data string) error {

	code, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
			fmt.Println("Erreur lors de la création du QR code:", err)
			return err
	}
	png, err := code.PNG(256)
	if err != nil {
			fmt.Println("Erreur lors de la conversion du QR code en PNG:", err)
			return err
	}

	err = save("qrcode.png", png)
	if err != nil {
			fmt.Println("Erreur lors de l'enregistrement du fichier PNG:", err)
			return err
	}

	fmt.Println("QR code généré avec succès et enregistré sous 'qrcode.png'")
	return nil
}

func (basics BucketBasics) GenerateAndUploadQRCode( data, s3BucketName, s3ObjectKeyPrefix string) error {
	
	// encryptionKey := os.Getenv("ENCRYPTION_KEY")
	// block, err := aes.NewCipher([]byte(encryptionKey))
	// if err != nil {
	// 		return err
	// }
	code, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		fmt.Println("Erreur lors de la création du QR code:", err)
		return err
	}

	png, err := code.PNG(256)
	if err != nil {
		fmt.Println("Erreur lors de la conversion du QR code en PNG:", err)
		return err
	}
	pngReader := bytes.NewReader(png)
	// Créez une session AWS
	if err != nil {
		fmt.Println("Erreur lors de la configuration AWS:", err)
		return err
	}
	_, err = basics.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String( s3BucketName),
		Key:    aws.String(s3ObjectKeyPrefix),
		Body:   pngReader,
	})
	if err != nil {
		log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
			png, s3BucketName, s3ObjectKeyPrefix, err)
	}

	if err != nil {
		fmt.Println("Erreur lors du téléchargement du fichier vers S3:", err)
		return err
	}

	fmt.Printf("QR code enregistré avec succès dans le compartiment S3 sous la clé : %s\n", s3ObjectKeyPrefix)
	return nil
}

func (basics BucketBasics) GenerateAndUploadMultipleQRCodes(data map[string]string, s3BucketName string) error {
	var wg sync.WaitGroup

	for dataValue, s3ObjectKeyPrefix := range data {
			wg.Add(1)
			go func(dataValue string, s3ObjectKeyPrefix string) {
					defer wg.Done()

					code, err := qrcode.New(dataValue, qrcode.Medium)
					if err != nil {
							fmt.Printf("Erreur lors de la création du QR code pour %s: %v\n", dataValue, err)
							return
					}

					png, err := code.PNG(256)
					if err != nil {
							fmt.Printf("Erreur lors de la conversion du QR code pour %s en PNG: %v\n", dataValue, err)
							return
					}

					s3ObjectKey := s3ObjectKeyPrefix + ".png"

					cfg, err := config.LoadDefaultConfig(context.TODO())
					if err != nil {
							fmt.Printf("Erreur lors de la configuration AWS pour le QR code pour %s: %v\n", dataValue, err)
							return
					}

					client := s3.NewFromConfig(cfg)

					input := &s3.PutObjectInput{
							Bucket: aws.String(s3BucketName),
							Key:    aws.String(s3ObjectKey),
							Body:   bytes.NewReader(png),
					}

					_, err = client.PutObject(context.TODO(), input)
					if err != nil {
							fmt.Printf("Erreur lors du téléchargement du fichier du QR code pour %s vers S3: %v\n", dataValue, err)
							return
					}

					fmt.Printf("QR code pour %s enregistré avec succès dans le compartiment S3 sous la clé : %s\n", dataValue, s3ObjectKey)
			}(dataValue, s3ObjectKeyPrefix)
	}

	wg.Wait()

	return nil
}