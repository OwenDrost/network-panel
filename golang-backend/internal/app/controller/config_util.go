package controller

import (
    "strconv"
    strings "strings"

    dbpkg "network-panel/golang-backend/internal/db"
    "network-panel/golang-backend/internal/app/model"
)

// getConfigString returns configuration value from vite_config by name, empty if missing.
func getConfigString(name string) string {
    var it model.ViteConfig
    if err := dbpkg.DB.Where("name = ?", name).First(&it).Error; err != nil {
        return ""
    }
    return strings.TrimSpace(it.Value)
}

// getConfigInt parses configuration value as int; returns def if empty/invalid.
func getConfigInt(name string, def int) int {
    s := getConfigString(name)
    if s == "" { return def }
    // allow suffix 's' for seconds
    if strings.HasSuffix(strings.ToLower(s), "s") {
        base := strings.TrimSpace(s[:len(s)-1])
        if v, err := strconv.Atoi(base); err == nil {
            return v
        }
        return def
    }
    if v, err := strconv.Atoi(s); err == nil {
        return v
    }
    return def
}

// readDiagLocalProbeTimeoutMs reads timeout for LocalProbe in milliseconds.
// Supported keys (priority):
//  - diag_local_probe_timeout_ms (milliseconds)
//  - diag_local_probe_timeout_s  (seconds)
// Default: 3000 ms
func readDiagLocalProbeTimeoutMs() int {
    if v := getConfigInt("diag_local_probe_timeout_ms", -1); v > 0 {
        return v
    }
    if v := getConfigInt("diag_local_probe_timeout_s", -1); v > 0 {
        return v * 1000
    }
    // also accept generic key without unit, treat as seconds for human readability
    if v := getConfigInt("diag_local_probe_timeout", -1); v > 0 {
        return v * 1000
    }
    return 3000
}

