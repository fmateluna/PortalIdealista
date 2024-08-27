package main

import (
	"botEs/ui"
	"botEs/webscraping"
	"fmt"
	"time"
)

func main() {
	principal()
}
func principal() {

	botUI := ui.BotGrafico{}
	botUI.MainLoop()

	if botUI.Return.Max > 0 {
		result := botUI.Return
		bot := webscraping.IdeaListaBot{}
		max := result.Max
		/*
			if len(botUI.Return.UrlFiltradas) < max {
				max = len(botUI.Return.UrlFiltradas)
			}
		*/
		now := time.Now().Format("2006-01-02" + " 15:04:05")
		bot.Cantidad = max
		bot.Second = botUI.Return.Second
		bot.SetExcelFile(result.ExcelPath)
		fmt.Println(now, "[IMPORTANTE!] No abra el archivo "+result.ExcelPath+" excel hasta que termine la aplicacion ")
		fmt.Println(now, "Segundos entre descargas :", bot.Second, " segundos, No debe interrumpir el proceso de descarga")
		time.Sleep(2 * time.Second)
		fmt.Println(now, " <-- Iniciando extraccion de datos...", max)
		time.Sleep(2 * time.Second)
		bot.Init()
		for index, url := range botUI.Return.UrlFiltradas {
			fmt.Println("Descargando desde -> ", url.URL, url.Text, "[", index+1, " / ", len(botUI.Return.UrlFiltradas), "]")
			bot.ObtenerPropiedades(url.URL)
			if len(bot.Propiedades) >= max {
				bot.Propiedades = []webscraping.PropiedadesIdealista{}
				bot.ProgressBarRun()
				if index+1 == len(botUI.Return.UrlFiltradas) {
					break
				}

			}

		}
		bot.ProgressBarStop()

		now = time.Now().Format("2006-01-02" + " 15:04:05")
		fmt.Println(now, "--> Finalizado ")

	}

}
