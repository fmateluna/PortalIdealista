package webscraping

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/schollz/progressbar/v3"
)

type IdeaListaBot struct {
	Propiedades       []PropiedadesIdealista
	Cantidad          int
	CookieHumana      string
	cantidadPorPagina int
	paginaPorUrl      map[string]int
	bar               *progressbar.ProgressBar
	BarON             bool
	report            FinalReport
	Second            int64
}

func (bot *IdeaListaBot) SetExcelFile(excelpath string) {
	bot.report = FinalReport{}
	bot.report.ExcelPath = excelpath
}

func (bot *IdeaListaBot) SendToExcel(propiedad PropiedadesIdealista) {
	bot.report.row++
	bot.report.SavePropiedad("Detalle", propiedad)
}

func (bot *IdeaListaBot) IniciaExcel() {
	bot.report.row = 1
	bot.report.CreateColumnaName("Detalle")
}

func (bot *IdeaListaBot) ProgressBarRun() {
	bot.BarON = true
	bot.bar = progressbar.NewOptions(bot.Cantidad,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]█[reset]",
			SaucerHead:    "[yellow]▒[reset]",
			SaucerPadding: "░",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	go func() {
		for bot.BarON {
			bot.bar.Set(len(bot.Propiedades))
		}
	}()

}
func (bot *IdeaListaBot) ProgressBarStop() {
	bot.BarON = false
}

func (bot *IdeaListaBot) Init() {
	bot.cantidadPorPagina = 30
	bot.paginaPorUrl = make(map[string]int)
	bot.ProgressBarRun()
	bot.IniciaExcel()
}

func (bot *IdeaListaBot) RunFromUI(localidad, tipoPropiedad, comercio string) []IdealistaURLs {
	url := "https://www.idealista.com/es/locationsSuggest/" + comercio + "/" + tipoPropiedad + "?searchField=" + localidad
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authority", "www.idealista.com")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Referer", "https://www.idealista.com/")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"117\", \"Not;A=Brand\";v=\"8\", \"Chromium\";v=\"117\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyRes, _ := ioutil.ReadAll(resp.Body)
	urlDetails := []IdealistaURLs{}
	contenido := string(bodyRes)

	decoder := json.NewDecoder(strings.NewReader(contenido))
	decoder.Decode(&urlDetails)

	defer resp.Body.Close()
	return urlDetails
}

func (bot *IdeaListaBot) guardarCookie(cookie string) {
	//save struct to json file
	humano := CookieHumana{bot.CookieHumana}
	humano.Cookie = cookie
	file, _ := json.MarshalIndent(humano, "", " ")
	_ = ioutil.WriteFile(".\\humano.json", file, 0644)
}

func (bot *IdeaListaBot) leeCookieHumana() string {

	jsonPath := ".\\humano.json"

	jsonFile, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		panic("Error en archivo : " + jsonPath)
	}
	//fmt.Printf("%s\n", string(jsonFile))
	content := string(jsonFile)
	galletaHumana := CookieHumana{}

	err = json.Unmarshal([]byte(content), &galletaHumana)

	if err == nil {
		return galletaHumana.Cookie
	} else {
		panic("Error en archivo : " + jsonPath)
	}
}

func (bot *IdeaListaBot) ExistID(ID string) bool {
	for _, propiedad := range bot.Propiedades {
		if propiedad.ID == ID {
			return true
		}
	}
	return false
}

func (bot *IdeaListaBot) ObtenerPropiedades(url string) {
	//venta-viviendas/madrid-madrid/mapa

	path := strings.Split(url, "/")
	comercioTipoPropiedad := path[1]
	comercio := strings.Split(comercioTipoPropiedad, "-")[0]
	cantidadPropxPage := 0

	if comercio != "geo" {

		tipoPropiedad := strings.Split(comercioTipoPropiedad, "-")[1]
		reqUrl := ""

		pagina, ok := bot.paginaPorUrl[url]
		if !ok {
			pagina = 1
		}

		if pagina == 1 {
			reqUrl = "https://www.idealista.com" + strings.TrimSuffix(url, "/mapa") + "/?ordenado-por=fecha-publicacion-desc"
		} else {
			reqUrl = "https://www.idealista.com" + strings.TrimSuffix(url, "/mapa") + "/pagina-" + fmt.Sprint(pagina) + ".htm?ordenado-por=fecha-publicacion-desc"
		}

		//fmt.Println(reqUrl)

		req, err := http.NewRequest("GET", reqUrl, nil)

		req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
		req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
		req.Header.Set("Referer", "https://www.idealista.com/")
		req.Header.Set("Sec-Ch-Device-Memory", "8")
		req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"117\", \"Not;A=Brand\";v=\"8\", \"Chromium\";v=\"117\"")
		req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
		req.Header.Set("Sec-Ch-Ua-Model", "\"\"")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
		req.Header.Set("Cookie", bot.leeCookieHumana())
		time.Sleep(2 * time.Second)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic("Actualice el archivo humano.json")
		}
		defer resp.Body.Close()

		bodyRes, _ := ioutil.ReadAll(resp.Body)
		html := string(bodyRes)

		catpcha := Captcha{}

		decoder := json.NewDecoder(strings.NewReader(html))
		decoder.Decode(&catpcha)
		if err == nil && catpcha.URL != "" {
			fmt.Println("El sitio a detectado que no soy humano, asi que necesito que visites esta url")
			fmt.Println("######################################################\n\n\n\n")
			fmt.Println(catpcha.URL)
			//Crear archivo de error con URL
			fmt.Println("\n\n\n\n######################################################")
			fmt.Println("resuelvas el captcha y  luego actualices el archivo humano.json")

			bot.ProgressBarStop()

			fmt.Println("Ingrese valor cookie para continuar, no puede equivocarse")
			fmt.Println("Cuando alla guardado el archivo, presione ENTER para continuar..")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			//bot.guardarCookie(cookie)
			fmt.Println("Un momento por favor...")
			err := GenerarArchivoTexto(catpcha.URL, "Humano.url.log")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			time.Sleep(5 * time.Second)
			bot.ProgressBarRun()
			bot.ObtenerPropiedades(url)
		}

		p := strings.NewReader(html)
		doc, _ := goquery.NewDocumentFromReader(p)
		doc.Find(".item-info-container .item-link").Each(func(i int, el *goquery.Selection) {
			urlPropiedad, ok := el.Attr("href")
			if ok {
				cantidadPropxPage++

				path := strings.Split(url, "/")
				id := path[len(path)-2]

				if !bot.ExistID(id) {
					urlPropiedad = "https://www.idealista.com" + urlPropiedad
					propiedad := bot.savePropiedad(urlPropiedad)
					propiedad.TipoPropiedad = tipoPropiedad
					propiedad.Comercio = comercio
					ahora := time.Now()
					propiedad.FechaExtraccion = ahora.Format("2006-01-02 15:04:05")

					bot.Propiedades = append(bot.Propiedades, propiedad)
					bot.SendToExcel(propiedad)
					//fmt.Print(len(bot.Propiedades), ",")

					if len(bot.Propiedades) >= bot.Cantidad {
						bot.ProgressBarStop()
						return
					}
				}
			}
		})
		bot.paginaPorUrl[url] = bot.paginaPorUrl[url] + 1

		if bot.BarON {
			if cantidadPropxPage >= 30 {
				bot.ObtenerPropiedades(url)
			}
		}
	}

}

func (bot *IdeaListaBot) savePropiedad(url string) PropiedadesIdealista {

	propiedad := PropiedadesIdealista{}

	propiedad.URL = url

	path := strings.Split(url, "/")
	id := path[len(path)-2]
	propiedad.ID = id

	req, err := http.NewRequest("GET", url, nil)

	//req.Header.Set("Authority", "www.idealista.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", bot.leeCookieHumana())
	req.Header.Set("Sec-Ch-Device-Memory", "8")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"")
	req.Header.Set("Sec-Ch-Ua-Arch", "\"x86\"")
	req.Header.Set("Sec-Ch-Ua-Full-Version-List", "\"Chromium\";v=\"118.0.5993.70\", \"Google Chrome\";v=\"118.0.5993.70\", \"Not=A?Brand\";v=\"99.0.0.0\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Model", "\"\"")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")

	time.Sleep(time.Duration(bot.Second) * time.Second)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic("Actualice el archivo humano.json")
	}
	defer resp.Body.Close()

	bodyRes, _ := ioutil.ReadAll(resp.Body)
	html := string(bodyRes)

	p := strings.NewReader(html)
	doc, _ := goquery.NewDocumentFromReader(p)
	atributos := []string{}

	doc.Find("#headerMap > ul").Each(func(i int, el *goquery.Selection) {
		value := el.Text()

		value = strings.ReplaceAll(value, "\n\n", "\n")
		value = strings.ReplaceAll(value, "\n", "#")
		value = strings.ReplaceAll(value, "##", "#")

		direcciones := strings.Split(value, "#")
		propiedad.Direccion = direcciones[1] + " - " + direcciones[2]
		propiedad.Ubicacion = direcciones[len(direcciones)-4]
		propiedad.Zona = direcciones[len(direcciones)-3]
		propiedad.Municipio = direcciones[len(direcciones)-2]
	})

	doc.Find("#details > div.details-property > div.details-property-feature-one > div:nth-child(2) > ul > li").Each(func(i int, el *goquery.Selection) {
		value, err := el.Html()
		if err == nil {
			input := value
			regex := regexp.MustCompile("[^0-9,]+")
			cleaned := regex.ReplaceAllString(input, "")
			if i == 0 {
				propiedad.Metros = cleaned
			}

			if i == 1 {
				propiedad.Habitaciones = cleaned
			}

			if i == 2 {
				propiedad.Banos = cleaned
			}

			atributos = append(atributos, value)
			//fmt.Println(i, value)
		}
	})

	doc.Find("#fixed-toolbar > div > div.main-info > p.info-data.txt-big > span.price-container .price").Each(func(i int, el *goquery.Selection) {
		value, err := el.Html()
		if err == nil {
			regex := regexp.MustCompile("[^0-9,]+")
			cleaned := regex.ReplaceAllString(value, "")
			propiedad.Valor = cleaned
		}
	})

	propiedad.Planta = bot.capturaTextBySelector(html, "#fixed-toolbar > div > div.main-info > p.info-data.txt-big > span:nth-child(4)")
	propiedad.Descripcion = bot.capturaTextBySelector(html, "#main > div > main > section.detail-info.ide-box-detail-first-picture.ide-box-detail--reset.overlay-box > section > div.commentsContainer > div.comment > div > p")
	propiedad.Titulo = bot.capturaTextBySelector(html, "#main > div > main > section.detail-info.ide-box-detail-first-picture.ide-box-detail--reset.overlay-box > section > div.commentsContainer > div.comment > div > p")
	propiedad.PrecioMetroCuadrado = bot.capturaTextBySelector(html, "#mortgages > div.ide-box-detail.overlay-box.auction-box-detail > div > article > section > p.flex-feature.squaredmeterprice > span:nth-child(2)")

	detalleDireccion := bot.capturaTextBySelector(html, "#fixed-toolbar > div > div.main-info > p.sticky-bar-detail-heading.txt-body > span")

	propiedad.Direccion = propiedad.Direccion + "," + detalleDireccion

	if propiedad.Metros == "" {
		propiedad.Metros = bot.capturaTextBySelector(html, "#fixed-toolbar > div > div.main-info > p.info-data.txt-big > span:nth-child(2) > span")
	}
	regex := regexp.MustCompile("[^0-9,]+")
	cleaned := regex.ReplaceAllString(propiedad.PrecioMetroCuadrado, "")
	propiedad.PrecioMetroCuadrado = cleaned

	//propiedad.Direccion = bot.capturaTextBySelector(html, "#headerMap > ul")

	total := len(atributos)
	if total > 0 {
		if total-1 > 0 {
			propiedad.Calefaccion = atributos[total-1]
		}
		if total-2 > 0 {
			propiedad.FechaConstruccion = atributos[total-2]
		}
		propiedad.Info = atributos
	}
	return propiedad
}

func (bot *IdeaListaBot) capturaTextBySelector(html string, selector string) string {
	value := ""
	p := strings.NewReader(html)
	doc, _ := goquery.NewDocumentFromReader(p)

	doc.Find(selector).Each(func(i int, el *goquery.Selection) {
		value = el.Text()

	})

	re := regexp.MustCompile(`\r?\n`)
	result := re.ReplaceAllString(value, " ")

	return result
}

func (bot *IdeaListaBot) ShowLocalidad(localidad, comercio, tipoPropiedad string) []IdealistaURLs {

	result := []IdealistaURLs{}

	req, err := http.NewRequest("GET", "https://www.idealista.com/es/locationsSuggest/"+comercio+"/"+tipoPropiedad+"?searchField="+localidad, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", "www.idealista.com")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Referer", "https://www.idealista.com/")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"117\", \"Not;A=Brand\";v=\"8\", \"Chromium\";v=\"117\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := http.DefaultClient.Do(req)
	bodyRes, _ := ioutil.ReadAll(resp.Body)
	urlDetails := []IdealistaURLs{}
	contenido := string(bodyRes)

	decoder := json.NewDecoder(strings.NewReader(contenido))
	decoder.Decode(&urlDetails)
	if err != nil {
		// handle err
	}

	for _, url := range urlDetails {
		//fmt.Println(url.Text, " = ", url.Count)
		//bot.ObtenerPropiedades(url.URL)
		result = append(result, url)

	}

	defer resp.Body.Close()

	return result
}

func (bot *IdeaListaBot) ShowLocalidadWithSub(localidad, comercio, tipoPropiedad string) []IdealistaURLs {

	result := []IdealistaURLs{}

	req, err := http.NewRequest("GET", "https://www.idealista.com/es/locationsSuggest/"+comercio+"/"+tipoPropiedad+"?searchField="+localidad, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", "www.idealista.com")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Referer", "https://www.idealista.com/")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"117\", \"Not;A=Brand\";v=\"8\", \"Chromium\";v=\"117\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := http.DefaultClient.Do(req)
	bodyRes, _ := ioutil.ReadAll(resp.Body)
	urlDetails := []IdealistaURLs{}
	contenido := string(bodyRes)

	decoder := json.NewDecoder(strings.NewReader(contenido))
	decoder.Decode(&urlDetails)
	if err != nil {
		// handle err
	}

	for _, url := range urlDetails {
		//fmt.Println(url.Text, " = ", url.Count)
		//bot.ObtenerPropiedades(url.URL)
		subLocals := bot.ShowSubLocalidad(url.URL)
		//fmt.Println(url.URL)
		for _, subLocal := range subLocals {
			//fmt.Println("\t", subLocal.URL)
			if bot.NotExistURLInList(result, subLocal.URL) {
				subLocal.Text = url.Text + " -> " + subLocal.Text
				result = append(result, subLocal)
			}
		}
	}

	defer resp.Body.Close()

	return result
}

func (bot *IdeaListaBot) NotExistURLInList(urls []IdealistaURLs, url string) bool {
	for _, urlObj := range urls {
		if urlObj.URL == url {
			return false
		}
	}
	return true
}

func (bot *IdeaListaBot) ShowSubLocalidad(localidadURL string) []IdealistaURLs {
	subURLs := []IdealistaURLs{}
	req, err := http.NewRequest("GET", "https://www.idealista.com"+localidadURL, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", "www.idealista.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", "didomi_token=eyJ1c2VyX2lkIjoiMThhZGQ4NTktMDcxNi02YTkzLTg3MWYtNGFmMjEzMThiMmU0IiwiY3JlYXRlZCI6IjIwMjMtMDktMjhUMjA6NDA6MzcuNDY3WiIsInVwZGF0ZWQiOiIyMDIzLTA5LTI4VDIwOjQwOjM3LjQ2N1oiLCJ2ZW5kb3JzIjp7ImVuYWJsZWQiOlsiZ29vZ2xlIiwiYzpsaW5rZWRpbi1tYXJrZXRpbmctc29sdXRpb25zIiwiYzptaXhwYW5lbCIsImM6YWJ0YXN0eS1MTGtFQ0NqOCIsImM6aG90amFyIiwiYzp5YW5kZXhtZXRyaWNzIiwiYzpiZWFtZXItSDd0cjdIaXgiLCJjOmFwcHNmbHllci1HVVZQTHBZWSIsImM6dGVhbGl1bWNvLURWRENkOFpQIiwiYzp0aWt0b2stS1pBVVFMWjkiLCJjOmlkZWFsaXN0YS1MenRCZXFFMyIsImM6aWRlYWxpc3RhLWZlUkVqZTJjIl19LCJwdXJwb3NlcyI6eyJlbmFibGVkIjpbImFuYWx5dGljcy1IcEJKcnJLNyIsImdlb2xvY2F0aW9uX2RhdGEiLCJkZXZpY2VfY2hhcmFjdGVyaXN0aWNzIl19LCJ2ZXJzaW9uIjoyLCJhYyI6IkFVYUFFQUZrQW93QS5BRm1BQ0FGayJ9; euconsent-v2=CPyz5QAPyz5QAAHABBENDYCsAP_AAE7AAAAAF5wBgAIAAqABaAFsAUgC8wAAAEBoAMAARBQJQAYAAiCgUgAwABEFAhABgACIKA6ADAAEQUAkAGAAIgoDIAMAARBQFQAYAAiCgAAA.f_gACdgAAAAA; _gcl_au=1.1.1162614620.1695933639; _fbp=fb.1.1695933638884.1973626187; _tt_enable_cookie=1; _ttp=9p1GezUNwiqWu1t_Xwa12J9X2Wy; _hjSessionUser_250321=eyJpZCI6IjZiOTcxY2RlLTc4OGQtNTY0MC1hZjhiLWU4MmY4MTk1MjhhOSIsImNyZWF0ZWQiOjE2OTU5MzM2Mzg5MjEsImV4aXN0aW5nIjp0cnVlfQ==; __rtbh.lid=%7B%22eventType%22%3A%22lid%22%2C%22id%22%3A%22LjezSjxkBExznpg13Ikg%22%7D; utag_main_v_id=018add858ebf009dae74dc0d8c880506f001a06700bd0; _pprv=eyJjb25zZW50Ijp7IjAiOnsibW9kZSI6Im9wdC1pbiJ9LCIxIjp7Im1vZGUiOiJvcHQtaW4ifSwiMiI6eyJtb2RlIjoib3B0LWluIn0sIjMiOnsibW9kZSI6Im9wdC1pbiJ9LCI0Ijp7Im1vZGUiOiJvcHQtaW4ifSwiNSI6eyJtb2RlIjoib3B0LWluIn0sIjYiOnsibW9kZSI6Im9wdC1pbiJ9LCI3Ijp7Im1vZGUiOiJvcHQtaW4ifX0sInB1cnBvc2VzIjpudWxsLCJfdCI6Im0zdThjdXQyfGxvNXRmZGgyIn0%3D; _pprv=eyJjb25zZW50Ijp7IjAiOnsibW9kZSI6Im9wdC1pbiJ9LCIxIjp7Im1vZGUiOiJvcHQtaW4ifSwiMiI6eyJtb2RlIjoib3B0LWluIn0sIjMiOnsibW9kZSI6Im9wdC1pbiJ9LCI0Ijp7Im1vZGUiOiJvcHQtaW4ifSwiNSI6eyJtb2RlIjoib3B0LWluIn0sIjYiOnsibW9kZSI6Im9wdC1pbiJ9LCI3Ijp7Im1vZGUiOiJvcHQtaW4ifX0sInB1cnBvc2VzIjpudWxsLCJfdCI6Im0zdThjdXQyfGxvNXRmZGgyIn0%3D; _pcid=%7B%22browserId%22%3A%22ln3n3i5mh36mz4sl%22%2C%22_t%22%3A%22m3u8cut3%7Clo5tfdh3%22%7D; _pcid=%7B%22browserId%22%3A%22ln3n3i5mh36mz4sl%22%2C%22_t%22%3A%22m3u8cut3%7Clo5tfdh3%22%7D; _pcus=eyJ1c2VyU2VnbWVudHMiOm51bGwsIl90IjoibTN1OGN1dDR8bG81dGZkaDQifQ%3D%3D; _pctx=%7Bu%7DN4IgrgzgpgThIC4B2YA2qA05owMoBcBDfSREQpAeyRCwgEt8oBJAE0RXSwH18yBbAMxgAHAGMw%2BAKwAfVJSn4AZqwAWUkAF8gA; _pctx=%7Bu%7DN4IgrgzgpgThIC4B2YA2qA05owMoBcBDfSREQpAeyRCwgEt8oBJAE0RXSwH18yBbAMxgAHAGMw%2BAKwAfVJSn4AZqwAWUkAF8gA; smc=\"{}\"; afUserId=41ec503b-51bc-49df-9745-a7cb41ccdeb1-p; utag_main__prevCompleteClickName=010-idealista/home > portal > home > button-search; _clck=fwookc|2|fgc|0|1366; _uetvid=489ea4f05e3f11eeb2fe3d657b6edc67; askToSaveAlertPopUp=true; userUUID=d7037ff9-d8af-408b-bd2e-ed190201f424; _hjHasCachedUserAttributes=true; contact0f778180-6163-4171-a5a0-28aa9bf1793c=\"{\\'maxNumberContactsAllow\\':10}\"; SESSION=7dde9dca2e01ec7c~0d277870-610d-4705-ab63-612ab23d9a4e; utag_main__sn=31; utag_main_ses_id=1699364206881%3Bexp-session; utag_main__prevVtUrl=https://www.idealista.com/%3Bexp-1699367807419; utag_main__prevVtUrlReferrer=%3Bexp-1699367807419; utag_main__prevVtSource=Direct traffic%3Bexp-1699367807419; utag_main__prevVtCampaignName=organicWeb%3Bexp-1699367807419; utag_main__prevVtCampaignCode=undefined%3Bexp-1699367807419; utag_main__prevVtCampaignLinkName=undefined%3Bexp-1699367807419; utag_main__prevVtRecipientId=undefined%3Bexp-1699367807419; _hjIncludedInSessionSample_250321=1; _hjSession_250321=eyJpZCI6ImQ5ZGIwOWFjLTAyZGItNDBhNC04ZTY5LTkwMTg3ZmEzNzFjOCIsImNyZWF0ZWQiOjE2OTkzNjQyMDc4ODIsImluU2FtcGxlIjp0cnVlLCJzZXNzaW9uaXplckJldGFFbmFibGVkIjpmYWxzZX0=; _hjAbsoluteSessionInProgress=0; utag_main__ss=0%3Bexp-session; ABTastySession=mrasn=&lp=https%253A%252F%252Fwww.idealista.com%252Fventa-viviendas%252Fmadrid-madrid%252Fmapa; cookieSearch-1=\"/venta-viviendas/madrid-madrid/:1699364508142\"; utag_main__se=6%3Bexp-session; utag_main__st=1699366309584%3Bexp-session; utag_main__pn=6%3Bexp-session; ABTasty=uid=gm0ca121h9jpqvmb&fst=1695933638703&pst=1699232142717&cst=1699364295586&ns=22&pvt=84&pvis=3&th=1086455.-1.8.3.2.1.1699232142731.1699364511814.1.22; utag_main__prevCompletePageName=005-idealista/portal > portal > viewGreyMap%3Bexp-1699368111829; utag_main__prevLevel2=005-idealista/portal%3Bexp-1699368111829; cto_bundle=yqRhuF9kN00xNjlyTTh2MHlXZ0p1QlRZZUNGWnBkTnUlMkZBbXdCcHh4OUJtVTVMMjRTOVlBSkQ1VmlVY3MlMkZLZiUyRlNzVVBab1R4RVFNTGRIelExQ0VZWXlWTDNwbXAlMkJpV0NXdUgxWjd3aVAzOVVpaXBWZTBUVUxzbmEyUW1UNnVtSWtCNkNEcFJSUGNVbzRGaGYlMkZJNWVZVGQ1TVdhb2RYSkplbU5qQnl2UXdLR1JhRWJBWmptR1dsWDJMRzZ1em9rbVBLcUhHYiUyQnZldEVPOUFMNnhEJTJGY3lVdzk4WnclM0QlM0Q; pbw=%24b%3d16999%3b%24o%3d11100; vs=33114=5706101; TestIfCookie=ok; TestIfCookieP=ok; pid=3687303593924752642; sasd2=q=%24qc%3D1500031917%3B%24ql%3DMedium%3B%24qt%3D11_2108_390030t%3B%24dma%3D0&c=1&l=-1269340260&lo=-935626210&lt=638349613145222179&o=1; sasd=%24qc%3D1500031917%3B%24ql%3DMedium%3B%24qt%3D11_2108_390030t%3B%24dma%3D0; datadome=la8PA2vtDAlCnMW33P7Qhuul2qcqYlZ1XtmhYotlSDA9YjSEeKGku4cJR~ktmkV5aWpomGvzDz1Sm7_9UbYuJxXdyaREWafeETMdIv7BKYbDPONK6iyoFSMMTrYIQg_g")
	req.Header.Set("Sec-Ch-Device-Memory", "8")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Arch", "\"x86\"")
	req.Header.Set("Sec-Ch-Ua-Full-Version-List", "\"Google Chrome\";v=\"119.0.6045.105\", \"Chromium\";v=\"119.0.6045.105\", \"Not?A_Brand\";v=\"24.0.0.0\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Model", "\"\"")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyRes, _ := ioutil.ReadAll(resp.Body)
	html := string(bodyRes)

	p := strings.NewReader(html)
	doc, _ := goquery.NewDocumentFromReader(p)

	doc.Find("#sublocations > li").Each(func(i int, el *goquery.Selection) {
		subHtml, ok := el.Html()
		if ok == nil {
			subULR := IdealistaURLs{}
			p := strings.NewReader(subHtml)
			subDoc, _ := goquery.NewDocumentFromReader(p)

			subDoc.Find("a").Each(func(i int, el *goquery.Selection) {
				href, ok := el.Attr("href")
				nombre := el.Text()
				if ok {
					subULR.Text = nombre
					subULR.URL = href
				}
			})

			subDoc.Find(".subdued").Each(func(i int, el *goquery.Selection) {
				nombre := el.Text()
				subULR.Count = nombre
			})

			subURLs = append(subURLs, subULR)
		}
	})

	return subURLs
}
