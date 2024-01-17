package utils

import (
	"fmt"
	"strings"
)

func NormalizeImageFromDB(pathName string) string {
	return fmt.Sprintf("https://storage.cloud.google.com/%s", pathName)
}

func NormalizeCPF(cpf string) string {
	t := strings.Replace(cpf, "-", "", -1)
	return strings.Replace(t, ".", "", -1)
}
