package utils

type extensionsConfig struct {
	Viewed    []string
	Edited    []string
	Converted []string
}

type extensionTypesConfig struct {
	Spreadsheet  []string
	Presentation []string
	Document     []string
}

//TODO: Rebuild with maps (N vs 1 complexity)
func isInList(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

var extensions = extensionsConfig{
	Viewed: []string{".pdf", ".djvu", ".xps"},
	Edited: []string{".docx", ".xlsx", ".csv", ".pptx", ".txt"},
}

var extensionTypes = extensionTypesConfig{
	Spreadsheet: []string{
		".xls", ".xlsx", ".xlsm",
		".xlt", ".xltx", ".xltm",
		".ods", ".fods", ".ots", ".csv",
	},
	Presentation: []string{
		".pps", ".ppsx", ".ppsm",
		".ppt", ".pptx", ".pptm",
		".pot", ".potx", ".potm",
		".odp", ".fodp", ".otp",
	},
	Document: []string{
		".doc", ".docx", ".docm",
		".dot", ".dotx", ".dotm",
		".odt", ".fodt", ".ott", ".rtf", ".txt",
		".html", ".htm", ".mht",
		".pdf", ".djvu", ".fb2", ".epub", ".xps",
	},
}

func GetFileType(fileExt string) string {
	if isInList(fileExt, extensionTypes.Document) {
		return "word"
	} else if isInList(fileExt, extensionTypes.Spreadsheet) {
		return "slide"
	} else if isInList(fileExt, extensionTypes.Spreadsheet) {
		return "cell"
	}
	return "word"
}
