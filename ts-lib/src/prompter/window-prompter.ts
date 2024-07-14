import { GoPrompter } from "@/go-interfaces";

export class WindowPrompter implements GoPrompter {
  shouldUpdateDependency(contractName: string): boolean {
    const result = window.prompt(
      `The latest version of ${contractName} is different from the one you have locally. Do you want to update it? (true/false)`,
      "false"
    );
    return result === "true";
  }

  addContractToDeployment(
    networkName: string,
    accountsJson: string,
    contractName: string
  ): string {
    const accounts = JSON.parse(accountsJson);
    const result = window.prompt(
      `Choose an account to deploy ${contractName} to on ${networkName} (${accounts.map((account: any) => `'${account.Name}'`).join(", ")} or 'none' to skip)`
    );

    if (!result) {
      throw new Error("No account selected");
    }

    return result;
  }

  addressPromptOrEmpty(label: string): string {
    const result = window.prompt(label);

    return result ?? "";
  }
}
