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
        const response = host.executeDQL(query, isMutation);
        return JSON.parse<DQLResponse<TData>>(response);
    }
}

export abstract class graphql {
    static execute<TData>(statement: string): GQLResponse<TData> {
        const response = host.executeGQL(statement);
        return JSON.parse<GQLResponse<TData>>(response);
    }
}
