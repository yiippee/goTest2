package main

import (
	"fmt"
	"github.com/smallnest/gofsm"
	"log"
)

const (
	EVENT_CATCH_ORDER   = "catch order"
	EVENT_MATCH_SKU     = "match sku"
	EVENT_MATCH_CARRIER = "match carrier"
	EVENT_UPLOAD_WMS    = "upload wms"
)

const (
	STATE_INIT          = "init"
	STATE_CATCH_ORDER   = "catch order"
	STATE_MATCH_SKU     = "match sku"
	STATE_MATCH_CARRIER = "match carrier"
	STATE_UPLOAD_WMS    = "upload wms"
)

type Package struct {
	ID           uint64
	CurrentState string
	States       []string
}

type PackageEventProcessor struct{}

func (p *PackageEventProcessor) OnExit(fromState string, args []interface{}) {
	t := args[0].(*Package)
	if t.CurrentState != fromState {
		panic(fmt.Errorf("订单 %v 的状态与期望的状态 %s 不一致，可能在状态机外被改变了", t, fromState))
	}

	log.Printf("订单 %d 从状态 %s 改变", t.ID, fromState)
}

func (p *PackageEventProcessor) Action(action string, fromState string, toState string, args []interface{}) error {
	fmt.Println("Action...")
	return nil
}

func (p *PackageEventProcessor) OnEnter(toState string, args []interface{}) {
	t := args[0].(*Package)
	t.CurrentState = toState
	t.States = append(t.States, toState)

	log.Printf("订单 %d 的状态改变为 %s ", t.ID, toState)
}

func (p *PackageEventProcessor) OnActionFailure(action string, fromState string, toState string, args []interface{}, err error) {
	t := args[0].(*Package)

	log.Printf("订单 %d 的状态从 %s to %s 改变失败， 原因: %v", t.ID, fromState, toState, err)
}

func initFSM() *fsm.StateMachine {
	delegate := &fsm.DefaultDelegate{P: &PackageEventProcessor{}}

	transitions := []fsm.Transition{
		{From: STATE_INIT, Event: EVENT_CATCH_ORDER, To: STATE_CATCH_ORDER, Action: "catch"},
		{From: STATE_CATCH_ORDER, Event: EVENT_MATCH_SKU, To: STATE_MATCH_SKU, Action: "matchSku"},
		{From: STATE_MATCH_SKU, Event: EVENT_MATCH_SKU, To: EVENT_MATCH_CARRIER, Action: "matchSku"},
	}

	return fsm.NewStateMachine(delegate, transitions...)
}

// 自身不具有状态的状态机，完全由对象本身维护状态，
// 状态机只提供转换动作
func main() {
	p := &Package{
		ID:           1,
		CurrentState: STATE_INIT,
		States:       []string{},
	}
	fsm := initFSM()

	// 根据当前的状态以及触发事件，找到匹配的规则
	err := fsm.Trigger(p.CurrentState, EVENT_CATCH_ORDER, p)
	if err != nil {
		log.Println("trigger err: %v", err)
	}

	err = fsm.Trigger(p.CurrentState, EVENT_CATCH_ORDER, p)
	if err != nil {
		log.Println("trigger err: %v", err)
	}

	err = fsm.Trigger(p.CurrentState, EVENT_MATCH_SKU, p)
	if err != nil {
		log.Println("trigger err: %v", err)
	}

	err = fsm.Trigger(p.CurrentState, EVENT_MATCH_CARRIER, p)
	if err != nil {
		log.Println("trigger err: %v", err)
	}

	fmt.Println(p, fsm)
}
