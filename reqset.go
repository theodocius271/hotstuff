package hotstuff

import (
	"container/list"
	"crypto/sha256"
	"encoding/hex"
	"sync"

	"github.com/hyperledger/fabric/protos/utils"
	pb "github.com/theodocius271/hotstuff/proto"
)

type ReqSet interface {
	Add(cmds ...*pb.Transaction)
	Remove(cmds ...*pb.Transaction)
	GetFirst(n int) []*pb.Transaction
	IsProposed(cmd *pb.Transaction) bool
	MarkProposed(cmds ...*pb.Transaction)
	UnMark(cmds ...*pb.Transaction)
}

type cmdElement struct {
	cmd      *pb.Transaction
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

// Using Sha-256 to mark unique Transaction
func getTransactionKey(tx *pb.Transaction) string {
	hash := sha256.Sum256(utils.MarshalOrPanic(tx.Envelope))
	return hex.EncodeToString(hash[:])
}

// Add add cmds to the list and set, duplicate one will be ignored
func (c *cmdSetImpl) Add(cmds ...*pb.Transaction) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}

		key := getTransactionKey(cmd)
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
func (c *cmdSetImpl) Remove(cmds ...*pb.Transaction) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}

		key := getTransactionKey(cmd)
		if e, ok := c.set[key]; ok {
			c.order.Remove(e)
			delete(c.set, key)
		}
	}
}

// GetFirst return the top n unused commands from the list
func (c *cmdSetImpl) GetFirst(n int) []*pb.Transaction {
	c.lock.Lock()
	defer c.lock.Unlock()

	if len(c.set) == 0 {
		return nil
	}
	cmds := make([]*pb.Transaction, 0, n)
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

func (c *cmdSetImpl) IsProposed(cmd *pb.Transaction) bool {
	if cmd == nil {
		return false
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	key := getTransactionKey(cmd)
	if e, ok := c.set[key]; ok {
		return e.Value.(*cmdElement).proposed
	}
	return false
}

// MarkProposed will mark the given commands as proposed and move them to the back of the queue
func (c *cmdSetImpl) MarkProposed(cmds ...*pb.Transaction) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}

		key := getTransactionKey(cmd)
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

func (c *cmdSetImpl) UnMark(cmds ...*pb.Transaction) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}

		key := getTransactionKey(cmd)
		if e, ok := c.set[key]; ok {
			e.Value.(*cmdElement).proposed = false
			c.order.MoveToFront(e)
		}
	}
}
