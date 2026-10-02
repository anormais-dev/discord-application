package model

type Format struct {
	Teams int
	Size  int
}

func (f Format) Slots() int {
	return f.Teams * f.Size
}

var Formats = []Format{
	{Teams: 2, Size: 1},
	{Teams: 2, Size: 2},
	{Teams: 2, Size: 3},
	{Teams: 2, Size: 4},
	{Teams: 2, Size: 5},
}

var DefaultFormat = Formats[len(Formats)-1]
