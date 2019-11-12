package space

import (
	"fmt"
	"testing"
)

const (
	TestEntityMask_BaseState   = 11
	TestEntityMask_ActionState = 12
)

type TestState struct {
	EntityState

	BaseState   byte
	ActionState byte
}

func newTestState() IEntityState {

	//state := &TestState{}
	// state.Init(state, newTestState)

	// state.Bind("BaseState", TestEntityMask_BaseState)
	// state.Bind("ActionState", TestEntityMask_ActionState)

	//return state
	return nil
}

func (s *TestState) String() string {
	return fmt.Sprint("time stamp ", s.GetTimeStamp(), " pos ", s.GetPos(), " rota ", s.GetRota(), " state ", s.BaseState, s.ActionState, " isDirty ", s.IsDirty(), "  isModify ", s.IsModify())
}

// func BenchmarkEntityStateBind1(b *testing.B) {
// 	state := &TestState{}
// 	state.Init(state, newTestState)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		state.Bind("Pos", TestEntityMask_BaseState)
// 	}
// }

// func BenchmarkEntityStateBind2(b *testing.B) {
// 	state := &TestState{}
// 	state.Init(state, newTestState)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		state.Bind("Pos", EntityStateMask_Pos_X)
// 	}
// }

// func BenchmarkEntityStateClone(b *testing.B) {
// 	state := newTestState()

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		state.Clone()
// 	}
// }

// func BenchmarkEntityStateDelta(b *testing.B) {
// 	a := newTestState().(*TestState)

// 	a.SetPos(linmath.NewVector3(1, 2, 3))
// 	a.SetRota(linmath.NewVector3(11, 22, 33))

// 	a.SetTimeStamp(3)
// 	a.BaseState = 1
// 	a.ActionState = 2

// 	bb := a.Clone().(*TestState)

// 	bb.SetPos(linmath.NewVector3(21, 22, 21))
// 	bb.BaseState = 21

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		a.Delta(bb)
// 	}
// }

// func BenchmarkEntityStateCombine(b *testing.B) {
// 	a := newTestState().(*TestState)

// 	a.SetPos(linmath.NewVector3(1, 2, 3))
// 	a.SetRota(linmath.NewVector3(11, 22, 33))

// 	a.SetTimeStamp(3)
// 	a.BaseState = 1
// 	a.ActionState = 2

// 	bb := a.Clone().(*TestState)

// 	bb.SetPos(linmath.NewVector3(21, 22, 21))
// 	bb.BaseState = 21

// 	d, _ := a.Delta(bb)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		a.Combine(d)
// 	}

// }

func TestEntityStates(t *testing.T) {
	// es := NewEntityStates()

	// state1 := newTestState()
	// state1.SetPos(linmath.NewVector3(1, 1, 1))
	// state1.SetRota(linmath.NewVector3(1, 1, 1))
	// state1.SetTimeStamp(1)
	// es.addEntityState(state1)
	// state1.SetDirty(false)

	// state5 := newTestState()
	// state5.SetPos(linmath.NewVector3(5, 5, 5))
	// state5.SetRota(linmath.NewVector3(5, 5, 5))
	// state5.SetTimeStamp(5)
	// es.addEntityState(state5)
	// state5.SetDirty(false)

	// state10 := newTestState()
	// state10.SetPos(linmath.NewVector3(10, 10, 10))
	// state10.SetRota(linmath.NewVector3(10, 10, 10))
	// state10.SetTimeStamp(10)
	// es.addEntityState(state10)
	// state10.SetDirty(false)

	// ns := es.GetHistoryState(2)
	// fmt.Println(ns)
}

func TestSimple(t *testing.T) {
	fmt.Println("hello")

	fmt.Println("world")

	fmt.Println("asdasd")
}

// func TestEntityState(t *testing.T) {

// 	a := newTestState().(*TestState)

// 	a.SetPos(linmath.NewVector3(1, 2, 3))
// 	a.SetRota(linmath.NewVector3(11, 22, 33))

// 	a.SetTimeStamp(3)
// 	a.BaseState = 1
// 	a.ActionState = 2

// 	b := a.Clone().(*TestState)

// 	b.SetPos(linmath.NewVector3(21, 22, 21))
// 	b.BaseState = 21

// 	d := a.Delta(b)

// 	fmt.Println(b)

// 	a.Combine(d)

// 	fmt.Println(a)

// }

func TestEntityStateEvents(t *testing.T) {
	state := &EntityState{}
	if err := state.MarshalEvent("testEvent", true, uint64(1111)); err != nil {
		fmt.Println(err)
	}

	fmt.Println(state.Events)

	fmt.Println(state.UnmarshalEvent())
}
