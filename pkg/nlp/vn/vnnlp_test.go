package vn

import (
	"testing"
)

func TestVnNLP(t *testing.T) {
	// var text = `Trả lời Zing.vn sau khi VKS đưa ra quan điểm về vấn đề ly hôn cũng như chia tài sản và cổ phần, bà Lê Hoàng Diệp Thảo cho hay đó là đề nghị nhân văn, tôn trọng nữ quyền.`

	var s = NewNLPServer()

	// 1. Test starts server
	if s.IsServerAlive() != true {
		t.Errorf(`Expect IsAlive: true`)
	}

	// Test stops server
	s.Stop()

	if s.IsServerAlive() != false {
		t.Errorf(`Expect IsAlive: false`)
	}
}
