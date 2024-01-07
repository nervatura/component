package base

type Spinner struct {
	Id      string `json:"id"`
	NoModal bool   `json:"no_modal"`
}

func (spn *Spinner) Render() (result string, err error) {
	spn.Id = ToString(spn.Id, "spinner")
	tpl := `<div id="{{ .Id }}" class="htmx-indicator{{ if eq .NoModal false }} modal{{ end }}" >
	<div class="loading-middle" ><div class="loading">
	<div></div><div></div><div></div><div></div><div></div><div></div><div></div><div></div>
	</div></div></div>`

	return TemplateBuilder("spinner", tpl, map[string]any{}, spn)
}
