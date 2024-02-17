package actions

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func TestUploadFile(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Erreur lors du chargement du fichier .env : %v", err)
	}

	bucketName := "honod"
	objectKey := "E-ticket/test.txt"
	fileName := "test.txt"

	// Récupérez les clés d'accès depuis les variables d'environnement
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Vérifiez si les clés d'accès sont présentes
	if accessKeyID == "" || secretAccessKey == "" {
		t.Fatalf("Les clés d'accès AWS ne sont pas configurées")
	}

	ctx := context.TODO()

	// Configurez manuellement la configuration avec les clés d'accès
	cfg, err := config.LoadDefaultConfig(ctx)
	fmt.Println(cfg)
	if err != nil {
		t.Fatalf("Erreur lors du chargement de la configuration par défaut AWS : %v", err)
	}
	// Initialisez le client S3 avec la configuration mise à jour
	s3Client := s3.NewFromConfig(cfg)

	// Initialisez l'objet BucketBasics avec votre client S3
	bucketBasics := BucketBasics{S3Client: s3Client}

	// Appelez la fonction UploadFile pour téléverser le fichier
	err = bucketBasics.UploadFile(bucketName, objectKey, fileName)
	if err != nil {
		t.Errorf("Erreur lors de l'upload du fichier : %v", err)
	}

	t.Log("Fichier téléversé avec succès")
}