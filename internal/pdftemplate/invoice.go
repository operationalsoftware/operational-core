package pdftemplate

func invoice(_params interface{}) PDFTemplate {

	pdfTemplate := PDFTemplate{
		Title: "Invoice",
		HTML:  "<html><body><h1>Hello Invoice</h1></body></html>",
	}

	return pdfTemplate
}
