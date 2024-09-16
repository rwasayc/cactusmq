package base

import (
	"testing"
)

type tmop int

const (
	opStore  tmop = 1
	opLoad   tmop = 2
	opDelete tmop = 3
	opRange  tmop = 4
	opLen    tmop = 5
)

type testAction[K comparable, V any] struct {
	op tmop
	k  K
	v  V
}

type testExpect[K comparable, V any] struct {
	op          tmop
	checkKey    K
	expectVal   any
	expectExist bool
}

// TestMapBasicOperations tests basic operations of the map.
func TestMapBasicOperations(t *testing.T) {
	caseList := []struct {
		name    string
		actions []testAction[int, string]
		expect  []testExpect[int, string]
	}{
		{
			name: "nil map",
			expect: []testExpect[int, string]{
				{
					op:        opLen,
					expectVal: 0,
				},
				{
					op:          opLoad,
					checkKey:    1,
					expectExist: false,
				},
				{
					op:          opLoad,
					checkKey:    2,
					expectExist: false,
				},
				{
					op:          opLoad,
					checkKey:    -1,
					expectExist: false,
				},
				{
					op:          opLoad,
					checkKey:    9990,
					expectExist: false,
				},
			},
		},
		{
			name: "map with values",
			actions: []testAction[int, string]{
				{op: opStore, k: 1, v: "value1"},
				{op: opStore, k: 2, v: "value2"},
				{op: opStore, k: 3, v: "value3"},
			},
			expect: []testExpect[int, string]{
				{
					op:        opLen,
					expectVal: 3,
				},
				{
					op:          opLoad,
					checkKey:    1,
					expectVal:   "value1",
					expectExist: true,
				},
				{
					op:          opLoad,
					checkKey:    2,
					expectVal:   "value2",
					expectExist: true,
				},
				{
					op:          opLoad,
					checkKey:    3,
					expectVal:   "value3",
					expectExist: true,
				},
			},
		},
		{
			name: "map with deletions",
			actions: []testAction[int, string]{
				{op: opStore, k: 1, v: "value1"},
				{op: opStore, k: 2, v: "value2"},
				{op: opStore, k: 3, v: "value3"},
				{op: opDelete, k: 2},
			},
			expect: []testExpect[int, string]{
				{
					op:        opLen,
					expectVal: 2,
				},
				{
					op:          opLoad,
					checkKey:    1,
					expectVal:   "value1",
					expectExist: true,
				},
				{
					op:          opLoad,
					checkKey:    2,
					expectExist: false,
				},
				{
					op:          opLoad,
					checkKey:    3,
					expectVal:   "value3",
					expectExist: true,
				},
			},
		},
	}

	for _, tc := range caseList {

		m := NewSyncMap[int, string]()
		t.Logf("Test case [%s]\n", tc.name)
		t.Logf("\tExecuting actions...\n")
		for _, act := range tc.actions {
			switch act.op {
			case opDelete:
				m.Delete(act.k)
			case opStore:
				m.Store(act.k, act.v)
			}
		}
		t.Logf("\tChecking expectations...\n")

		for idx, exp := range tc.expect {
			t.Logf("\t\tChecking result [%d]", idx)
			switch exp.op {
			case opLen:
				if m.Len() != exp.expectVal {
					t.Errorf("\t\t\t[Failed] m.Len()->%d  exp.expectVal->%d", m.Len(), exp.expectVal)
					t.FailNow()
				}
			case opLoad:
				val, exist := m.Load(exp.checkKey)
				if exist != exp.expectExist {
					t.Errorf("\t\t\t[Failed] exist->%v != exp.expectExist->%v", exist, exp.expectExist)
					t.FailNow()
				}
				if exist {
					if val != exp.expectVal {
						t.Errorf("\t\t\t[Failed] val->%v != exp.expectVal->%v", val, exp.expectVal)
						t.FailNow()
					}
				}
			}
		}
	}
}
