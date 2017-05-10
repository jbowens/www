package blog

import "bytes"

// snippetRenderer is a blackfriday markdown Renderer that extracts
// the paragraphs as plain text. The blackfriday package has a really
// terrible API, so this is fairly gross.
type snippetRenderer struct{}

func (_ snippetRenderer) GetFlags() int { return 0 }

func (_ snippetRenderer) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("\n\n")
}

func (_ snippetRenderer) AutoLink(out *bytes.Buffer, link []byte, _ int) {
	out.Write(link)
}

func (_ snippetRenderer) CodeSpan(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (_ snippetRenderer) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (_ snippetRenderer) Emphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (_ snippetRenderer) Link(out *bytes.Buffer, _ []byte, _ []byte, content []byte) {
	out.Write(content)
}

func (_ snippetRenderer) TripleEmphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (_ snippetRenderer) Entity(out *bytes.Buffer, entity []byte) {
	out.Write(entity)
}

func (_ snippetRenderer) NormalText(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

// should be no-ops, but blackfriday requires the callbacks
// to be called.

func (_ snippetRenderer) Header(out *bytes.Buffer, text func() bool, _ int, _ string) {
	v := out.Len()
	text()
	out.Truncate(v)
}
func (_ snippetRenderer) List(out *bytes.Buffer, text func() bool, _ int) {
	v := out.Len()
	text()
	out.Truncate(v)
}
func (_ snippetRenderer) Footnotes(out *bytes.Buffer, text func() bool) {
	v := out.Len()
	text()
	out.Truncate(v)
}

// no-ops

func (_ snippetRenderer) BlockCode(*bytes.Buffer, []byte, string)         {}
func (_ snippetRenderer) TitleBlock(*bytes.Buffer, []byte)                {}
func (_ snippetRenderer) BlockQuote(*bytes.Buffer, []byte)                {}
func (_ snippetRenderer) BlockHtml(*bytes.Buffer, []byte)                 {}
func (_ snippetRenderer) HRule(*bytes.Buffer)                             {}
func (_ snippetRenderer) ListItem(*bytes.Buffer, []byte, int)             {}
func (_ snippetRenderer) Table(*bytes.Buffer, []byte, []byte, []int)      {}
func (_ snippetRenderer) TableRow(*bytes.Buffer, []byte)                  {}
func (_ snippetRenderer) TableHeaderCell(*bytes.Buffer, []byte, int)      {}
func (_ snippetRenderer) TableCell(*bytes.Buffer, []byte, int)            {}
func (_ snippetRenderer) FootnoteRef(*bytes.Buffer, []byte, int)          {}
func (_ snippetRenderer) FootnoteItem(*bytes.Buffer, []byte, []byte, int) {}
func (_ snippetRenderer) Image(*bytes.Buffer, []byte, []byte, []byte)     {}
func (_ snippetRenderer) LineBreak(*bytes.Buffer)                         {}
func (_ snippetRenderer) RawHtmlTag(*bytes.Buffer, []byte)                {}
func (_ snippetRenderer) StrikeThrough(*bytes.Buffer, []byte)             {}
func (_ snippetRenderer) DocumentHeader(*bytes.Buffer)                    {}
func (_ snippetRenderer) DocumentFooter(*bytes.Buffer)                    {}
