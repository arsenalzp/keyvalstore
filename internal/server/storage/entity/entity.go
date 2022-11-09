// Storage entity
// It is used for Import and Export operation.

package entity

type ImportData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ExportData ImportData
