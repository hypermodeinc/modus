import * as host from "./hypermode";
import { JSON } from "json-as/assembly";

export function queryDQL<TData>(query: string): TData {
    const results = host.queryDQL(query);
    return JSON.parse<TData>(results);
}

export function queryGQL<TResults>(query: string): TResults {

    // Preferablly this would return `GQLResponse<TResults>`,
    // but we're blocked by https://github.com/JairusSW/as-json/issues/53

    const response = host.queryGQL(query);
    return JSON.parse<TResults>(response);
}

// // @ts-ignore
// @json
// export class GQLResponse<T> {
//     data!: T;
//     extensions: GQLExtensions | null = null;
// }

// @ts-ignore
@json
export class GQLExtensions {
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
