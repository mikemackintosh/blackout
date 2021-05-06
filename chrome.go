package blackout

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	aescbcSalt            = `saltysalt`
	aescbcIV              = `                `
	aescbcIterationsLinux = 1
	aescbcIterationsMacOS = 1003
	aescbcLength          = 16
)

// Cookies is a collection of Cookie.
type Cookies []*Cookie

// Cookie maps to the sqlite3 database.
type Cookie struct {
	Creation   int64 `json:"-" csv:"-" gorm:"column:creation_utc"`
	LastAccess int64 `json:"-" csv:"-" gorm:"column:last_access_utc"`
	Expires    int64 `json:"-" csv:"-" gorm:"column:expires_utc"`

	Domain   string `json:"domain" csv:"domain" gorm:"column:host_key"`
	Name     string `json:"name" csv:"name"`
	Value    string `json:"value" csv:"value" `
	Path     string `json:"path" csv:"path" `
	Priority int64  `json:"priority" csv:"priority" `

	IsSecure     int `json:"secure" csv:"secure" gorm:"column:is_secure"`
	IsHTTPOnly   int `json:"httponly" csv:"httponly" gorm:"column:is_httponly"`
	IsPersistent int `json:"persistent" csv:"persistent" gorm:"column:is_persistent"`
	IsSameParty  int `json:"sameparty" csv:"sameparty" gorm:"column:is_same_party"`
	HasExpires   int `json:"-" csv:"-" gorm:"column:has_expires"`

	EncryptedValue string `json:"-" csv:"-" gorm:"column:encrypted_value"`
	SameSite       int    `json:"-" csv:"samesite" gorm:"column:samesite"`
	SourceScheme   int    `json:"-" csv:"-" gorm:"column:source_scheme"`
	SourcePort     int    `json:"-" csv:"-" gorm:"column:source_port"`
}

// GetCookies will query cookies, and return an error.
func GetCookies(profilePath string) (Cookies, error) {
	var cookies Cookies

	db, err := gorm.Open(sqlite.Open(profilePath), &gorm.Config{})
	if err != nil {
		return cookies, err
	}

	if err = db.Find(&cookies).Error; err != nil {
		return cookies, err
	}

	key, err := getDecryptKey()
	if err != nil {
		return cookies, err
	}

	for _, cookie := range cookies {
		encryptedValue := strings.TrimPrefix(cookie.EncryptedValue, "v10")

		value, err := decryptCookieValue(key, encryptedValue)
		if err != nil {
			continue
		}
		cookie.Value = value

	}

	return cookies, nil
}

// decryptCookieValue will decrypt the cookie value.
func decryptCookieValue(password, encrypted string) (string, error) {
	key := pbkdf2.Key([]byte(password), []byte(aescbcSalt), aescbcIterationsMacOS, aescbcLength, sha1.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decrypted := make([]byte, len(encrypted))
	cbc := cipher.NewCBCDecrypter(block, []byte(aescbcIV))
	cbc.CryptBlocks(decrypted, []byte(encrypted))

	if len(decrypted) == 0 {
		return "", fmt.Errorf("not enough bits")
	}

	if len(decrypted)%aescbcLength != 0 {
		return "", fmt.Errorf("decrypted data block length is not a multiple of %d", aescbcLength)
	}
	paddingLen := int(decrypted[len(decrypted)-1])
	if paddingLen > 16 {
		return "", fmt.Errorf("invalid last block padding length: %d", paddingLen)
	}

	return string(decrypted[:len(decrypted)-paddingLen]), nil
}
