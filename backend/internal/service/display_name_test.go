package service

import "testing"

func TestCleanName(t *testing.T) {
	cases := []struct {
		in   []string
		want string
	}{
		{[]string{"Ana ✨"}, "Ana"},
		{[]string{"🔥 João 🔥"}, "João"},
		{[]string{"x_Pedro_x"}, "x Pedro x"},
		{[]string{"Maria Clara"}, "Maria Clara"},
		{[]string{"João Silva Pereira"}, "João"},
		{[]string{"Bartholomeuzinho"}, "Bartholomeuz"},
		{[]string{"", "Bia"}, "Bia"},
		{[]string{"✨✨", "", "bruno.dev"}, "bruno dev"},
		{[]string{"★", ""}, "Jogador"},
	}
	for _, c := range cases {
		if got := CleanName(c.in...); got != c.want {
			t.Errorf("CleanName(%q) = %q, esperado %q", c.in, got, c.want)
		}
	}
}
