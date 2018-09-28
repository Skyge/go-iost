package vm

import (
	"testing"

	"time"

	. "github.com/golang/mock/gomock"
	"github.com/iost-official/go-iost/core/contract"
	"github.com/iost-official/go-iost/vm/database"
	"github.com/iost-official/go-iost/vm/host"
)

func Init(t *testing.T) (*Monitor, *MockVM, *database.MockIMultiValue, *database.Visitor) {
	mc := NewController(t)
	defer mc.Finish()
	vm := NewMockVM(mc)
	db := database.NewMockIMultiValue(mc)
	vi := database.NewVisitor(100, db)
	pm := NewMonitor()
	pm.vms[""] = vm
	return pm, vm, db, vi
}

func TestMonitor_Call(t *testing.T) {
	monitor, vm, db, vi := Init(t)

	ctx := host.NewContext(nil)
	ctx.Set("gas_price", int64(1))

	h := host.NewHost(ctx, vi, monitor, nil)

	flag := false

	vm.EXPECT().LoadAndCall(Any(), Any(), Any(), Any()).DoAndReturn(func(h *host.Host, c *contract.Contract, api string, args ...string) (rtn []string, cost *contract.Cost, err error) {
		flag = true
		return []string{"world"}, cost, nil
	})

	c := contract.Contract{
		ID:   "Contract",
		Code: "codes",
		Info: &contract.Info{
			Lang:    "",
			Version: "1.0.0",
			Abi: []*contract.ABI{
				{
					Name:     "abi",
					Args:     []string{"string"},
					Payment:  0,
					GasPrice: int64(1000),
					Limit:    contract.NewCost(100, 100, 100),
				},
			},
		},
	}

	db.EXPECT().Get(Any(), Any()).DoAndReturn(func(table string, key string) (string, error) {
		return c.Encode(), nil
	})

	monitor.Call(h, "Contract", "abi", "[\"1\"]")

	if !flag {
		t.Fatal(flag)
	}
}

func TestMonitor_Context(t *testing.T) {
	monitor, vm, db, vi := Init(t)
	ctx := host.NewContext(nil)
	ctx.Set("gas_price", int64(1))

	h := host.NewHost(ctx, vi, monitor, nil)

	outerFlag := false
	innerFlag := false

	vm.EXPECT().LoadAndCall(Any(), Any(), "outer", Any()).DoAndReturn(func(h *host.Host, c *contract.Contract, api string, args ...interface{}) (rtn []string, cost *contract.Cost, err error) {
		outerFlag = true
		monitor.Call(h, "Contract", "inner", "[\"hello\"]")

		return []string{"world"}, cost, nil
	})

	vm.EXPECT().LoadAndCall(Any(), Any(), "inner", Any()).DoAndReturn(func(h *host.Host, c *contract.Contract, api string, args ...interface{}) (rtn []string, cost *contract.Cost, err error) {
		innerFlag = true
		return []string{"world"}, cost, nil
	})
	c := contract.Contract{
		ID:   "Contract",
		Code: "codes",
		Info: &contract.Info{
			Lang:    "",
			Version: "1.0.0",
			Abi: []*contract.ABI{
				{
					Name:     "outer",
					Args:     []string{"number"},
					Payment:  0,
					GasPrice: int64(1000),
					Limit:    contract.NewCost(100, 100, 100),
				},
				{
					Name:     "inner",
					Args:     []string{"string"},
					Payment:  0,
					GasPrice: int64(1000),
					Limit:    contract.NewCost(100, 100, 100),
				},
			},
		},
	}

	db.EXPECT().Get(Any(), Any()).DoAndReturn(func(table string, key string) (string, error) {
		return c.Encode(), nil
	})

	monitor.Call(h, "Contract", "outer", "[1]")

	if !outerFlag || !innerFlag {
		t.Fatal(outerFlag, innerFlag)
	}
}

func TestMonitor_HostCall(t *testing.T) {
	monitor, vm, db, vi := Init(t)
	staticMonitor = monitor

	ctx := host.NewContext(nil)

	ctx.Set("gas_price", int64(1))
	ctx.Set("stack_height", 1)
	ctx.Set("stack0", "test")

	h := host.NewHost(ctx, vi, monitor, nil)
	outerFlag := false
	innerFlag := false

	vm.EXPECT().LoadAndCall(Any(), Any(), "outer", Any()).DoAndReturn(func(h *host.Host, c *contract.Contract, api string, args ...interface{}) (rtn []string, cost *contract.Cost, err error) {
		outerFlag = true
		h.Call("Contract", "inner", "[\"hello\"]")

		return []string{"world"}, cost, nil
	})

	vm.EXPECT().LoadAndCall(Any(), Any(), "inner", Any()).DoAndReturn(func(h *host.Host, c *contract.Contract, api string, args ...interface{}) (rtn []string, cost *contract.Cost, err error) {
		innerFlag = true
		if h.Context().Value("abi_name") != "inner" {
			t.Fatal(h.Context())
		}

		return []string{"world"}, cost, nil
	})
	c := contract.Contract{
		ID:   "Contract",
		Code: "codes",
		Info: &contract.Info{
			Lang:    "",
			Version: "1.0.0",
			Abi: []*contract.ABI{
				{
					Name:     "outer",
					Args:     []string{"number"},
					Payment:  0,
					GasPrice: int64(1000),
					Limit:    contract.NewCost(100, 100, 100),
				},
				{
					Name:     "inner",
					Args:     []string{"string"},
					Payment:  0,
					GasPrice: int64(1000),
					Limit:    contract.NewCost(100, 100, 100),
				},
			},
		},
	}

	db.EXPECT().Get(Any(), Any()).DoAndReturn(func(table string, key string) (string, error) {
		return c.Encode(), nil
	})

	monitor.Call(h, "Contract", "outer", "[1]")

	if !outerFlag || !innerFlag {
		t.Fatal(outerFlag, innerFlag)
	}
}

func TestJSM(t *testing.T) {
	monitor, _, db, vi := Init(t)

	ctx := host.NewContext(nil)
	ctx.Set("gas_price", int64(1))
	ctx.GSet("gas_limit", int64(1000))

	h := host.NewHost(ctx, vi, monitor, nil)
	h.SetDeadline(time.Now().Add(time.Second))

	c := contract.Contract{
		ID: "Contract",
		Code: `
class Contract {
 init() {
  
 }
 hello() {
  return "world";
 }
}

module.exports = Contract;
`,
		Info: &contract.Info{
			Lang:    "javascript",
			Version: "1.0.0",
			Abi: []*contract.ABI{
				{
					Name:     "hello",
					Args:     []string{},
					Payment:  0,
					GasPrice: int64(1000),
					Limit:    contract.NewCost(100, 100, 100),
				},
			},
		},
	}

	db.EXPECT().Get(Any(), Any()).DoAndReturn(func(table string, key string) (string, error) {
		return c.Encode(), nil
	})

	rs, co, e := monitor.Call(h, "Contract", "hello", `[]`)
	if rs[0] != "world" {
		t.Fatal(rs, co, e)
	}

}