package main

import (
	_ "embed"

	"github.com/117503445/goutils"
	"github.com/imroc/req/v3"
	"github.com/rs/zerolog/log"
)

func Download(client *req.Client, url string, filePath string) {
	if goutils.FileExists(filePath) {
		return
	}

	log.Info().Str("url", url).Str("filePath", filePath).Msg("Downloading file")

	resp, err := client.R().
		SetOutputFile(filePath).
		Get(url)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to download file")
	}
	defer resp.Body.Close()
}

func main() {
	client := req.C()
	Download(client, "https://raw.githubusercontent.com/BrotatoMods/Brotato-Translations/refs/heads/main/translations-v1.1.6.0-dlc_1.csv", "../../pkg/assets/translations-v1.1.6.0-dlc_1.csv")
	Download(client, "https://raw.githubusercontent.com/BrotatoMods/Brotato-Translations/refs/heads/main/translations-v1.1.6.0.csv", "../../pkg/assets/translations-v1.1.6.0.csv")
}
