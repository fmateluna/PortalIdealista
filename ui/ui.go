package ui

import (
	"botEs/webscraping"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	term "github.com/nsf/termbox-go"

	ui "github.com/VladimirMarkelov/clui"
)

type BotGrafico struct {
	brain                 webscraping.IdeaListaBot
	Return                ResultUI
	UniversoIdealistaURLs []webscraping.IdealistaURLs
}

type ResultUI struct {
	Comercios       []string
	TipoPropiedades []string
	Max             int
	ExcelPath       string
	Localidad       string
	Second          int64
	UrlFiltradas    []webscraping.IdealistaURLs
}

func (g *BotGrafico) createViewBot() {
	g.Return.Second = 1
	tipoComercio := ""
	cantidadPropiedades := 1200
	universoDatos := 0

	view := ui.AddWindow(5, 2, 180, 7, "WebScraping Idealista")
	view.SetPack(ui.Vertical)

	frmMain := ui.CreateFrame(view, 1, 1, ui.BorderNone, 1)
	frmMain.SetPack(ui.Vertical)

	frameTipoPropiedad := ui.CreateFrame(frmMain, 1, 1, ui.BorderThin, ui.AutoSize)
	frameTipoPropiedad.SetTitle("Tipo de propiedad")
	frameTipoPropiedad.SetPack(ui.Horizontal)

	frameTipoPropiedadA := ui.CreateFrame(frameTipoPropiedad, ui.Fixed, 1, ui.BorderNone, 1)
	frameTipoPropiedadA.SetPack(ui.Vertical)

	newdevelopment := ui.CreateCheckBox(frameTipoPropiedadA, ui.AutoSize, "Obra nueva", 0)
	newdevelopment.SetState(0)

	home := ui.CreateCheckBox(frameTipoPropiedadA, ui.AutoSize, "Viviendas", 0)
	home.SetState(1)

	room := ui.CreateCheckBox(frameTipoPropiedadA, ui.AutoSize, "Habitación", 0)
	room.SetState(0)

	frameTipoPropiedadB := ui.CreateFrame(frameTipoPropiedad, ui.Fixed, 1, ui.BorderNone, 1)
	frameTipoPropiedadB.SetPack(ui.Vertical)
	garage := ui.CreateCheckBox(frameTipoPropiedadB, ui.AutoSize, "Garajes", 0)
	garage.SetState(0)

	storageroom := ui.CreateCheckBox(frameTipoPropiedadB, ui.AutoSize, "Trasteros", 0)
	storageroom.SetState(0)

	office := ui.CreateCheckBox(frameTipoPropiedadB, ui.AutoSize, "Oficinas", 0)
	office.SetState(0)

	frameTipoPropiedadC := ui.CreateFrame(frameTipoPropiedad, ui.Fixed, 1, ui.BorderNone, 1)
	frameTipoPropiedadC.SetPack(ui.Vertical)
	warehouse := ui.CreateCheckBox(frameTipoPropiedadC, ui.AutoSize, "Locales o naves", 0)
	warehouse.SetState(0)
	land := ui.CreateCheckBox(frameTipoPropiedadC, ui.AutoSize, "Terrenos", 0)
	land.SetState(0)
	building := ui.CreateCheckBox(frameTipoPropiedadC, ui.AutoSize, "Edificios", 0)
	building.SetState(0)

	frameComercio := ui.CreateFrame(frmMain, 1, 1, ui.BorderThin, 1)

	frameComercio.SetTitle("Tipo de Comercio")
	frameComercio.SetPack(ui.Vertical)
	frameComercio.SetGaps(ui.KeepValue, 1)
	frameComercio.SetPaddings(1, 1)

	radioVenta := ui.CreateCheckBox(frameComercio, ui.AutoSize, "Comprar", 0)
	radioVenta.SetState(0)
	radioAlquiler := ui.CreateCheckBox(frameComercio, ui.AutoSize, "Alquiler", 0)
	radioAlquiler.SetState(0)

	frameUbicacion := ui.CreateFrame(frmMain, 1, 1, ui.BorderThin, 1)
	frameUbicacion.SetTitle("Atributos de Descarga")
	frameUbicacion.SetPack(ui.Vertical)
	frameUbicacion.SetAlign(ui.AlignLeft)

	frameUbicacion.SetSize(64, 2)
	ui.CreateLabel(frameUbicacion, ui.AutoSize, 1, "Localidad ", 1)
	editLocalidad := ui.CreateEditField(frameUbicacion, 64, "", ui.AutoSize)

	subLocalidad := ui.CreateCheckBox(frameUbicacion, ui.AutoSize, "Buscar en subLocalidades", 0)
	subLocalidad.SetState(0)

	ui.CreateLabel(frameUbicacion, ui.AutoSize, 1, "Cantidad a descargar", 1)
	editCantidad := ui.CreateEditField(frameUbicacion, ui.AutoSize, fmt.Sprintf("%d", cantidadPropiedades), 1)

	ui.CreateLabel(frameUbicacion, ui.AutoSize, 1, "Archivo Excel", 1)
	editExcelPath := ui.CreateEditField(frameUbicacion, ui.AutoSize, "", 1)
	now := time.Now()
	formattedDateTime := now.Format("20060102_150405")
	fileName := formattedDateTime + ".xlsx"
	editExcelPath.SetTitle(fileName)

	labelUniverso :=
		ui.CreateLabel(frameUbicacion, 1, 2, "Seleccione y luego marque con la tecla Espacio ", ui.Fixed)
	logBox := ui.CreateListBox(frameUbicacion, 1, 15, ui.Fixed)

	//logBox.OnSelectItem(func(ev ui.Event) {)

	ui.ActivateControl(view, editLocalidad)

	/*
		editLocalidad.OnKeyPress(func(key term.Key, ch rune) bool {
			if key == term.KeyCtrlM {
				actualizaListado()

				return true
			}
			return false
		})
	*/
	frameBotonera := ui.CreateFrame(frmMain, ui.ButtonBottom, 1, ui.BorderNone, 1)

	frameSegundos := ui.CreateFrame(frameBotonera, 1, 1, ui.BorderThin, 1)
	frameSegundos.SetTitle("Segundos entre descargas")
	frameSegundos.SetPack(ui.Vertical)
	frameSegundos.SetAlign(ui.AlignLeft)
	frameSegundos.SetSize(1, 1)
	editSegundos := ui.CreateEditField(frameSegundos, 64, "", ui.AutoSize)
	editSegundos.SetTitle("1")

	frameBotonera.SetPack(ui.Horizontal)
	frameBotonera.SetGaps(ui.KeepValue, 1)
	frameBotonera.SetPaddings(1, 1)

	btnActualiza := ui.CreateButton(frameBotonera, ui.AutoSize, ui.Fixed, "Mostrar Propiedades", ui.Fixed)
	btnActualiza.SetActiveBackColor(ui.ColorRed)

	btnExit := ui.CreateButton(frameBotonera, ui.AutoSize, ui.Fixed, "Cerrar", ui.Fixed)
	btnExit.OnClick(func(ev ui.Event) {
		g.Return.Max = 0
		go ui.Stop()
	})

	frmTask := ui.CreateFrame(frmMain, ui.AutoSize, 1, ui.BorderNone, ui.Fixed)

	pb := ui.CreateProgressBar(frmTask, 11, 1, 1)

	actualizaListado := func() {
		g.UniversoIdealistaURLs = []webscraping.IdealistaURLs{}
		g.Return.Comercios = []string{}
		g.Return.TipoPropiedades = []string{}

		g.Return.Localidad = editLocalidad.Title()
		if g.Return.Localidad != "" {

			if radioVenta.State() == 1 && radioAlquiler.State() == 1 {
				g.Return.Comercios = []string{"sale", "rent"}
				tipoComercio = "ambos"
			} else {
				if radioVenta.State() == 1 {
					g.Return.Comercios = []string{"sale"}
					tipoComercio = "compra"
				}
				if radioAlquiler.State() == 1 {
					g.Return.Comercios = []string{"rent"}
					tipoComercio = "alquiler"
				}
			}

			if newdevelopment.State() == 1 {
				g.Return.TipoPropiedades = append(g.Return.TipoPropiedades, "newdevelopment")
			}
			if home.State() == 1 {
				g.Return.TipoPropiedades = append(g.Return.TipoPropiedades, "home")
			}
			if room.State() == 1 {
				g.Return.TipoPropiedades = append(g.Return.TipoPropiedades, "room")
			}
			if garage.State() == 1 {
				g.Return.TipoPropiedades = append(g.Return.TipoPropiedades, "garage")
			}
			if storageroom.State() == 1 {
				g.Return.TipoPropiedades = append(g.Return.TipoPropiedades, "storageroom")
			}
			if office.State() == 1 {
				g.Return.TipoPropiedades = append(g.Return.TipoPropiedades, "office")
			}
			if warehouse.State() == 1 {
				g.Return.TipoPropiedades = append(g.Return.TipoPropiedades, "warehouse")
			}
			if land.State() == 1 {
				g.Return.TipoPropiedades = append(g.Return.TipoPropiedades, "land")
			}
			if building.State() == 1 {
				g.Return.TipoPropiedades = append(g.Return.TipoPropiedades, "building")
			}

			logBox.Clear()
			universoDatos = 0

			totales := make(map[string]int)

			for _, comercio := range g.Return.Comercios {

				for _, tipoPropiedad := range g.Return.TipoPropiedades {
					idealistaURLS := []webscraping.IdealistaURLs{}
					if subLocalidad.State() == 1 {
						idealistaURLS = g.brain.ShowLocalidadWithSub(g.Return.Localidad, comercio, tipoPropiedad)
					} else {
						idealistaURLS = g.brain.ShowLocalidad(g.Return.Localidad, comercio, tipoPropiedad)
					}
					for _, idealistaURL := range idealistaURLS {

						textCon := strings.ReplaceAll(idealistaURL.Text, "<b>", "")
						textCon = strings.ReplaceAll(textCon, "</b>", "")

						countStr := strings.ReplaceAll(idealistaURL.Count, ".", "")

						count, _ := strconv.Atoi(countStr)
						//count := idealistaURL.TotalResults
						universoDatos = universoDatos + count

						idealistaURL.Text = textCon

						g.UniversoIdealistaURLs = append(g.UniversoIdealistaURLs, idealistaURL)

						contador, ok := totales[textCon]
						if !ok {
							contador = count
						} else {
							contador = contador + count
						}

						totales[textCon] = contador

						if count > 0 {
							showProp := ""
							if tipoPropiedad == "newdevelopment" {
								showProp = "Obra nueva "
							}
							if tipoPropiedad == "home" {
								showProp = "Viviendas "
							}
							if tipoPropiedad == "room" {
								showProp = "Habitación "
							}
							if tipoPropiedad == "garage" {
								showProp = "Garajes "
							}
							if tipoPropiedad == "storageroom" {
								showProp = "Trasteros "
							}
							if tipoPropiedad == "office" {
								showProp = "Oficinas "
							}
							if tipoPropiedad == "warehouse" {
								showProp = "Locales o naves "
							}
							if tipoPropiedad == "land" {
								showProp = "Terrenos "
							}
							if tipoPropiedad == "building" {
								showProp = "Edificios "
							}
							showComercio := comercio
							if comercio == "sale" {
								showComercio = "Compra"
							}
							if comercio == "rent" {
								showComercio = "Alquiler"
							}
							logBox.AddItem(strconv.Itoa(logBox.ItemCount()+1) + " - [" + textCon + "] --> " + showProp + " + " + showComercio + " = " + strconv.Itoa(contador))

						}

						fmt.Sscan(editSegundos.Title(), &g.Return.Second)
					}
				}

			}

			logBox.Clear()

			keys := make([]string, 0, len(totales))
			for k := range totales {

				keys = append(keys, k)
			}

			sort.Slice(keys, func(i, j int) bool {
				return keys[i] < keys[j]
			})

			for _, key := range keys {
				if totales[key] > 0 {
					logBox.AddItem("[ ] " + key + " = " + strconv.Itoa(totales[key]) + " Propiedades")
				}
			}

			labelUniverso.SetTitle("Universo total de datos  " + strconv.Itoa(universoDatos) + " Propiedades en " + strconv.Itoa(logBox.ItemCount()) + " Localidades")
			pb.SetLimits(0, logBox.ItemCount())

		}
	}
	btnActualiza.OnClick(func(ev ui.Event) {
		actualizaListado()
	})
	btnWebScraping := ui.CreateButton(frameBotonera, ui.AutoSize, ui.Fixed, "WEBSCRAPING!", ui.Fixed)
	btnWebScraping.OnClick(func(ev ui.Event) {
		isOK := true

		if editLocalidad.Title() == "" {
			dialog := ui.CreateAlertDialog("Falta ingreso", "LOCALIDAD", "Aceptar")
			editLocalidad.Active()
			dialog.Result()
			isOK = false
		}

		if isOK && editCantidad.Title() == "" {
			dialog := ui.CreateAlertDialog("Falta ingreso", "CANTIDAD PROPIEDADES", "Aceptar")
			editCantidad.Active()
			dialog.Result()
			isOK = false
		}

		if isOK && radioVenta.State() != 1 && radioAlquiler.State() != 1 {
			dialog := ui.CreateAlertDialog("Falta ingreso", "TIPO COMERCIO", "Aceptar")
			dialog.Result()
			frameComercio.Active()
			isOK = false
		}

		if isOK && universoDatos == 0 {
			actualizaListado()
			if universoDatos == 0 {
				dialog := ui.CreateAlertDialog("ERROR", "Universo de datos vacio!", "Aceptar")
				logBox.Active()
				dialog.Result()
				isOK = false
			}
		}

		cantidadPropiedades, err := strconv.Atoi(editCantidad.Title())
		if err != nil || cantidadPropiedades <= 0 {

			dialog := ui.CreateAlertDialog("Falta ingreso", "La cantidad de propiedades debe ser valor numerico o mayor de 0", "Aceptar")
			dialog.Result()
			editCantidad.Active()
			isOK = false
		}

		if isOK && len(g.Return.UrlFiltradas) == 0 {
			dialog := ui.CreateAlertDialog("Falta ingreso", "Y la ubicacion?, no la selecciono!", "Aceptar")
			dialog.Result()
			logBox.Active()
			isOK = false

		}

		if isOK {
			//actualizaListado()
			botones := []string{"Descargar", "Cancelar"}

			propiedadesDescarga := fmt.Sprintf("Localidad %s ,  %s , Cantidad de registros = %d ", editLocalidad.Title(), tipoComercio, cantidadPropiedades)

			result := ui.CreateConfirmationDialog("Iniciar descarga", propiedadesDescarga, botones, 1)

			g.Return.Max = cantidadPropiedades

			result.OnClose(func() {
				if result.Result() == ui.DialogButton1 {
					//logBox.Clear()
					g.Return.ExcelPath = editExcelPath.Title()
					ui.Stop()
				}

			})

		}

	})

	logBox.OnKeyPress(func(key term.Key) bool {
		if key == 32 {

			//Guardar item seleccionado
			selectText := logBox.SelectedItemText()
			if strings.Contains(selectText, "[X]") {
				selectText = strings.ReplaceAll(selectText, "[X]", "[ ]")
			} else {
				selectText = strings.ReplaceAll(selectText, "[ ]", "[X]")
			}
			index := logBox.SelectedItem()

			//Cambiar el test [ ] [X]
			listado := []string{}
			for i := 0; i < logBox.ItemCount(); i++ {
				if i == index {
					listado = append(listado, selectText)
				} else {
					text, ok := logBox.Item(i)
					if ok {
						listado = append(listado, text)
					}
				}
			}
			logBox.Clear()
			g.Return.UrlFiltradas = []webscraping.IdealistaURLs{}

			for nuevoIndex, item := range listado {

				logBox.AddItem(item)
				if strings.Contains(item, "[X]") {
					itemSinX := strings.ReplaceAll(item, "[X] ", "")
					itemSinX = strings.Split(itemSinX, " = ")[0]
					for _, url := range g.UniversoIdealistaURLs {
						if strings.Contains(url.Text, itemSinX) {
							//panic(itemSinX + "   ---   " + url.Text)
							g.Return.UrlFiltradas = append(g.Return.UrlFiltradas, url)

						}
					}
				}
				if nuevoIndex == index {
					logBox.SelectItem(nuevoIndex)
				}
			}
			pb.SetLimits(0, logBox.ItemCount())
			pb.SetTitle("{{value}} / {{max}}")
			pb.SetValue(len(g.Return.UrlFiltradas))

			return true
		} else {
			return false
		}
	})

	editLocalidad.Active()

}

func (g *BotGrafico) MainLoop() {
	g.brain = webscraping.IdeaListaBot{}
	ui.InitLibrary()
	defer ui.DeinitLibrary()
	ui.SetCurrentTheme("basic")
	g.createViewBot()
	ui.MainLoop()
}
