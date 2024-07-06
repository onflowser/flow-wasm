import {GoPrompter} from "@/go-interfaces";

export class NullPrompter implements GoPrompter {
    shouldUpdateDependency(_contractName: string): boolean {
        return false;
    }

    addContractToDeployment(_networkName: string, accountsJson: string, _contractName: string): string {
        const accounts = JSON.parse(accountsJson);
        return accounts[0].Name;
    }

    addressPromptOrEmpty(_label: string): string {
        return ""
    }
}
