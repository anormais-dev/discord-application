package dto

type MapsResponse struct {
	Maps []string `json:"maps"`
	// Artes de cada mapa pelo nome; pode faltar mapa se a API não tiver imagem.
	Images map[string]MapImagesResponse `json:"images"`
}

type MapImagesResponse struct {
	Card   string `json:"card"`
	Splash string `json:"splash"`
}
