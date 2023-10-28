package utils

import(
	"fmt"
	"github.com/skip2/go-qrcode"
	"os"
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

func GenerateAndSaveQRCode( data string) error {

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
