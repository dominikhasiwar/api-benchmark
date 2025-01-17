package models

type ExportDataModel struct {
	Content  []byte `json:"content"`
	MimeType string `json:"mimeType"`
	FileName string `json:"fileName"`
}
