import {GoPrompter, PrompterAccount} from "@/go-interfaces";

export class DefaultPrompter implements GoPrompter {
    shouldUpdateDependency(contractName: string): boolean {
        const result = window.prompt(`The latest version of ${contractName} is different from the one you have locally. Do you want to update it? (true/false)`, "false");
        return result === "true";
    }

    addContractToDeployment(networkName: string, accounts: PrompterAccount[], contractName: string): string {
        const result = window.prompt(`Choose an account to deploy ${contractName} to on ${networkName} (${accounts.map(account => `'${account.Name}'`).join(", ")} or 'none' to skip)`);

        if (!result) {
            throw new Error("No account selected")
        }

        return result;
    }

    addressPromptOrEmpty(label: string): string {
        const result = window.prompt(label);

        if (!result) {
            throw new Error("No address entered")
        }

        return result;
    }

}
