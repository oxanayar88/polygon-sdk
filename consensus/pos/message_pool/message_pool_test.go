package pool

import "testing"

func TestMessagePool_Insert(t *testing.T) {
	data1 := []byte{0x1}

	set := newMockValidatorSet([]string{"A", "B"})

	m := NewMessagePool(NodeID("A"), nil)
	m.Reset(set)

	m.addImpl(&Message{
		Data: data1,
		From: NodeID("A"),
	})
}

type mockValidatorSet struct {
	ids []string
}

func newMockValidatorSet(ids []string) ValidatorSet {
	return &mockValidatorSet{ids: ids}
}

func (m *mockValidatorSet) Includes(n NodeID) bool {
	for _, i := range m.ids {
		if i == string(n) {
			return true
		}
	}
	return false
}

func (m *mockValidatorSet) Size() int {
	return len(m.ids)
}
