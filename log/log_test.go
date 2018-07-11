// Package log is wrapper for logging
package log

import (
	"testing"

	"bitbucket.org/coinplugin/proxy/common"
	"bitbucket.org/coinplugin/proxy/json"
)

func TestFormatd(t *testing.T) {
	var a uint64
	a = 55555555555555
	Warnf("%d", a)
}

func TestId(t *testing.T) {
	id := common.RandomUint64()
	// Ascending log level
	Debugd(id, "debug", "1")
	Infod(id, "info", "2")
	Warnd(id, "warn", "3")
	Errord(id, "error", "4")
}

func TestJson(t *testing.T) {
	resp := &json.RPCResponse{}
	Error(resp.String())
}

func TestGeneral(t *testing.T) {
	// Ascending log level
	Debug("debug", "1")
	Info("info", "2")
	Warn("warn", "3")
	Error("error", "4")
}

func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("Recovered in ", r)
		}
	}()

	Panic("panic")
}

func TestFatal(t *testing.T) {

	Fatal("fatal")
}
