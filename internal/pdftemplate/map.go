package pdftemplate

type PDFTemplate struct {
	Title string
	HTML  string
}

type PDFTemplateFunc func(interface{}) PDFTemplate

var TemplateRegistry = map[string]PDFTemplateFunc{
	"invoice": invoice,
	"receipt": receipt,
}
