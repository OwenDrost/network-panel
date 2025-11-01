package controller

import (
	"fmt"
	"log"
	"net"
	"time"
)

// diagnoseFromNode asks a node via WS to perform TCP connect tests.
// Node agent should support command: {type: "Diagnose", data: {requestId, host, port, protocol, count, timeoutMs}}
// and reply: {type: "DiagnoseResult", requestId, data: {success, averageTime, packetLoss, message}}
func diagnoseFromNode(nodeID int64, host string, port int, count int, timeoutMs int) (avg float64, loss float64, ok bool, msg string) {
	avg, loss, ok, msg, _ = diagnoseFromNodeCtx(nodeID, host, port, count, timeoutMs, nil)
	return
}

func diagnoseFromNodeCtx(nodeID int64, host string, port int, count int, timeoutMs int, ctx map[string]interface{}) (avg float64, loss float64, ok bool, msg string, reqId string) {
	if nodeID == 0 {
		return 0, 100, false, "节点不可用", ""
	}
	rid := RandUUID()
	payload := map[string]interface{}{
		"requestId": rid,
		"host":      host,
		"port":      port,
		"protocol":  "tcp",
		"count":     count,
		"timeoutMs": timeoutMs,
	}
	if ctx != nil {
		payload["ctx"] = ctx
	}
	log.Printf("%s", fmt.Sprintf("{\"event\":\"diagnose_begin\",\"mode\":\"tcp\",\"nodeId\":%d,\"reqId\":\"%s\",\"host\":\"%s\",\"port\":%d,\"count\":%d,\"timeoutMs\":%d,\"ctx\":%v}", nodeID, rid, host, port, count, timeoutMs, ctx))
	if err := sendWSCommand(nodeID, "Diagnose", payload); err != nil {
		return 0, 100, false, "节点未在线或密钥不匹配", rid
	}
	res, ok2 := RequestDiagnose(nodeID, payload, 8*time.Second)
	if !ok2 {
		return 0, 100, false, "节点未响应诊断", rid
	}
	data, _ := res["data"].(map[string]interface{})
	if data == nil {
		return 0, 100, false, "诊断结果无数据", rid
	}
	if succ, _ := data["success"].(bool); !succ {
		msg, _ := data["message"].(string)
		return 0, 100, false, msg, rid
	}
	if v, ok := toFloat(data["averageTime"]); ok {
		avg = v
	}
	if v, ok := toFloat(data["packetLoss"]); ok {
		loss = v
	}
	if m, ok := data["message"].(string); ok {
		msg = m
	}
	return avg, loss, true, msg, rid
}

// tcpDialFallback performs multiple TCP connections from backend host to estimate latency and loss.
func tcpDialFallback(host string, port int, count int, timeoutMs int) (avg float64, loss float64, ok bool, msg string) {
	if host == "" || port == 0 {
		return 0, 100, false, "无效目标"
	}
	var sum float64
	var succ int
	for i := 0; i < count; i++ {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Duration(timeoutMs)*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			elapsed := time.Since(start).Seconds() * 1000
			sum += elapsed
			succ++
		}
	}
	if succ == 0 {
		return 0, 100, false, "后端直连失败"
	}
	avg = sum / float64(succ)
	loss = float64(count-succ) * 100.0 / float64(count)
	return avg, loss, true, "后端直连测量"
}

// diagnosePingFromNode asks a node to perform ICMP ping tests to host.
func diagnosePingFromNode(nodeID int64, host string, count int, timeoutMs int) (avg float64, loss float64, ok bool, msg string) {
	avg, loss, ok, msg, _ = diagnosePingFromNodeCtx(nodeID, host, count, timeoutMs, nil)
	return
}

func diagnosePingFromNodeCtx(nodeID int64, host string, count int, timeoutMs int, ctx map[string]interface{}) (avg float64, loss float64, ok bool, msg string, reqId string) {
	if nodeID == 0 {
		return 0, 100, false, "节点不可用", ""
	}
	rid := RandUUID()
	payload := map[string]interface{}{
		"requestId": rid,
		"host":      host,
		"mode":      "icmp",
		"count":     count,
		"timeoutMs": timeoutMs,
	}
	if ctx != nil {
		payload["ctx"] = ctx
	}
	log.Printf("%s", fmt.Sprintf("{\"event\":\"diagnose_begin\",\"mode\":\"icmp\",\"nodeId\":%d,\"reqId\":\"%s\",\"host\":\"%s\",\"count\":%d,\"timeoutMs\":%d,\"ctx\":%v}", nodeID, rid, host, count, timeoutMs, ctx))
	if err := sendWSCommand(nodeID, "Diagnose", payload); err != nil {
		return 0, 100, false, "节点未在线或密钥不匹配", rid
	}
	res, ok2 := RequestDiagnose(nodeID, payload, 8*time.Second)
	if !ok2 {
		return 0, 100, false, "节点未响应诊断", rid
	}
	data, _ := res["data"].(map[string]interface{})
	if data == nil {
		return 0, 100, false, "诊断结果无数据", rid
	}
	if succ, _ := data["success"].(bool); !succ {
		m, _ := data["message"].(string)
		return 0, 100, false, m, rid
	}
	if v, ok := toFloat(data["averageTime"]); ok {
		avg = v
	}
	if v, ok := toFloat(data["packetLoss"]); ok {
		loss = v
	}
	if m, ok := data["message"].(string); ok {
		msg = m
	}
	return avg, loss, true, msg, rid
}

func toFloat(v interface{}) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	default:
		return 0, false
	}
}
