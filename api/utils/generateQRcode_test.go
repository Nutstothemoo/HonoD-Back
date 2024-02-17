package utils

import (
    "context"
    "testing"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/joho/godotenv"
)

func TestGenerateAndUploadQRCode(t *testing.T) {
    err := godotenv.Load("../../.env")
    if err != nil {
        t.Fatalf("Erreur lors du chargement du fichier .env : %v", err)
    }

    s3BucketName := "honod"
    s3ObjectKeyPrefix := "E-ticket/test.png"
    data := "Test QR Code"

    ctx := context.TODO()

    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        t.Fatalf("Erreur lors du chargement de la configuration par défaut AWS : %v", err)
    }
		s3Client := s3.NewFromConfig(cfg)

		// Initialisez l'objet BucketBasics avec votre client S3
		bucketBasics := BucketBasics{S3Client: s3Client}
    err = bucketBasics.GenerateAndUploadQRCode(data, s3BucketName, s3ObjectKeyPrefix)
    if err != nil {
        t.Errorf("Erreur lors de la génération et du téléchargement du QR code : %v", err)
    }

    t.Log("QR code généré et téléversé avec succès")
}