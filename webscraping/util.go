package webscraping

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var lock = &sync.Mutex{}

type FinalReport struct {
	ReporteName string
	ExcelPath   string
	Properties  []PropiedadesIdealista
	row         int
}

func (fr *FinalReport) SavePropiedad(hoja string, propiedad PropiedadesIdealista) {
	f, err := excelize.OpenFile(fr.ExcelPath)
	if err != nil {
		f = excelize.NewFile()
	}
	f.DeleteSheet("Sheet1")
	excelIndex := strconv.Itoa(fr.row)
	t := time.Now()
	sheet := hoja + " " + fmt.Sprintf((t.Format("2006-01-02")))
	f.SetCellValue(sheet, "A"+excelIndex, propiedad.Comercio)
	f.SetCellValue(sheet, "B"+excelIndex, propiedad.Municipio)
	f.SetCellValue(sheet, "C"+excelIndex, propiedad.Ubicacion)
	f.SetCellValue(sheet, "D"+excelIndex, propiedad.Direccion)
	f.SetCellValue(sheet, "E"+excelIndex, propiedad.Zona)
	f.SetCellValue(sheet, "F"+excelIndex, propiedad.Valor)
	f.SetCellValue(sheet, "G"+excelIndex, propiedad.Habitaciones)
	f.SetCellValue(sheet, "H"+excelIndex, propiedad.Banos)
	f.SetCellValue(sheet, "I"+excelIndex, propiedad.Metros)
	f.SetCellValue(sheet, "J"+excelIndex, propiedad.PrecioMetroCuadrado)
	f.SetCellValue(sheet, "K"+excelIndex, propiedad.Descripcion)
	f.AddComment(sheet, "K"+excelIndex, `{"author":"Excelize: ","text":"`+propiedad.Descripcion+` ."}`)
	f.SetCellValue(sheet, "L"+excelIndex, propiedad.URL)
	f.SetCellValue(sheet, "M"+excelIndex, propiedad.ID)
	f.SetCellValue(sheet, "N"+excelIndex, propiedad.FechaExtraccion)
	f.SetCellValue(sheet, "O"+excelIndex, propiedad.TipoPropiedad)
	if err := f.SaveAs(fr.ExcelPath); err != nil {
		log.Fatal(err)
	}
}

func (fr *FinalReport) CreateColumnaName(hoja string) {
	f, err := excelize.OpenFile(fr.ExcelPath)
	if err != nil {
		f = excelize.NewFile()
	}

	f.DeleteSheet("Sheet1")
	t := time.Now()
	sheet := hoja + " " + fmt.Sprintf((t.Format("2006-01-02")))
	index := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.SetCellValue(sheet, "A1", "Mercado")
	f.SetCellValue(sheet, "B1", "Municipio")
	f.SetCellValue(sheet, "C1", "Ubicación")
	f.SetCellValue(sheet, "D1", "Direccion")
	f.SetCellValue(sheet, "E1", "Zona")
	f.SetCellValue(sheet, "F1", "Valor")
	f.SetCellValue(sheet, "G1", "Habitaciones")
	f.SetCellValue(sheet, "H1", "Baños")
	f.SetCellValue(sheet, "I1", "Metrtos Cuadrados")
	f.SetCellValue(sheet, "J1", "Precio Metro Cuadrado")
	f.SetCellValue(sheet, "K1", "Descripcion")
	f.SetCellValue(sheet, "L1", "URL")
	f.SetCellValue(sheet, "M1", "ID")
	f.SetCellValue(sheet, "N1", "fecha de extraccion")
	f.SetCellValue(sheet, "O1", "Tipo Propiedad")

	if err := f.SaveAs(fr.ExcelPath); err != nil {
		log.Fatal(err)
	}

}

func (fr *FinalReport) CreateExcel(name string, address string) {

	f, err := excelize.OpenFile(name)
	if err != nil {
		f = excelize.NewFile()
	}

	f.DeleteSheet("Sheet1")
	t := time.Now()
	sheet := address + " " + fmt.Sprintf((t.Format("2006-01-02")))
	index := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.SetCellValue(sheet, "A1", "Mercado")
	f.SetCellValue(sheet, "B1", "Municipio")
	f.SetCellValue(sheet, "C1", "Ubicación")
	f.SetCellValue(sheet, "D1", "Zona")
	f.SetCellValue(sheet, "E1", "Valor")
	f.SetCellValue(sheet, "F1", "Habitaciones")
	f.SetCellValue(sheet, "G1", "Baños")
	f.SetCellValue(sheet, "H1", "Metrtos Cuadrados")
	f.SetCellValue(sheet, "I1", "Precio Metro Cuadrado")
	f.SetCellValue(sheet, "J1", "Descripcion")
	f.SetCellValue(sheet, "K1", "URL")
	f.SetCellValue(sheet, "L1", "ID")
	f.SetCellValue(sheet, "M1", "fecha de extraccion")
	f.SetCellValue(sheet, "N1", "Tipo Propiedad")

	for index, propiedad := range fr.Properties {
		excelIndex := strconv.Itoa(index + 2)
		f.SetCellValue(sheet, "A"+excelIndex, propiedad.Comercio)
		f.SetCellValue(sheet, "B"+excelIndex, propiedad.Municipio)
		f.SetCellValue(sheet, "C"+excelIndex, propiedad.Ubicacion)
		f.SetCellValue(sheet, "D"+excelIndex, propiedad.Zona)
		f.SetCellValue(sheet, "E"+excelIndex, propiedad.Valor)
		f.SetCellValue(sheet, "F"+excelIndex, propiedad.Habitaciones)
		f.SetCellValue(sheet, "G"+excelIndex, propiedad.Banos)
		f.SetCellValue(sheet, "H"+excelIndex, propiedad.Metros)
		f.SetCellValue(sheet, "I"+excelIndex, propiedad.PrecioMetroCuadrado)
		f.SetCellValue(sheet, "J"+excelIndex, propiedad.Descripcion)
		f.AddComment(sheet, "J"+excelIndex, `{"author":"Idealista Bot: ","text":"`+propiedad.Descripcion+`"}`)
		f.SetCellValue(sheet, "K"+excelIndex, propiedad.URL)
		f.SetCellValue(sheet, "L"+excelIndex, propiedad.ID)
		f.SetCellValue(sheet, "M"+excelIndex, propiedad.FechaExtraccion)
		f.SetCellValue(sheet, "N"+excelIndex, propiedad.TipoPropiedad)
	}

	if err := f.SaveAs(name); err != nil {
		log.Fatal(err)
	}
}

var singleReport *FinalReport

func getFinalReport() *FinalReport {
	if singleReport == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleReport == nil {
			singleReport = &FinalReport{}
		}
	}
	return singleReport
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func GenerarArchivoTexto(contenido string, nombreArchivo string) error {
	// Crea el archivo con el nombre especificado
	err := ioutil.WriteFile(nombreArchivo, []byte(contenido), 0644)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo: %v", err)
	}
	fmt.Printf("Se ha creado el archivo %s con éxito.\n", nombreArchivo)
	return nil
}
