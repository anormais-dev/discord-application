package dto

type WSRequest struct {
	Type        string `json:"type"`
	AccessToken string `json:"accessToken"`
	Format      Format `json:"format"`
	UserID      string `json:"userId"`
	Enabled     bool   `json:"enabled"`
	GuildID     string `json:"guildId"`
}
