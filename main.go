package main

import (
	"strings"

	"github.com/117503445/goutils"
	"github.com/imroc/req/v3"
	"github.com/rs/zerolog/log"
)

func Download(client *req.Client, url string) {
	fileName := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
	if goutils.FileExists(fileName) {
		return
	}

	log.Info().Str("url", url).Str("fileName", fileName).Msg("Downloading file")

	resp, err := client.R().
		SetOutputFile(fileName).
		Get(url)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to download file")
	}
	defer resp.Body.Close()
}

func main() {
	goutils.InitZeroLog()

	log.Info().Msg("Hello World")
	client := req.C()
	Download(client, "https://raw.githubusercontent.com/BrotatoMods/Brotato-Translations/refs/heads/main/translations-v1.1.6.0-dlc_1.csv")
	Download(client, "https://raw.githubusercontent.com/BrotatoMods/Brotato-Translations/refs/heads/main/translations-v1.1.6.0.csv")

}
