package util

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "os"
    "strings"
    "time"
    "strconv"
)

var jwtSecret = func() string { return os.Getenv("JWT_SECRET") }

type jwtHeader struct {
    Alg string `json:"alg"`
    Typ string `json:"typ"`
}

type jwtPayload struct {
    Sub    string `json:"sub"`
    Iat    int64  `json:"iat"`
    Exp    int64  `json:"exp"`
    User   string `json:"user"`
    Name   string `json:"name"`
    RoleID int    `json:"role_id"`
}

func b64(data []byte) string {
    return base64.RawURLEncoding.EncodeToString(data)
}

func sign(content string, secret string) string {
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(content))
    return b64(mac.Sum(nil))
}

func GenerateToken(userID int64, username string, roleID int) string {
    head := jwtHeader{Alg: "HmacSHA256", Typ: "JWT"}
    iat := time.Now().Unix()
    exp := time.Now().Add(90 * 24 * time.Hour).Unix()
    payload := jwtPayload{Sub: toStr(userID), Iat: iat, Exp: exp, User: username, Name: username, RoleID: roleID}

    hb, _ := json.Marshal(head)
    pb, _ := json.Marshal(payload)
    eh, ep := b64(hb), b64(pb)
    sig := sign(eh+"."+ep, jwtSecret())
    return eh + "." + ep + "." + sig
}

func ValidateToken(token string) bool {
    parts := strings.Split(token, ".")
    if len(parts) != 3 { return false }
    if jwtSecret() == "" { return false }
    sig := sign(parts[0]+"."+parts[1], jwtSecret())
    if sig != parts[2] { return false }
    // check exp
    var p jwtPayload
    dec, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil { return false }
    if err := json.Unmarshal(dec, &p); err != nil { return false }
    if p.Exp <= time.Now().Unix() { return false }
    return true
}

func GetUserID(token string) int64 {
    var p jwtPayload
    parts := strings.Split(token, ".")
    dec, _ := base64.RawURLEncoding.DecodeString(parts[1])
    _ = json.Unmarshal(dec, &p)
    return toInt64(p.Sub)
}

func GetRoleID(token string) int {
    var p jwtPayload
    parts := strings.Split(token, ".")
    dec, _ := base64.RawURLEncoding.DecodeString(parts[1])
    _ = json.Unmarshal(dec, &p)
    return p.RoleID
}

func toStr(v int64) string { return strconv.FormatInt(v, 10) }
func toInt64(s string) int64 {
    i, _ := strconv.ParseInt(s, 10, 64)
    return i
}
