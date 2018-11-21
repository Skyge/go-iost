package vm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bitly/go-simplejson"

	"github.com/iost-official/go-iost/common"
	"github.com/iost-official/go-iost/core/contract"
	"github.com/iost-official/go-iost/ilog"
	"github.com/iost-official/go-iost/vm/host"
	"github.com/iost-official/go-iost/vm/native"
	"github.com/iost-official/go-iost/vm/v8vm"
)

var (
	errGasPriceIllegal = errors.New("gas price illegal")
	errGasLimitIllegal = errors.New("gas limit illegal")
)

// Monitor ...
type Monitor struct {
	vms map[string]VM
}

// NewMonitor ...
func NewMonitor() *Monitor {
	m := &Monitor{
		vms: make(map[string]VM),
	}
	jsvm := Factory("javascript")
	m.vms["javascript"] = jsvm
	return m
}

func (m *Monitor) prepareContract(h *host.Host, contractName, api, jarg string) (c *contract.Contract, abi *contract.ABI, args []interface{}, err error) {
	var cid string
	if h.IsDomain(contractName) {
		cid = h.ContractID(contractName)
	} else {
		cid = contractName
	}

	c = h.DB().Contract(cid)
	if c == nil {
		return nil, nil, nil, fmt.Errorf("contract %s not found", cid)
	}

	abi = c.ABI(api)

	if abi == nil {
		return nil, nil, nil, fmt.Errorf("abi %s not found", api)
	}

	args, err = unmarshalArgs(abi, jarg)

	return
}

// Call ...
// nolint
func (m *Monitor) Call(h *host.Host, contractName, api string, jarg string) (rtn []interface{}, cost contract.Cost, err error) {

	c, abi, args, err := m.prepareContract(h, contractName, api, jarg)

	if err != nil {
		return nil, host.Costs["GetCost"], fmt.Errorf("prepare contract: %v", err)
	}

	h.PushCtx()

	h.Context().Set("contract_name", c.ID)
	h.Context().Set("abi_name", api)

	vm, ok := m.vms[c.Info.Lang]
	if !ok {
		vm = Factory(c.Info.Lang)
		m.vms[c.Info.Lang] = vm
		err := m.vms[c.Info.Lang].Init()
		if err != nil {
			panic(err)
		}
	}
	// check amount limit
	signerList := map[string]int{}
	if h.Context().Value("signer_list") != nil {
		signerList = h.Context().Value("signer_list").(map[string]int)
	}
	amountLimit := abi.AmountLimit
	if amountLimit == nil {
		amountLimit = []*contract.Amount{}
	}
	var userAmountLimit []*contract.Amount
	if h.Context().Value("amount_limit") != nil {
		userAmountLimit = h.Context().Value("amount_limit").([]*contract.Amount)
	}
	var fixedAmountLimit []contract.FixedAmount
	beforeBalance := make(map[string][]int64)
	cost = contract.Cost0()

	// only check amount limit when executing action, not system call
	if h.Context().Value("stack_height") == 1 {
		cost0 := host.CommonOpCost(len(signerList) * len(amountLimit))
		cost.AddAssign(cost0)
		for _, limit := range amountLimit {
			decimal := h.DB().Decimal(limit.Token)
			fixedAmount, err := common.NewFixed(limit.Val, decimal)
			if err == nil {
				fixedAmountLimit = append(fixedAmountLimit, contract.FixedAmount{limit.Token, fixedAmount})
			}
		}
		for _, limit := range userAmountLimit {
			decimal := h.DB().Decimal(limit.Token)
			fixedAmount, err := common.NewFixed(limit.Val, decimal)
			if err == nil {
				fixedAmountLimit = append(fixedAmountLimit, contract.FixedAmount{limit.Token, fixedAmount})
			}
		}
		for acc := range signerList {
			beforeBalance[acc] = []int64{}
			for _, limit := range fixedAmountLimit {
				beforeBalance[acc] = append(beforeBalance[acc], h.DB().TokenBalance(limit.Token, acc))
			}
		}
	}

	rtn, cost0, err := vm.LoadAndCall(h, c, api, args...)
	cost.AddAssign(cost0)

	//payment, ok := h.Context().GValue("abi_payment").(int)
	//if !ok {
	//	payment = int(abi.Payment)
	//}
	//var gasPrice = h.Context().Value("gas_price").(int64)

	//if payment == 1 &&
	//	abi.GasPrice > gasPrice &&
	//	!{
	//	b := h.DB().TokenBalance("iost",host.ContractGasPrefix + contractName)
	//	if b > gasPriceCost.ToGas() {
	//		h.PayCost(cost, host.ContractGasPrefix+contractName)
	//		cost = contract.Cost0()
	//	}
	//}

	// check amount limit
	if h.Context().Value("stack_height") == 1 {
		for acc := range signerList {
			for i, limit := range fixedAmountLimit {
				afterBalance := h.DB().TokenBalance(limit.Token, acc)
				delta := common.Fixed{
					Value:   beforeBalance[acc][i] - afterBalance,
					Decimal: fixedAmountLimit[i].Val.Decimal,
				}
				if delta.Value > fixedAmountLimit[i].Val.Value {
					err = errors.New(fmt.Sprintf("token %s exceed amountLimit in abi. limit %s, need %s",
						limit.Token,
						fixedAmountLimit[i].Val.ToString(),
						delta.ToString()))
					return nil, cost, err
				}
			}
		}
	}

	h.PopCtx()

	return
}

// Compile ...
func (m *Monitor) Compile(con *contract.Contract) (string, error) {
	switch con.Info.Lang {
	case "native":
		return "", nil
	case "javascript":
		jsvm, _ := m.vms["javascript"]
		return jsvm.Compile(con)
	}
	return "", errors.New("vm unsupported")
}

// Factory ...
func Factory(lang string) VM {
	switch lang {
	case "native":
		vm := native.Impl{}
		vm.Init()
		return &vm
	case "javascript":
		vm := v8.NewVMPool(10, 200)
		vm.Init()
		//vm.SetJSPath(jsPath)
		return vm
	}
	return nil
}

func unmarshalArgs(abi *contract.ABI, data string) ([]interface{}, error) {
	if strings.HasSuffix(data, ",]") {
		data = data[:len(data)-2] + "]"
	}
	js, err := simplejson.NewJson([]byte(data))
	if err != nil {
		return nil, fmt.Errorf("error in data: %v, %v", err, data)
	}

	rtn := make([]interface{}, 0)
	arr, err := js.Array()
	if err != nil {
		ilog.Error(js.EncodePretty())
		return nil, err
	}

	if len(arr) != len(abi.Args) {
		return nil, fmt.Errorf("args length unmatched to abi %v. need %v, got %v", abi.Name, len(abi.Args), len(arr))
	}
	for i := range arr {
		switch abi.Args[i] {
		case "string":
			s, err := js.GetIndex(i).String()
			if err != nil {
				return nil, err
			}
			rtn = append(rtn, s)
		case "bool":
			s, err := js.GetIndex(i).Bool()
			if err != nil {
				return nil, err
			}
			rtn = append(rtn, s)
		case "number":
			s, err := js.GetIndex(i).Int64()
			if err != nil {
				return nil, err
			}
			rtn = append(rtn, s)
		case "json":
			s, err := js.GetIndex(i).Encode()
			if err != nil {
				return nil, err
			}
			// make sure s is a valid json
			_, err = simplejson.NewJson(s)
			if err != nil {
				ilog.Error(string(s))
				return nil, err
			}
			rtn = append(rtn, s)
		}
	}

	return rtn, nil
}
