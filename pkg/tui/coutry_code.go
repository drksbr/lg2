package tui

import "github.com/biter777/countries"

// CountryToFlag converte o código ISO 3166-1 alpha-2 do país (BR, US, UK, etc.) em um emoji de bandeira.
func CountryToFlag(countryCode string) string {
	// // Take only first 2 characters and convert to uppercase
	// if len(countryCode) < 2 {
	// 	return ""
	// }
	// countryCode = strings.ToUpper(countryCode[:2])

	// if countryCode[0] < 'A' || countryCode[0] > 'Z' || countryCode[1] < 'A' || countryCode[1] > 'Z' {
	// 	return "" // Retorna vazio se o código não for válido
	// }

	// // Convert country code to flag emoji and return with the code
	// flagOffset := 127397 // Offset for flag emojis
	// flag := string([]rune{
	// 	rune(countryCode[0]) + rune(flagOffset),
	// 	rune(countryCode[1]) + rune(flagOffset),
	// })
	// return fmt.Sprintf("(%s %s)", countryCode, flag)

	// Convert country code to flag emoji and return with the code
	countryes := countries.All()
	for _, country := range countryes {
		if country.Alpha2() == countryCode {
			return country.Emoji3()
		}
	}

	return countryCode

}
