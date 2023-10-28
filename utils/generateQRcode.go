package utils

import(
	"context"
	"fmt"
	"github.com/skip2/go-qrcode"
	"os"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode/qr/gf256"
)

func savePNG(filename string, data []byte) error {
	
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

	err = savePNG("qrcode.png", png)
	if err != nil {
			fmt.Println("Erreur lors de l'enregistrement du fichier PNG:", err)
			return err
	}

	fmt.Println("QR code généré avec succès et enregistré sous 'qrcode.png'")
	return nil
}

func GenerateAndUploadQRCode(data string, s3BucketName, s3ObjectKey string) error {

	// Générer le QR code
	qrCode, err := qr.Encode(data, qr.M, qr.Auto)
	if err != nil {
		return err
	}

	qrCode, _ = barcode.Scale(qrCode, 256, 256)
	pngBytes := qrCode.PNG()

	// Configurer la session AWS
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	// Créer un client S3
	client := s3.NewFromConfig(cfg)

	// Préparer les paramètres de chargement S3
	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(s3ObjectKey),
		Body:   bytes.NewReader(pngBytes),
		ACL:    s3.BucketCannedACLPrivate, // ACL privé, modifiez selon vos besoins
	}

	// Charger le QR code dans S3

	_, err = client.PutObject(context.TODO(), uploadInput)
	if err != nil {
		return err
	}

	fmt.Println("QR code généré avec succès et enregistré dans S3 sous le chemin :", s3ObjectKey)
	return nil
}