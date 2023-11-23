import * as host from "./hypermode";
import { JSON } from "json-as";

export function queryDQL<TData>(query: string): TData {
    const results = host.queryDQL(query);
    return JSON.parse<TData>(results);
}

export function queryGQL<TData>(query: string): GQLResponse<TData> {
    const response = host.queryGQL(query);
    return JSON.parse<GQLResponse<TData>>(response);
}

// @ts-ignore
@json
class GQLResponse<T> {
    data!: T;
    extensions: GQLExtensions | null = null;
}

// @ts-ignore
@json
class GQLExtensions {
    touched_uids: u32 = 0;
    tracing!: GQLTracing;
}

// @ts-ignore
@json
class GQLTracing {
    version!: u32;
    startTime!: Date
    endTime!: Date
    duration!: u32;
    execution: GQLExecution | null = null;
}

// @ts-ignore
@json
class GQLExecution {
    resolvers!: GQLResolver[];
}

// @ts-ignore
@json
class GQLResolver {
    path!: string[];
    parentType!: string;
    fieldName!: string;
    returnType!: string;
    startOffset!: u32;
    duration!: u32;
    dgraph!: GQLDgraph[];
}

// @ts-ignore
@json
class GQLDgraph {
    label!: string;
    startOffset!: u32;
    duration!: u32;
}
