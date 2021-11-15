package pool

import (
	"crypto/sha256"
	"encoding/hex"
)

type NodeID string

type Message struct {
	// Hash is the hash of the data
	Hash string

	// Arbitrary data to be pooled
	Data []byte

	// From is the sender of the message
	From NodeID
}

type ValidatorSet interface {
	Includes(NodeID) bool
	Size() int
}

type MessagePool struct {
	local        NodeID
	transport    Transport
	messages     map[string]*messageTally
	validatorSet ValidatorSet
}

func NewMessagePool(local NodeID, transport Transport) *MessagePool {
	pool := &MessagePool{
		local:     local,
		transport: transport,
		messages:  make(map[string]*messageTally),
	}
	return pool
}

func (m *MessagePool) GetReady() []*Message {
	res := []*Message{}

	for _, msg := range m.messages {
		if msg.ready {
			res = append(res, nil) // include int he message tally the raw proposal so that we have easy access
		}
	}

	// Do we have to remove the message from the pool now,
	// or we acknoledge later on that it has been included? This part does not sound deterministic, we should if possible at any time include everything
	return res
}

func (m *MessagePool) Add(msg *Message) {
	// gossip
	// m.transport.Gossip(msg)

	// add to pool
	m.addImpl(msg)
}

func (m *MessagePool) addImpl(msg *Message) {
	if !m.validatorSet.Includes(msg.From) {
		return
	}
	if msg.Hash == "" {
		// hash the message
		hashRaw := sha256.Sum256(msg.Data)
		msg.Hash = hex.EncodeToString(hashRaw[:])
	}
	tally, ok := m.messages[msg.Hash]
	if !ok {
		tally = newMessageTally(msg.Data)
		m.messages[msg.Hash] = tally
	}
	count := tally.addMsg(msg)
	if count > m.validatorSet.Size()/2 { // mock value (depending on validatorset)
		tally.ready = true
	}
}

func (m *MessagePool) Reset(validatorSet ValidatorSet) {
	m.validatorSet = validatorSet

	// remove the values from the pool but reschedule the ones we have seen
	reschedule := []*Message{}
	for id, tally := range m.messages {
		if msg, ok := tally.hasLocal(m.local); ok {
			reschedule = append(reschedule, msg)
		}
		delete(m.messages, id)
	}

	// Should we send this if we are not a validator anymore? TODO

	// send again the rescheduled messages
	for _, msg := range reschedule {
		m.Add(msg)
	}
}

type messageTally struct {
	// tally of seen messages
	tally map[NodeID]*Message

	// arbitrary bytes of the proposal
	proposal []byte

	// ready selects whether the message is valid
	ready bool
}

func newMessageTally(proposal []byte) *messageTally {
	return &messageTally{
		tally:    map[NodeID]*Message{},
		proposal: proposal,
	}
}

func (m *messageTally) addMsg(msg *Message) int {
	if _, ok := m.tally[msg.From]; !ok {
		m.tally[msg.From] = msg
	}
	return len(m.tally)
}

func (m *messageTally) hasLocal(local NodeID) (*Message, bool) {
	msg, ok := m.tally[local]
	return msg, ok
}

// 1. validate the message in the pool? i.e. it decodes to an Ethereum event.
// 2. message pool is the one connected with the Transport
// 3. signing is done already in the gossip protocol.
// 4. being able to reset the pool.
// 5. add storage to save the data while there are reorgs.
//		- Also useful to get stats of what is being send and when.
// 		- This should only be valid for local ones.
// 6. we have to store messages even if we do not have it yet.
//		- This should only be reset at each epoch.
// 7. Connected with a ValidatorSet to know if it is a valid message?
// 8. We have two types of messages, we need to abstract that (arbitrary bytes?)
// 9. If the validator set changes the threshold for message valid changes too.
// 10. Even when there are slashings do we have to reset the message tally?
// 		- that moment would be good to pass the validator set.

type Transport interface {
	Gossip(msg *Message)
}
