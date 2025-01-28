package component

import (
	"html/template"
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [Upload] constants
const (
	ComponentTypeUpload = "upload"

	UploadEventUpload         = "upload"
	UploadDefaultPlaceholder  = "Choose file to upload"
	UploadDefaultToastMessage = "Successful file upload"
)

// Creates a file upload control
type Upload struct {
	BaseComponent
	// Specifies a filter for what file types the user can pick from the file input dialog box
	// Valid values: [file_extension], "audio/*", "video/*", "image/*", [media_type]
	Accept string `json:"accept"`
	// Specifies a short hint that describes the expected value of the input element
	Placeholder string `json:"placeholder"`
	// The text of the Toast message after a successful file upload
	ToastMessage string `json:"toast_message"`
	// Specifies that the input should be disabled
	Disabled bool `json:"disabled"`
	// Specifies the maximum number of characters allowed in the input element
	MaxLength int64 `json:"max_length"`
	// Full width cell (100%)
	Full bool `json:"full"`
}

/*
Returns all properties of the [Upload]
*/
func (upl *Upload) Properties() ut.IM {
	return ut.MergeIM(
		upl.BaseComponent.Properties(),
		ut.IM{
			"accept":        upl.Accept,
			"placeholder":   upl.Placeholder,
			"toast_message": upl.ToastMessage,
			"disabled":      upl.Disabled,
			"max_length":    upl.MaxLength,
			"full":          upl.Full,
		})
}

/*
Returns the value of the property of the [Upload] with the specified name.
*/
func (upl *Upload) GetProperty(propName string) interface{} {
	return upl.Properties()[propName]
}

/*
Setting a property of the [Upload] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (upl *Upload) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"accept": func() interface{} {
			upl.Accept = ut.ToString(propValue, "")
			return upl.Accept
		},
		"placeholder": func() interface{} {
			upl.Placeholder = ut.ToString(propValue, UploadDefaultPlaceholder)
			return upl.Placeholder
		},
		"disabled": func() interface{} {
			upl.Disabled = ut.ToBoolean(propValue, false)
			return upl.Disabled
		},
		"max_length": func() interface{} {
			upl.MaxLength = ut.ToInteger(propValue, 0)
			return upl.MaxLength
		},
		"full": func() interface{} {
			upl.Full = ut.ToBoolean(propValue, false)
			return upl.Full
		},
		"toast_message": func() interface{} {
			upl.ToastMessage = ut.ToString(propValue, UploadDefaultToastMessage)
			return upl.ToastMessage
		},
	}
	if _, found := pm[propName]; found {
		return upl.SetRequestValue(propName, pm[propName](), []string{})
	}
	if upl.BaseComponent.GetProperty(propName) != nil {
		return upl.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

/*
If the OnResponse function of the [Upload] is implemented, the function calls it after the [TriggerEvent]
is processed, otherwise the function's return [ResponseEvent] is the processed [TriggerEvent].
*/
func (upl *Upload) OnRequest(te TriggerEvent) (re ResponseEvent) {
	evt := ResponseEvent{
		Trigger: &Toast{
			Type:  ToastTypeSuccess,
			Value: upl.ToastMessage,
		},
		TriggerName: upl.Name,
		Name:        UploadEventUpload,
		Header: ut.SM{
			HeaderRetarget: "#toast-msg",
			HeaderReswap:   SwapInnerHTML,
		},
	}
	if upl.OnResponse != nil {
		return upl.OnResponse(evt)
	}
	return evt
}

func (upl *Upload) getComponent(name string) (html template.HTML, err error) {
	ccMap := map[string]func() ClientComponent{
		"submit": func() ClientComponent {
			return &Button{
				BaseComponent: BaseComponent{
					Id: upl.Id + "_" + name, Name: name,
					Style: ut.SM{"padding": "8px"},
				},
				ButtonStyle: ButtonStyleDefault,
				Type:        ButtonTypeSubmit,
				Icon:        "Upload",
			}
		},
	}
	cc := ccMap[name]()
	html, err = cc.Render()
	return html, err
}

/*
Based on the values, it will generate the html code of the [Upload] or return with an error message.
*/
func (upl *Upload) Render() (html template.HTML, err error) {
	upl.InitProps(upl)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(upl.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(upl.Class, " ")
		},
		"uploadComponent": func(name string) (template.HTML, error) {
			return upl.getComponent(name)
		},
	}
	tpl := `
	<form id="{{ .Id }}" name="{{ .Name }}" method="POST" enctype="multipart/form-data" 
	{{ if eq .Disabled false }}{{ if ne .EventURL "" }} hx-post="{{ .EventURL }}" hx-target="{{ .Target }}" {{ if ne .Sync "none" }} hx-sync="{{ .Sync }}"{{ end }} hx-swap="{{ .Swap }}"{{ end }}
	{{ if ne .Indicator "none" }} hx-indicator="#{{ .Indicator }}"{{ end }}{{ end }}>
	<div class="upload{{ if .Full }} full{{ end }} {{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}>
	  <div class="row{{ if .Full }} full{{ end }}"><div class="cell">
		<label id="{{ .Id }}_label" for="{{ .Id }}_input" 
		class="{{ if .Disabled }}upload-disabled{{else}}link{{end}}" >{{ .Placeholder }}</label>
		</div><div class="cell" style="width: 1px;">
	  <input id="{{ .Id }}_input" type="file" {{ if .Disabled }}disabled{{ end }} name="file"></input>
		</div><div id="{{ .Id }}_submit_cell" class="hide">
	  {{ uploadComponent "submit" }}
		</div></div>
		<div class="row full"><div id="{{ .Id }}_progress_cell" class="hide">
	  <progress id="{{ .Id }}_progress" value="0" max="100"></progress>
		</div></div>
	</div>
	</form>
  <script>
	htmx.on('#{{ .Id }}_input', 'change', function(evt) {
		var maxlen = {{ .MaxLength }};
    if (evt.target.files[0]) {
			var fName = evt.target.files[0].name;
			if ((maxlen > 0) && (fName.length > maxlen)) {
				fName = fName.substring(0,maxlen)+'...';
			} 
			htmx.find('#{{ .Id }}_label').textContent = fName;
			htmx.find('#{{ .Id }}_submit_cell').className = 'cell';
		} else {
			htmx.find('#{{ .Id }}_label').textContent = '{{ .Placeholder }}';
			htmx.find('#{{ .Id }}_submit_cell').className = 'hide';
		}
  });
  htmx.on('#{{ .Id }}', 'htmx:xhr:progress', function(evt) {
		//console.log('Here ' + evt.detail.loaded/evt.detail.total * 100)
    htmx.find('#{{ .Id }}_progress').setAttribute('value', evt.detail.loaded/evt.detail.total * 100)
  });
	htmx.on('#{{ .Id }}', 'htmx:xhr:loadstart', function(evt) {
		htmx.find('#{{ .Id }}_progress_cell').className = 'cell';
  });
	htmx.on('#{{ .Id }}', 'htmx:xhr:abort', function(evt) {
		htmx.find('#{{ .Id }}_progress_cell').className = 'hide';
  });
	htmx.on('#{{ .Id }}', 'htmx:xhr:loadend', function(evt) {
		htmx.find('#{{ .Id }}_progress_cell').className = 'hide';
		htmx.find('#{{ .Id }}_submit_cell').className = 'hide';
		htmx.find('#{{ .Id }}_input').setAttribute('value', null);
		htmx.find('#{{ .Id }}_label').textContent = '{{ .Placeholder }}';
  });
  </script>`

	if html, err = ut.TemplateBuilder("upload", tpl, funcMap, upl); err == nil && upl.EventURL != "" {
		upl.SetProperty("request_map", upl)
	}
	return html, nil
}

// [Upload] test and demo data
func TestUpload(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeUpload,
			Component: &Upload{
				BaseComponent: BaseComponent{
					Id:           id + "_upload_default",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Accept:   "",
				Full:     false,
				Disabled: false,
			}},
		{
			Label:         "Disabled, full",
			ComponentType: ComponentTypeUpload,
			Component: &Upload{
				BaseComponent: BaseComponent{
					Id:           id + "_upload_disabled",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Accept:   "image/*",
				Full:     true,
				Disabled: true,
			}},
	}
}
