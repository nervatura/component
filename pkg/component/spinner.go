package component

import (
	ut "github.com/nervatura/component/pkg/util"
)

// The htmx AJAX request indicator.
type Spinner struct {
	Id string `json:"id"`
	// Modal appearance by default, other elements of the page are not available
	NoModal bool `json:"no_modal"`
}

/*
Based on the values, it will generate the html code of the [Spinner] or return with an error message.
*/
func (spn *Spinner) Render() (result string, err error) {
	spn.Id = ut.ToString(spn.Id, "spinner")
	tpl := `<div id="{{ .Id }}" class="htmx-indicator{{ if eq .NoModal false }} modal{{ end }}" >
	<div class="loading-middle" ><div class="loading">
	<div></div><div></div><div></div><div></div><div></div><div></div><div></div><div></div>
	</div></div></div>`

	return ut.TemplateBuilder("spinner", tpl, map[string]any{}, spn)
}
