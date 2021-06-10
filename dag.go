package main

import "sync"

type edge struct {
	prevs []*edge
	next  []*edge
	tx    *Transaction
}

type DAG struct {
	tips  []*edge
	txs   map[TxRef]*edge
	mutex sync.Mutex
}

func NewDag() *DAG {
	return &DAG {
		tips: make([]*edge, 0),
		txs: make(map[TxRef]*edge),
		mutex: sync.Mutex{},
	}
}

func (d *DAG) Has(tx TxRef) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, ok := d.txs[tx]
	return ok
}

func (d *DAG) Tips() []TxRef {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	refs := make([]TxRef, len(d.tips))
	for i, e := range d.tips {
		refs[i] = e.tx.Reference()
	}
	return refs
}

func (d *DAG) Offer(tx *Transaction) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, ok := d.txs[tx.Reference()]; ok {
		return
	}

	newEdge := &edge{
		tx: tx,
	}

	// update next
	var removeTips = make(map[TxRef]bool)
	for _, ref := range tx.refs {
		refEdge := d.txs[ref]
		newEdge.prevs = append(newEdge.prevs, refEdge)
		refEdge.next = append(refEdge.next, newEdge)

		// remove tips
		for _, t := range d.tips {
			if t.tx.Reference() == ref {
				removeTips[ref] = true
			}
		}
	}
	j := 0
	for _, e := range d.tips {
		if _, ok := removeTips[e.tx.Reference()]; !ok {
			d.tips[j] = e
			j++
		}
	}
	d.tips = d.tips[:j]

	// add to all Txs
	d.txs[tx.Reference()] = newEdge

	// add tx to tips
	d.tips = append(d.tips, newEdge)
}

func (d *DAG) GetTransaction(txRef TxRef) *Transaction {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.txs[txRef].tx
}