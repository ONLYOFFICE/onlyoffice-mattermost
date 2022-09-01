package validator

import "strings"

func IsMobile(useragent string) string {
	mobiles := []string{"Mobile Explorer", "Palm", "Motorola", "Nokia", "Palm", "Apple iPhone", "iPad", "Apple iPod Touch", "Sony Ericsson", "Sony Ericsson", "BlackBerry", "O2 Cocoon", "Treo", "LG", "Amoi", "XDA", "MDA", "Vario", "HTC", "Samsung",
		"Sharp", "Siemens", "Alcatel", "BenQ", "HP iPaq", "Motorola", "PlayStation Portable", "PlayStation 3", "PlayStation Vita", "Danger Hiptop", "NEC", "Panasonic", "Philips", "Sagem", "Sanyo", "SPV", "ZTE", "Sendo", "Nintendo DSi", "Nintendo DS", "Nintendo 3DS",
		"Nintendo Wii", "Open Web", "OpenWeb", "Android", "Symbian", "SymbianOS", "Palm", "Symbian S60", "Windows CE", "Obigo", "Netfront Browser", "Openwave Browser", "Mobile Explorer", "Opera Mini", "Opera Mobile", "Firefox Mobile",
		"Digital Paths", "AvantGo", "Xiino", "Novarra Transcoder", "Vodafone", "NTT DoCoMo", "O2", "mobile", "wireless", "j2me", "midp", "cldc", "up.link", "up.browser", "smartphone", "cellphone", "Generic Mobile"}

	for _, device := range mobiles {
		if strings.Contains(useragent, device) {
			return "mobile"
		}
	}
	return "desktop"
}
