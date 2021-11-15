package pos

import pool "github.com/0xPolygon/polygon-sdk/consensus/pos/message_pool"

type Consensus struct {
	stateSyncPool *pool.MessagePool
}

func (c *Consensus) AddStateSync() {
	// this message is received from the bridge?
	// would it be good to also do this arbitrary?
	c.stateSyncPool.Add(&pool.Message{})
}

/*
1. API to connect to message pools to notify when new things are up
2. APIs are event specific with the typing (not arbitrary as pool)
3. How do we decide when we seal?
	- How does Heimdall does it? Each 10 seconds?
	- Here we would have more control rather than using seconds (i.e. each 10 blocks)
	- That differs between SDK and Bor. One can replace all of this with this module
		while the other still has to think about execution layer.
		Though we could combine both in a single sealing module
*/
