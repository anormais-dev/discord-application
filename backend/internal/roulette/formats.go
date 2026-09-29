package roulette

import "strconv"

// Format define quantos times e quantas pessoas por time.
type Format struct {
	Teams int `json:"teams"`
	Size  int `json:"size"`
}

// Slots é o total de vagas do formato.
func (f Format) Slots() int {
	return f.Teams * f.Size
}

// Label devolve o nome exibido, ex: "5x5".
func (f Format) Label() string {
	return strconv.Itoa(f.Size) + "x" + strconv.Itoa(f.Size)
}

// Formats são os presets que o admin pode escolher.
var Formats = []Format{
	{Teams: 2, Size: 1},
	{Teams: 2, Size: 2},
	{Teams: 2, Size: 3},
	{Teams: 2, Size: 4},
	{Teams: 2, Size: 5},
}

// DefaultFormat é o formato inicial de uma room nova.
var DefaultFormat = Formats[len(Formats)-1]

func validFormat(f Format) bool {
	for _, p := range Formats {
		if p == f {
			return true
		}
	}
	return false
}
