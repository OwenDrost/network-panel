package util

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
)

// AESDecrypt decrypts base64-encoded AES-GCM ciphertext using SHA256(secret) key.
func AESDecrypt(secret string, encryptedBase64 string) ([]byte, error) {
    if secret == "" { return nil, fmt.Errorf("empty secret") }
    if encryptedBase64 == "" { return nil, fmt.Errorf("empty data") }

    enc, err := base64.StdEncoding.DecodeString(encryptedBase64)
    if err != nil { return nil, fmt.Errorf("base64 decode: %w", err) }
    key := sha256.Sum256([]byte(secret))
    block, err := aes.NewCipher(key[:])
    if err != nil { return nil, fmt.Errorf("new cipher: %w", err) }
    gcm, err := cipher.NewGCM(block)
    if err != nil { return nil, fmt.Errorf("new gcm: %w", err) }
    if len(enc) < gcm.NonceSize() { return nil, fmt.Errorf("cipher too short") }
    nonce := enc[:gcm.NonceSize()]
    ciphertext := enc[gcm.NonceSize():]
    return gcm.Open(nil, nonce, ciphertext, nil)
}

