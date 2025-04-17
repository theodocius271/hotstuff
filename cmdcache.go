package hotstuff

import (
	"container/list"
	"crypto/sha256"
	"encoding/hex"
	"sync"

	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/utils"
)

type CmdSet interface {
	Add(cmds ...*common.Envelope)
	Remove(cmds ...*common.Envelope)
	GetFirst(n int) []*common.Envelope
	IsProposed(cmd *common.Envelope) bool
	MarkProposed(cmds ...*common.Envelope)
	UnMark(cmds ...*common.Envelope)
}

type cmdElement struct {
	cmd      *common.Envelope
	proposed bool
}

type cmdSetImpl struct {
	lock  sync.Mutex
	order list.List
	set   map[string]*list.Element
}

func NewCmdSet() *cmdSetImpl {
	c := &cmdSetImpl{
		set: make(map[string]*list.Element),
	}
	c.order.Init()
	return c
}

// Using Sha-256 to mark unique Envelope
func getEnvelopeKey(envelope *common.Envelope) string {
	hash := sha256.Sum256(utils.MarshalOrPanic(envelope))
	return hex.EncodeToString(hash[:])
}

// Add add cmds to the list and set, duplicate one will be ignored
func (c *cmdSetImpl) Add(cmds ...*common.Envelope) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}

		key := getEnvelopeKey(cmd)
		// avoid duplication
		if _, ok := c.set[key]; ok {
			continue
		}
		e := c.order.PushBack(&cmdElement{
			cmd:      cmd,
			proposed: false,
		})
		c.set[key] = e
	}
}

// Remove remove commands from set and list
func (c *cmdSetImpl) Remove(cmds ...*common.Envelope) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}

		key := getEnvelopeKey(cmd)
		if e, ok := c.set[key]; ok {
			c.order.Remove(e)
			delete(c.set, key)
		}
	}
}

// GetFirst return the top n unused commands from the list
func (c *cmdSetImpl) GetFirst(n int) []*common.Envelope {
	c.lock.Lock()
	defer c.lock.Unlock()

	if len(c.set) == 0 {
		return nil
	}
	cmds := make([]*common.Envelope, 0, n)
	i := 0
	// get the first element of list
	e := c.order.Front()
	for i < n {
		if e == nil {
			break
		}
		if cmd := e.Value.(*cmdElement); !cmd.proposed {
			cmds = append(cmds, cmd.cmd)
			i++
		}
		e = e.Next()
	}
	return cmds
}

func (c *cmdSetImpl) IsProposed(cmd *common.Envelope) bool {
	if cmd == nil {
		return false
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	key := getEnvelopeKey(cmd)
	if e, ok := c.set[key]; ok {
		return e.Value.(*cmdElement).proposed
	}
	return false
}

// MarkProposed will mark the given commands as proposed and move them to the back of the queue
func (c *cmdSetImpl) MarkProposed(cmds ...*common.Envelope) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}

		key := getEnvelopeKey(cmd)
		if e, ok := c.set[key]; ok {
			e.Value.(*cmdElement).proposed = true
			// Move to back so that it's not immediately deleted by a call to TrimToLen()
			c.order.MoveToBack(e)
		} else {
			// new cmd, store it to back
			e := c.order.PushBack(&cmdElement{cmd: cmd, proposed: true})
			c.set[key] = e
		}
	}
}

func (c *cmdSetImpl) UnMark(cmds ...*common.Envelope) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}

		key := getEnvelopeKey(cmd)
		if e, ok := c.set[key]; ok {
			e.Value.(*cmdElement).proposed = false
			c.order.MoveToFront(e)
		}
	}
}
