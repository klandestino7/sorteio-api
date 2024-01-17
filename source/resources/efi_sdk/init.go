package efi

import (
	"os"
	"sorteio-api/source/utils"
)

var Credentials = map[string]interface{}{}

func InitializePayment() {
	configsPayment := map[string]interface{}{
		"client_id":     os.Getenv("EFI_CLIENT_ID"),
		"client_secret": os.Getenv("EFI_CLIENT_SECRET"),
		"sandbox":       os.Getenv("EFI_SANDBOX") == "true",
		"timeout":       20,
		"CA":            os.Getenv("EFI_CA"),  // caminho da chave publica da gerencianet
		"Key":           os.Getenv("EFI_KEY"), // caminho da chave privada da sua conta Gerencianet
	}

	Credentials = configsPayment
	utils.DebugPrint("PAYMENT CONFIG LOADED")
}
