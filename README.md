# Otter
Otter is a toolkit library to build UIs with go and templ

## CLI
`go install github.com/martinmunillas/otter/cmd/otter`
Run otter to run a development server for go and templ

## Packages

### otter
Almost unstyled common UI components
### otter/env
Environment variables utilities
### otter/response/send
JSON and HTML response sender utilities
### otter/api
Api utilities
### otter/i18n
I18n utilities
#### How to use
##### Set up
First make sure you have some translations files
```json
// en.json
{
    "helloWorld": "Hello World",
}
```
```json
// es.json
{
    "helloWorld": "Hola Mundo",
}
```
Then read the file either at runtime or compile time and pass it to i18n.AddLocale()
###### Compile time
```go
// translations.go
package translations

import (
	_ "embed"

	"github.com/martinmunillas/otter/i18n"
)
//go:embed en.json
var EnJson []byte

//go:embed es.json
var EsJson []byte


func init() {
	i18n.AddLocale("en", bytes.NewReader(translations.EnJson))
	i18n.AddLocale("es", bytes.NewReader(translations.EsJson))
}
```
###### Runtime

```go
// translations.go
package translations

import (
	"os"

	"github.com/martinmunillas/otter/i18n"
)

func main() {
	file, _ := os.Open("./en.json")
	i18n.AddLocale("en", file)
	_ = file.Close()

	file, _ = os.Open("./es.json")
	i18n.AddLocale("es", file)
	_ = file.Close()
}

```
You could even fetch the translations from a remote source
Once you added your locales, make sure you use the i18n.Middleware()
```go
// main.go
	_ = http.ListenAndServe(api.PortString(8080), i18n.Middleware(mux))
```
This middleware will make sure to set the user's locale to the context, making it available with i18n.FromContext()

##### Locale switching
The user's locale is determined by the `otter-lang` cookie but when it hasn't been set the `'Accept-Language'` header will be used.
The `otter-lang` cookie will be set only whenever the user decides to change their default locale.
The middleware will enable a `/set-locale` endpoint which will read the `locale` key from the request body and set it to the `otter-lang` cookie.
You can create your own locale setters but there is already a i18n.LanguageSelector component that wraps this endpoint. You can style it from css as this is a `<select />` element with a `"language-selector"` class.
##### Translations usage
For the use of translations in your templ files you only need to call i18n.T() with your context and the translation key.
```templ
// hello_world.templ
templ HelloWorld() {
    <h1>
        @i18n.T(ctx, "helloWorld")
    </h1>
}
```
If your translations contain raw html, for example `"hello <b>world</b>"`, you can use i18n.RawT().
```templ
// hello_world.templ
templ HelloWorld() {
    <h1>
        @i18n.RawT(ctx, "helloWorld")
    </h1>
}
```
There are also replacements available, these are super useful when the content needs further styling or dynamic content.
```json
// en.json
{
    "hello": "Hello {name}, we are {logo} and we are stoked to have you with us!"s
}
```
```templ
// hello.templ
css logoStyles() {
    color: tomato;
    font-weight: 900;
    font-size: 0.875rem;
}

templ Logo() {
    <p class={ logoStyles() }>AmazingStartup</p>
}

templ Hello(name string) {
    <h1>
        @i18n.T(ctx, "hello", map[string]any{
            name: name,
            logo: Logo(),
        })
    </h1>
}
```