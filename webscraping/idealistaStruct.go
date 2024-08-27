package webscraping

type CookieHumana struct {
	Cookie string `json:"cookie"`
}

type Captcha struct {
	URL string `json:"url"`
}

type IdealistaURLs struct {
	Count               string `json:"count"`
	URL                 string `json:"url"`
	LocationID          string `json:"locationId"`
	Text                string `json:"text"`
	SuggestedLocationID int    `json:"suggestedLocationId"`
	Category            string `json:"category"`
	CategoryID          any    `json:"categoryId"`
	CategoryAlias       any    `json:"categoryAlias"`
	TotalResults        int    `json:"totalResults"`
	ZoneOfInterest      bool   `json:"zoneOfInterest"`
}

type PropiedadesIdealista struct {
	ID                  string
	Titulo              string   `json:"title"`
	TipoPropiedad       string   `json:"propertyType"`
	Comercio            string   `json:"businessType"`
	Municipio           string   `json:"municipality"`
	Ubicacion           string   `json:"location"`
	Zona                string   `json:"zone"`
	Direccion           string   `json:"address"`
	Valor               string   `json:"value"`
	Moneda              string   `json:"currency"`
	Habitaciones        string   `json:"rooms"`
	Banos               string   `json:"bathrooms"`
	Metros              string   `json:"size"`
	Descripcion         string   `json:"description"`
	Info                []string `json:"info"`
	Calefaccion         string   `json:"heating"`
	FechaConstruccion   string   `json:"constructionDate"`
	Planta              string   `json:"floor"`
	URL                 string   `json:"url"`
	PrecioMetroCuadrado string
	FechaExtraccion     string
}
