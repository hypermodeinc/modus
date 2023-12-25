import * as host from "./hypermode";
import { DQLResponse, DQLMutationResponse } from "./dqltypes";
import { GQLResponse } from "./gqltypes";
import { JSON } from "json-as";

export abstract class dql {
    public static mutate(query: string): DQLResponse<DQLMutationResponse> {
        return this.execute<DQLMutationResponse>(query, true);
    }

    public static query<TData>(query: string): DQLResponse<TData> {
        return this.execute<TData>(query, false);
    }

    private static execute<TData>(query: string, isMutation: bool): DQLResponse<TData> {
        let response = host.executeDQL(query, isMutation);
        response = response.replace("@groupby","groupby");
        console.log(response);
        return JSON.parse<DQLResponse<TData>>(response);
    }
}

export abstract class graphql {
    static execute<TData>(statement: string): GQLResponse<TData> {
        const response = host.executeGQL(statement);
        console.log(`graphql ${response}`)
        return JSON.parse<GQLResponse<TData>>(response);
    }
}

export function invokeClassifier(modelId: string, text: string): ClassificationResult {

    // Preferablly this would return `GQLResponse<TResults>`,
    // but we're blocked by https://github.com/JairusSW/as-json/issues/53
    console.log("Invoking invokeClassifier")
    const response = host.invokeClassifier(modelId, text);
    console.log(`response ${response}`)
    const result = JSON.parse<ClassificationResult>(response)

    return result;
}
// @ts-ignore
@json
export class ClassificationProbability { // must be defined in the library
  label!: string;
  probability!: f32;
};

// @ts-ignore
@json
export class ClassificationResult { // must be defined in the library
  label: string = "";
  confidence: f32 = 0.0;
  probabilities!: ClassificationProbability[]
};
