package js

import (
	"encoding/json"
	"github.com/onflow/flowkit/v2/accounts"
	"github.com/onflow/flowkit/v2/deps"
	"syscall/js"
)

type Prompter struct {
	target js.Value
}

func NewPrompter(target js.Value) *Prompter {
	return &Prompter{target: target}
}

func (p *Prompter) ShouldUpdateDependency(contractName string) bool {
	result := p.target.Call("shouldUpdateDependency", contractName)
	return result.Bool()
}

func (p *Prompter) AddContractToDeployment(networkName string, accounts accounts.Accounts, contractName string) *deps.DeploymentData {
	accountsJson, err := json.Marshal(accounts)

	if err != nil {
		panic(err)
	}

	result := p.target.Call("addContractToDeployment", networkName, string(accountsJson), contractName)
	return &deps.DeploymentData{
		Network:   networkName,
		Account:   result.String(),
		Contracts: []string{contractName},
	}
}

func (p *Prompter) AddressPromptOrEmpty(label string, validate deps.InputValidator) string {
	var result string
	for {
		result = p.target.Call("addressPromptOrEmpty", label).String()
		err := validate(result)
		// TODO: Let the callee validate the input
		if err == nil {
			break
		}
	}

	return result
}

var _ deps.Prompter = &Prompter{}
