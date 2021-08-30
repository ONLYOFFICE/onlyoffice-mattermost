/**
 *
 * (c) Copyright Ascensio System SIA 2021
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"sync"
)

var aesEncryptors sync.Map = sync.Map{}

func (e EncryptorAES) create(key []byte) (cipher.AEAD, error) {
	byteString := string(key[:])

	loaded, ok := aesEncryptors.Load(byteString)

	if ok {
		return loaded.(cipher.AEAD), nil
	}

	aes, err := aes.NewCipher(key)

	if err != nil {
		return nil, errors.New(ONLYOFFICE_LOGGER_PREFIX + "AES Encryptor Could not create a new cipher with the given key" + ONLYOFFICE_LOGGER_ENCRYPTION_SUFFIX)
	}

	aesGCM, err := cipher.NewGCM(aes)

	if err != nil {
		return nil, errors.New(ONLYOFFICE_LOGGER_PREFIX + "AES GCM Error" + ONLYOFFICE_LOGGER_ENCRYPTION_SUFFIX)
	}

	aesEncryptors.Store(byteString, aesGCM)

	return aesGCM, nil
}

func (e EncryptorAES) Encrypt(text string, key []byte) (string, error) {
	byteText := []byte(text)
	aesGCM, aesErr := e.create(key)

	if aesErr != nil {
		return "", aesErr
	}

	nonce := make([]byte, aesGCM.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.New(ONLYOFFICE_LOGGER_PREFIX + "AES GCM Nonce error" + ONLYOFFICE_LOGGER_ENCRYPTION_SUFFIX)
	}

	encrypted := aesGCM.Seal(nonce, nonce, byteText, nil)
	data := encryptedTextConversion(encrypted)

	return data, nil
}

func (e EncryptorAES) Decrypt(text string, key []byte) (string, error) {
	cipherBytes := textToEncryptedConversion(text)

	aesGCM, aesErr := e.create(key)

	if aesErr != nil {
		return "", aesErr
	}

	nonceSize := aesGCM.NonceSize()
	if len(cipherBytes) < nonceSize {
		return "", errors.New(ONLYOFFICE_LOGGER_PREFIX + "AES GCM Nonce size error" + ONLYOFFICE_LOGGER_DECRYPTION_SUFFIX)
	}

	nonce, ciphertext := cipherBytes[:nonceSize], cipherBytes[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return "", errors.New(ONLYOFFICE_LOGGER_PREFIX + "AES GCM Error" + ONLYOFFICE_LOGGER_DECRYPTION_SUFFIX)
	}

	return string(plaintext), nil
}
