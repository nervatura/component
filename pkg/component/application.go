package component

import (
	"fmt"
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [Application] HeadLink data
type HeadLink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
	Type string `json:"type"`
}

/*
The Application component is a top-level element to which all other components belong. This element is
completely never replaced, only some of its parts can change. Its task is to load and manage all static
elements required for the operation and display of the components, such as style sheets and the htmx package.
*/
type Application struct {
	BaseComponent
	// The title value of the html document
	Title string `json:"title"`
	/*
		The theme of the application.
		[Theme] variable constants: [ThemeLight], [ThemeDark]. Default value: [ThemeLight]
	*/
	Theme string `json:"theme"`
	/*
		The htmx hx-headers attribute allows you to add to the headers that will be submitted with an AJAX request.
		Example: ut.SM{"X-CSRF-Token": "TOKEN0123456789"}
	*/
	Header ut.SM `json:"header"`
	// Any javascript libraries that the app still needs to load. The htmx javascript library is loaded automatically
	Script []string `json:"script"`
	/*
		Any additional stylesheets or resource files. The path to the style sheets of the components in the
		pkg/static/css package must be specified.
		Example: []ct.HeadLink{
			{Rel: "icon", Href: "/static/favicon.svg", Type: "image/svg+xml"},
			{Rel: "stylesheet", Href: "/public/demo.css"},
			{Rel: "stylesheet", Href: "/static/css/index.css"}}
	*/
	HeadLink []HeadLink `json:"head_link"`
	// The main component of the application, to which all other components belong.
	MainComponent ClientComponent `json:"main_component"`
}

/*
Returns all properties of the [Application]
*/
func (app *Application) Properties() ut.IM {
	return ut.MergeIM(
		app.BaseComponent.Properties(),
		ut.IM{
			"title":  app.Title,
			"theme":  app.Theme,
			"header": app.Header,
			"script": app.Script,
			"link":   app.HeadLink,
			"main":   app.MainComponent,
		})
}

/*
Returns the value of the property of the [Application] with the specified name.
*/
func (app *Application) GetProperty(propName string) interface{} {
	return app.Properties()[propName]
}

/*
It checks the value given to the property of the [Application] and always returns a valid value
*/
func (app *Application) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"theme": func() interface{} {
			return app.CheckEnumValue(ut.ToString(propValue, ""), ThemeLight, Theme)
		},
		"header": func() interface{} {
			value := ut.SetSMValue(app.Header, "", "")
			if smap, valid := propValue.(ut.SM); valid {
				value = ut.MergeSM(value, smap)
			}
			return value
		},
		"script": func() interface{} {
			value := []string{
				"https://unpkg.com/htmx.org@latest",
				"https://unpkg.com/htmx.org/dist/ext/remove-me.js",
				//"static/js/htmx.min.js",
				//"static/js/remove-me.js",
			}
			if script, valid := propValue.([]string); valid {
				value = append(value, script...)
			}
			return value
		},
		"link": func() interface{} {
			value := []HeadLink{
				{Rel: "preconnect", Href: "https://fonts.gstatic.com"},
				{Rel: "stylesheet", Href: "https://fonts.googleapis.com/css2?family=Noto+Sans:ital,wght@0,400;0,700;1,400;1,700&display=swap"},
			}
			if link, valid := propValue.([]HeadLink); valid {
				value = append(value, link...)
			}
			return value
		},
		"main": func() interface{} {
			if mc, valid := propValue.(ClientComponent); valid {
				return mc
			}
			return nil
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if app.BaseComponent.GetProperty(propName) != nil {
		return app.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [Application] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (app *Application) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"title": func() interface{} {
			app.Title = ut.ToString(propValue, "")
			return app.Title
		},
		"theme": func() interface{} {
			app.Theme = app.Validation(propName, propValue).(string)
			return app.Theme
		},
		"header": func() interface{} {
			app.Header = app.Validation(propName, propValue).(ut.SM)
			return app.Header
		},
		"script": func() interface{} {
			app.Script = app.Validation(propName, propValue).([]string)
			return app.Script
		},
		"link": func() interface{} {
			app.HeadLink = app.Validation(propName, propValue).([]HeadLink)
			return app.HeadLink
		},
		"main": func() interface{} {
			if mc, valid := propValue.(ClientComponent); valid {
				app.MainComponent = mc
			}
			return app.MainComponent
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if app.BaseComponent.GetProperty(propName) != nil {
		return app.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

/*
The [TriggerEvent] event of the user interface is forwarded to the child component registered in the
RequestMap based on the component id. If there is no component associated with the received component ID,
or the processing of the component returns an error, it returns the error message by creating a [Toast]
component.
*/
func (app *Application) OnRequest(te TriggerEvent) (re ResponseEvent) {
	if cc, found := app.RequestMap[te.Id]; found {
		return cc.OnRequest(te)
	}
	re = ResponseEvent{
		Trigger: &Toast{
			Type:  ToastTypeError,
			Value: fmt.Sprintf("Invalid parameter: %s", te.Id),
		},
		TriggerName: te.Name,
		Name:        te.Name,
		Header: ut.SM{
			HeaderRetarget: "#toast-msg",
			HeaderReswap:   SwapInnerHTML,
		},
	}
	return re
}

func (app *Application) getComponent() (res string, err error) {
	if app.MainComponent != nil {
		return app.MainComponent.Render()
	}
	return res, err
}

/*
Based on the values, it will generate the html code of the [Application] or return with an error message.
*/
func (app *Application) Render() (res string, err error) {
	app.InitProps(app)
	spinner := Spinner{}

	funcMap := map[string]any{
		"headerKeys": func() string {
			values := []string{}
			for key, value := range app.Header {
				values = append(values, key, value)
			}
			if len(values) > 0 {
				return fmt.Sprintf(`hx-headers='{"%s"}'`, strings.Join(values, `":"`))
			}
			return ""
		},
		"styleMap": func() bool {
			return len(app.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(app.Class, " ")
		},
		"spinner": func() (string, error) {
			return spinner.Render()
		},
		"main": func() (res string, err error) {
			return app.getComponent()
		},
	}
	tpl := `<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0, viewport-fit=cover" />
			<meta http-equiv="X-UA-Compatible" content="ie=edge" />
			{{ range $index, $value := .Script }}<script src="{{ $value }}"></script>{{ end }}
			{{ range $index, $link := .HeadLink }}
			<link rel="{{ $link.Rel }}" href="{{ $link.Href }}" {{ if ne $link.Type "" }}type="{{ $link.Type }}"{{ end }} />
			{{ end }}
			<title>{{ .Title }}</title>
		</head>
		<body>
		<div id="{{ .Id }}" theme="{{ .Theme }}" 
		{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }} 
		hx-ext="remove-me" {{ headerKeys }} class="{{ customClass }}">
		<div id="toast-msg"></div><div>{{ spinner }}</div>
		{{ main }}
		</div>
		</body>
	</html>`

	return ut.TemplateBuilder("application", tpl, funcMap, app)
}
