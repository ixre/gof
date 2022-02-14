package report

var _ IExportProvider = new(ExcelProvider)

type ExcelProvider struct {
	csv IExportProvider
}

func NewExcelProvider() IExportProvider {
	return &ExcelProvider{
		csv: NewCsvProvider(),
	}
}

func (e *ExcelProvider) Export(rows []map[string]interface{},
	fields []string, names []string, formatter []IExportFormatter) (binary []byte) {
	return e.csv.Export(rows, fields, names, formatter)
}
