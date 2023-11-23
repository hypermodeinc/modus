import * as host from "./hypermode";
import { GQLResponse } from "./gqltypes";
import { JSON } from "json-as";

export function queryDQL<TData>(query: string): TData {
    const results = host.queryDQL(query);
    return JSON.parse<TData>(results);
}

export function queryGQL<TData>(query: string): GQLResponse<TData> {
    const response = host.queryGQL(query);
    return JSON.parse<GQLResponse<TData>>(response);
}
