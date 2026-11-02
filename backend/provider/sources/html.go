package sources

import "golang.org/x/net/html"

// Helpers for reading and editing parsed pages, shared by the sources that
// rework a chapter page before it is translated.

// collectParagraphs appends the <p> elements under node, in document order,
// to paragraphs. It does not look inside a <p>.
func collectParagraphs(node *html.Node, paragraphs *[]*html.Node) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "p" {
			*paragraphs = append(*paragraphs, child)
			continue
		}
		collectParagraphs(child, paragraphs)
	}
}

// findHTMLNode returns the first node under node, in document order, that
// match accepts, or nil.
func findHTMLNode(node *html.Node, match func(*html.Node) bool) *html.Node {
	if match(node) {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := findHTMLNode(child, match); found != nil {
			return found
		}
	}
	return nil
}

// hasHTMLAttr reports whether node is an element with the attribute key.
func hasHTMLAttr(node *html.Node, key string) bool {
	if node.Type != html.ElementNode {
		return false
	}
	for _, a := range node.Attr {
		if a.Key == key {
			return true
		}
	}
	return false
}

// htmlAttr returns the value of node's attribute key, or "" if it has none.
func htmlAttr(node *html.Node, key string) string {
	for _, a := range node.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// setHTMLAttr sets node's attribute key to value, adding it if needed.
func setHTMLAttr(node *html.Node, key, value string) {
	for i := range node.Attr {
		if node.Attr[i].Key == key {
			node.Attr[i].Val = value
			return
		}
	}
	node.Attr = append(node.Attr, html.Attribute{Key: key, Val: value})
}
