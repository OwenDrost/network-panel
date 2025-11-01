package util

import (
    "bufio"
    "os"
    "path/filepath"
    "strings"
)

// LoadEnv loads KEY=VALUE pairs from .env-like files into process env.
// It tries the provided paths, and if none provided, tries common locations: .env, ../.env, ../../.env
func LoadEnv(paths ...string) {
    candidates := paths
    if len(candidates) == 0 {
        candidates = []string{".env", filepath.Join("..", ".env"), filepath.Join("..", "..", ".env")}
    }
    for _, p := range candidates {
        loadEnvFile(p)
    }
}

func loadEnvFile(path string) {
    f, err := os.Open(path)
    if err != nil { return }
    defer f.Close()
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") { continue }
        // support export KEY=VALUE
        if strings.HasPrefix(line, "export ") { line = strings.TrimSpace(strings.TrimPrefix(line, "export ")) }
        kv := strings.SplitN(line, "=", 2)
        if len(kv) != 2 { continue }
        key := strings.TrimSpace(kv[0])
        val := strings.TrimSpace(kv[1])
        // strip surrounding quotes
        if len(val) >= 2 {
            if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
                val = val[1:len(val)-1]
            }
        }
        // do not override existing env
        if _, exists := os.LookupEnv(key); !exists {
            _ = os.Setenv(key, val)
        }
    }
}

