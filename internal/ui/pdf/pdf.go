package pdf

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/olle/docops/internal/view"
)

const (
	pageWidth  = 595.28
	pageHeight = 841.89
	margin     = 54.0
	footerY    = 30.0
)

type Presenter struct{}

func (Presenter) Render(w io.Writer, snapshot view.Snapshot) error {
	pages := buildPages(snapshot)
	return writePDF(w, pages, snapshot.Greeting)
}

type page struct {
	content bytes.Buffer
}

type document struct {
	pages []page
	y     float64
}

func newDocument() *document {
	d := &document{}
	d.newPage()
	return d
}

func (d *document) newPage() {
	d.pages = append(d.pages, page{})
	d.y = pageHeight - margin
}

func (d *document) current() *bytes.Buffer {
	return &d.pages[len(d.pages)-1].content
}

func (d *document) line(text, font string, size float64) {
	lines := wrapText(text, maxChars(size))
	if len(lines) == 0 {
		d.space(size)
		return
	}
	for _, line := range lines {
		if d.y < margin+size {
			d.newPage()
		}
		fmt.Fprintf(d.current(), "BT /%s %.2f Tf 1 0 0 1 %.2f %.2f Tm %s Tj ET\n", font, size, margin, d.y, pdfLiteral(line))
		d.y -= size * 1.5
	}
}

func (d *document) space(size float64) {
	if d.y-size < margin {
		d.newPage()
	}
	d.y -= size
}

func (d *document) rule() {
	if d.y < margin+24 {
		d.newPage()
	}
	fmt.Fprintf(d.current(), "0.6 w 0.75 0.75 0.75 RG %.2f %.2f m %.2f %.2f l S\n", margin, d.y, pageWidth-margin, d.y)
	d.y -= 16
}

func buildPages(snapshot view.Snapshot) []page {
	d := newDocument()

	d.line(snapshot.Greeting, "F2", 22)
	d.line("Published specification catalog", "F1", 10)
	d.space(10)

	if len(snapshot.Warnings) > 0 {
		d.line("Warnings", "F2", 12)
		for _, warning := range snapshot.Warnings {
			d.line("- "+warning, "F1", 10)
		}
		d.space(8)
	}

	summaries := view.Summaries(snapshot)
	d.line(fmt.Sprintf("Specifications (%d)", len(summaries)), "F2", 14)
	for _, doc := range summaries {
		d.line(fmt.Sprintf("- %s: %s", doc.KindLabel, doc.DisplayName), "F1", 10)
	}
	d.space(14)

	d.line("Details", "F2", 16)
	for _, doc := range summaries {
		d.rule()
		d.line(doc.KindLabel+": "+doc.DisplayName, "F2", 14)
		field(d, "Source", doc.Source)
		field(d, "Specification version", doc.SpecVersion)
		field(d, "ID", doc.ID)
		field(d, "Version", doc.Version)
		field(d, "Title", doc.Title)
		field(d, "Name", doc.Name)
		field(d, "Description", doc.Description)
		structure(d, "Channels", doc.Structure.Channels)
		structure(d, "Operations", doc.Structure.Operations)
		structure(d, "Paths", doc.Structure.Paths)
		structure(d, "Properties", doc.Structure.Properties)
		structure(d, "Fields", doc.Structure.Fields)
		d.space(8)
	}

	for i := range d.pages {
		fmt.Fprintf(&d.pages[i].content, "BT /F1 8 Tf 1 0 0 1 %.2f %.2f Tm %s Tj ET\n", margin, footerY, pdfLiteral(fmt.Sprintf("docops - %d / %d", i+1, len(d.pages))))
	}
	return d.pages
}

func field(d *document, label, value string) {
	if value == "" {
		return
	}
	d.line(label+": "+value, "F1", 10)
}

func structure(d *document, label string, values []string) {
	if len(values) == 0 {
		return
	}
	d.line(label+": "+strings.Join(values, ", "), "F1", 10)
}

func maxChars(size float64) int {
	width := pageWidth - 2*margin
	chars := int(width / (size * 0.52))
	if chars < 12 {
		return 12
	}
	return chars
}

func wrapText(text string, maxChars int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	var current strings.Builder
	flush := func() {
		if current.Len() > 0 {
			lines = append(lines, current.String())
			current.Reset()
		}
	}
	for _, word := range words {
		for len([]rune(word)) > maxChars {
			flush()
			runes := []rune(word)
			lines = append(lines, string(runes[:maxChars]))
			word = string(runes[maxChars:])
		}
		currentLen := len([]rune(current.String()))
		switch {
		case current.Len() == 0:
			current.WriteString(word)
		case currentLen+1+len([]rune(word)) <= maxChars:
			current.WriteByte(' ')
			current.WriteString(word)
		default:
			flush()
			current.WriteString(word)
		}
	}
	flush()
	return lines
}

func pdfLiteral(text string) string {
	var b strings.Builder
	b.WriteByte('(')
	for _, r := range text {
		switch r {
		case '(', ')', '\\':
			b.WriteByte('\\')
			b.WriteByte(byte(r))
		case '\t':
			b.WriteString("    ")
		case '\n', '\r':
			b.WriteByte(' ')
		case '\u2018', '\u2019', '\u201c', '\u201d':
			b.WriteByte('\'')
		case '\u2013', '\u2014':
			b.WriteByte('-')
		case '\u00a0':
			b.WriteByte(' ')
		default:
			if r < 32 || r > 255 || (r >= 128 && r <= 159) || !unicode.IsPrint(r) {
				b.WriteByte('?')
				continue
			}
			b.WriteByte(byte(r))
		}
	}
	b.WriteByte(')')
	return b.String()
}

func writePDF(w io.Writer, pages []page, title string) error {
	pageCount := len(pages)
	fontRegularID := 3 + 2*pageCount
	fontBoldID := fontRegularID + 1
	infoID := fontBoldID + 1
	objectCount := infoID

	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")
	buf.WriteString("%\xe2\xe3\xcf\xd3\n")

	offsets := make([]int, objectCount+1)
	writeObject := func(id int, body string) {
		offsets[id] = buf.Len()
		fmt.Fprintf(&buf, "%d 0 obj\n%s\nendobj\n", id, body)
	}

	writeObject(1, "<< /Type /Catalog /Pages 2 0 R >>")

	kids := make([]string, pageCount)
	for i := range kids {
		kids[i] = fmt.Sprintf("%d 0 R", 3+2*i)
	}
	writeObject(2, fmt.Sprintf("<< /Type /Pages /Kids [%s] /Count %d >>", strings.Join(kids, " "), pageCount))

	for i, p := range pages {
		pageID := 3 + 2*i
		contentID := pageID + 1
		writeObject(pageID, fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %.2f %.2f] /Resources << /Font << /F1 %d 0 R /F2 %d 0 R >> >> /Contents %d 0 R >>", pageWidth, pageHeight, fontRegularID, fontBoldID, contentID))
		stream := p.content.Bytes()
		body := fmt.Sprintf("<< /Length %d >>\nstream\n", len(stream)) + string(stream) + "endstream"
		writeObject(contentID, body)
	}

	writeObject(fontRegularID, "<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding /WinAnsiEncoding >>")
	writeObject(fontBoldID, "<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica-Bold /Encoding /WinAnsiEncoding >>")
	writeObject(infoID, fmt.Sprintf("<< /Producer (docops) /Title %s >>", pdfLiteral(title)))

	xrefOffset := buf.Len()
	fmt.Fprintf(&buf, "xref\n0 %d\n", objectCount+1)
	buf.WriteString("0000000000 65535 f \n")
	for id := 1; id <= objectCount; id++ {
		fmt.Fprintf(&buf, "%010d 00000 n \n", offsets[id])
	}
	fmt.Fprintf(&buf, "trailer\n<< /Size %d /Root 1 0 R /Info %d 0 R >>\nstartxref\n%d\n%%%%EOF\n", objectCount+1, infoID, xrefOffset)

	_, err := w.Write(buf.Bytes())
	return err
}
