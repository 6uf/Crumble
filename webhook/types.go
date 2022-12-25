package webhook

type Web struct {
	Content   string   `json:"content"`
	Embeds    []Embeds `json:"embeds"`
	Username  string   `json:"username"`
	AvatarURL string   `json:"avatar_url"`
}

type Embeds struct {
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Image       Image     `json:"image"`
	Thumbnail   Thumbnail `json:"thumbnail"`
	Fields      []Fields  `json:"fields"`
	Color       int       `json:"color"`
	Author      Author    `json:"author"`
	Footer      Footer    `json:"footer"`
}

type Fields struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Footer struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url"`
}

type Author struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
}

type Image struct {
	URL string `json:"url"`
}

type Thumbnail struct {
	URL string `json:"url"`
}
