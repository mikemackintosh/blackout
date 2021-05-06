package blackout

import (
	"fmt"

	"github.com/keybase/go-keychain"
)

// getDecryptKey will connect to the keychain and query
// the Chrome Safe Storage password to decrypt cookies.
func getDecryptKey() (string, error) {
	var err error

	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService("Chrome Safe Storage")
	query.SetAccount("Chrome")
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)
	results, err := keychain.QueryItem(query)
	if err != nil {
		return "", err
	} else if len(results) != 1 {
		return "", fmt.Errorf("password not found")
	}

	return string(results[0].Data), nil
}
