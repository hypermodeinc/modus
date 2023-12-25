// This file should only export functions from the "hypermode" host module.

export declare function executeDQL(statement: string, isMutation: bool): string;
export declare function executeGQL(statement: string): string;
export declare function invokeClassifier(modelId: string, sentence: string): string;
