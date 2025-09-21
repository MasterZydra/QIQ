package optionsInfo

import (
	"QIQ/cmd/qiq/config"
	"QIQ/cmd/qiq/ini"
	"QIQ/cmd/qiq/runtime"
	"fmt"
	"slices"
)

func printPhpInfo(context runtime.Context) {
	printfln := func(format string, a ...any) {
		context.Interpreter.Println(fmt.Sprintf(format, a...))
	}

	println := func(str string) {
		context.Interpreter.Println(str)
	}

	println("<!DOCTYPE html>")
	println("<html>")
	println("<head>")
	printfln("<title>%s - phpinfo()</title>", config.SoftwareVersion)

	// MARK: Styling
	println(`<style type="text/css">`)
	println("h1, h2 { text-align: center; }")
	println("table { width: 100%; max-width: 50em; margin: 1em auto auto auto; border-collapse: collapse; }")
	println("tr:nth-child(even) {background-color: #2d9de049;}")
	println("th { padding: 0.5em; border: solid 1px gray; background-color: #2e9cdfff; }")
	println("td { padding: 0.5em; border: solid 1px gray; }")
	println("</style>")

	println("</head>")
	println("<body>")
	println("<main>")

	printfln("<h1>%s</h1>", config.SoftwareVersion)
	println("<table>")
	printfln("<tr><td><strong>Production mode</strong></td><td>%t</td></tr>", !config.IsDevMode)
	println("</table>")

	// MARK: Configuration
	println("<h1>Configuration</h1>")

	// MARK: Configuration - Core
	println("<h2>Core</h2>")
	println("<table>")
	printfln("<tr><td><strong>QIQ Version</strong></td><td>%s</td></tr>", config.QIQVersion)
	printfln("<tr><td><strong>PHP Version</strong></td><td>%s</td></tr>", config.Version)
	println("</table>")

	println("<table>")
	println("<tr><th>Directive</th><th>Local Value</th></tr>")
	directives := ini.GetDirectives()
	slices.Sort(directives)
	for _, directive := range directives {
		if ini.IsBool(directive) {
			printfln("<tr><td><strong>%s</strong></td><td>%t</td></tr>", directive, context.Interpreter.GetIni().GetBool(directive))
			continue
		}

		value, err := context.Interpreter.GetIni().Get(directive)
		if err != nil {
			value = ""
		}
		printfln("<tr><td><strong>%s</strong></td><td>%s</td></tr>", directive, value)
	}
	println("</table>")

	// MARK: Credits
	println("<h1>Credits</h1>")

	println("<table>")
	println("<tr><th>Authors</th></tr>")
	println("<tr><td>David Hein</td></tr>")
	println("</table>")

	// MARK: Credits - License
	println("<h2>QIQ License</h2>")
	println("<table>")
	println(`<tr><td>MIT License<br>
<br>
Copyright (c) 2024 MasterZydra
<br>
Permission is hereby granted, free of charge, to any person obtaining a copy<br>
of this software and associated documentation files (the "Software"), to deal<br>
in the Software without restriction, including without limitation the rights<br>
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell<br>
copies of the Software, and to permit persons to whom the Software is<br>
furnished to do so, subject to the following conditions:<br>
<br>
The above copyright notice and this permission notice shall be included in all<br>
copies or substantial portions of the Software.<br>
<br>
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR<br>
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,<br>
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE<br>
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER<br>
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,<br>
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE<br>
SOFTWARE.
</td></tr>`)
	println("</table>")

	println("</main>")
	println("</body>")
	println("</html>")
}
