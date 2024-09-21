/* Component static files
 */
package component

import "embed"

//go:embed css js favicon.svg
var Static embed.FS

var JSLibs []string = []string{
	"https://unpkg.com/htmx.org@2.0.2",
	"https://unpkg.com/htmx-ext-remove-me@2.0.0/remove-me.js",
}
