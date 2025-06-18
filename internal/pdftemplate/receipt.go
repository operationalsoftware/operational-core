package pdftemplate

func receipt(_params interface{}) PDFTemplate {

	pdfTemplate := PDFTemplate{
		Title: "Receipt",
		HTML:  "<html><body><h1>Hello Receipt</h1></body></html>",
	}

	return pdfTemplate
}
