// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package router

import (
	"reflect"
	"testing"

	msg "github.com/ChainSafe/ChainBridge/message"
)

type mockWriter struct {
	msgs []msg.Message
}

func (w *mockWriter) Start() error { return nil }
func (w *mockWriter) Stop() error  { return nil }

func (w *mockWriter) ResolveMessage(msg msg.Message) bool {
	w.msgs = append(w.msgs, msg)
	return true
}

func TestRouter(t *testing.T) {
	router := NewRouter()

	ethW := &mockWriter{msgs: *new([]msg.Message)}
	router.Listen(msg.ChainId(0), ethW)

	ctfgW := &mockWriter{msgs: *new([]msg.Message)}
	router.Listen(msg.ChainId(1), ctfgW)

	msgEthToCtfg := msg.Message{
		Source:      msg.ChainId(0),
		Destination: msg.ChainId(1),
	}

	msgCtfgToEth := msg.Message{
		Source:      msg.ChainId(1),
		Destination: msg.ChainId(0),
	}

	err := router.Send(msgCtfgToEth)
	if err != nil {
		t.Fatal(err)
	}
	err = router.Send(msgEthToCtfg)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(ethW.msgs[0], msgCtfgToEth) {
		t.Error("Unexpected message")
	}

	if !reflect.DeepEqual(ctfgW.msgs[0], msgEthToCtfg) {
		t.Error("Unexpected message")
	}
}
