package discord

import "testing"

func TestAvatarURL(t *testing.T) {
	cases := []struct {
		name string
		user User
		want string
	}{
		{"estático", User{ID: "1", Avatar: "abc"}, "https://cdn.discordapp.com/avatars/1/abc.png?size=128"},
		{"animado", User{ID: "1", Avatar: "a_abc"}, "https://cdn.discordapp.com/avatars/1/a_abc.gif?size=128"},
		{"padrão", User{ID: "4194304"}, "https://cdn.discordapp.com/embed/avatars/1.png"},
	}
	for _, c := range cases {
		if got := c.user.AvatarURL(); got != c.want {
			t.Errorf("%s: %q, esperado %q", c.name, got, c.want)
		}
	}
}
