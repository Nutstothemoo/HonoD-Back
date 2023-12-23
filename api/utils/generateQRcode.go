package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sync"
	"ginapp/api/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/skip2/go-qrcode"
)

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

func GenerateAndSaveQRCodeOnDisk( data string) error {

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

func GenerateAndUploadQRCode( event models.Event, data, s3BucketName, s3ObjectKeyPrefix string) error {
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

	// Créez une session AWS
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Erreur lors de la configuration AWS:", err)
		return err
	}

	// Initialisez un client S3
	client := s3.NewFromConfig(cfg)

	// Créez un nom unique pour le fichier PNG
	s3ObjectKey := s3ObjectKeyPrefix + "qrcode.png"

	// Configurez les paramètres de téléchargement
	input := &s3.PutObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(s3ObjectKey),
		Body:   bytes.NewReader(png), // Chargez le contenu du PNG
	}

	// Téléchargez le fichier PNG vers S3
	_, err = client.PutObject(context.TODO(), input)
	if err != nil {
		fmt.Println("Erreur lors du téléchargement du fichier vers S3:", err)
		return err
	}

	fmt.Printf("QR code enregistré avec succès dans le compartiment S3 sous la clé : %s\n", s3ObjectKey)
	return nil
}

func GenerateAndUploadMultipleQRCodes( event models.Event, data, n int, s3BucketName, s3ObjectKeyPrefix string) error {
	// Créez une "wait group" pour attendre la fin de toutes les goroutines
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
			wg.Add(1) // Incrémentez le compteur de la "wait group" pour chaque goroutine

			go func(index int) {
					defer wg.Done() // Décrémentez le compteur de la "wait group" lorsque la goroutine se termine

					// Générez le contenu du QR code (utilisez un contenu unique pour chaque QR code)
					data := fmt.Sprintf("QR Code %d", index)
			// 				// Générez le contenu du QR code (incluez le nom de l'événement et le numéro de la place)
					// data := fmt.Sprintf("Event: %s\nDate: %s\nPlace: %d", event.Name, event.Date.Format("2006-01-02"), index+1)

			// // Signez numériquement le contenu avec la clé privée
			// signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, []byte(data))
					code, err := qrcode.New(data, qrcode.Medium)
					if err != nil {
							fmt.Printf("Erreur lors de la création du QR code %d: %v\n", index, err)
							return
					}

					png, err := code.PNG(256)
					if err != nil {
							fmt.Printf("Erreur lors de la conversion du QR code %d en PNG: %v\n", index, err)
							return
					}

					// Créez un nom unique pour le fichier PNG
					s3ObjectKey := s3ObjectKeyPrefix + fmt.Sprintf("qrcode%d.png", index)

					// Créez une session AWS
					cfg, err := config.LoadDefaultConfig(context.TODO())
					if err != nil {
							fmt.Printf("Erreur lors de la configuration AWS pour le QR code %d: %v\n", index, err)
							return
					}

					// Initialisez un client S3
					client := s3.NewFromConfig(cfg)

					// Configurez les paramètres de téléchargement
					input := &s3.PutObjectInput{
							Bucket: aws.String(s3BucketName),
							Key:    aws.String(s3ObjectKey),
							Body:   bytes.NewReader(png), // Chargez le contenu du PNG
					}

					// Téléchargez le fichier PNG vers S3
					_, err = client.PutObject(context.TODO(), input)
					if err != nil {
							fmt.Printf("Erreur lors du téléchargement du fichier du QR code %d vers S3: %v\n", index, err)
							return
					}

					fmt.Printf("QR code %d enregistré avec succès dans le compartiment S3 sous la clé : %s\n", index, s3ObjectKey)
			}(i)
	}

	// Attendez que toutes les goroutines se terminent
	wg.Wait()

	return nil
}
