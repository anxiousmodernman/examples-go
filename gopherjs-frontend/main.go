package main

import "honnef.co/go/js/dom"
import "fmt"
import "time"

func main() {

	d := dom.GetWindow().Document()
	content := d.GetElementByID("content")
	fmt.Println(content)
	elm := d.CreateElement("h1")
	elm.SetInnerHTML("Hello World")
	content.AppendChild(elm)

	elm2 := d.CreateElement("div")
	elm2.SetInnerHTML(tmpl)
	content.AppendChild(elm2)

	go func() {
		var i int
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-ticker.C:
				if i%2 == 0 {
					elm2.SetInnerHTML(tmpl2)
				} else {
					elm2.SetInnerHTML(tmpl)
				}
				i++
			}
		}
	}()

	b := d.CreateElement("button")
	b.SetInnerHTML("click me")
	b.AddEventListener("click", true, func(e dom.Event) {
		fmt.Println("GOT EVENT", e)
	})
	content.AppendChild(b)

}

var tmpl = `
<ul>
  <li>foo</li>
  <li>baz</li>
  <li>zed</li>
</ul>
`

var tmpl2 = `
<ul>
  <li>zerb</li>
  <li>derb</li>
  <li>birb</li>
</ul>
`
