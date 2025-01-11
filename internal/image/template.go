package img

import (
	"strings"

	"golang.org/x/net/html"
)

type HTMLTemplate string

var Template HTMLTemplate = `<img width="796" height="769" src="https://cloudageio.com/wp-content/uploads/2023/11/about-image.png" alt="" title="" sizes="(min-width: 0px) and (max-width: 480px) 480px, (min-width: 481px) 796px, 100vw" srcset="https://cloudageio.com/wp-content/uploads/2023/11/about-image.png 796w,https://cloudageio.com/wp-content/uploads/2023/11/about-image-480x464.png 480w" />`
var Template2 HTMLTemplate = `<p class="hello">salut</p>`

func (templ HTMLTemplate) Parse() (*html.Node, error) {

	node, err := html.Parse(strings.NewReader(string(templ)))
	if err != nil {
		return nil, err
	}
	return node, nil
}
