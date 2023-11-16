// @ts-ignore
@json
export class GQLResponse<T> {
    data!: T;
    extensions: GQLExtensions | null = null;
}

// @ts-ignore
@json
export class GQLExtensions {
    touched_uids: u32 = 0;
    tracing!: GQLTracing;
}

// TODO: Use Date instead of string when this is fixed:
// https://github.com/JairusSW/as-json/issues/54

// @ts-ignore
@json
class GQLTracing {
    version!: u32;
    startTime!: string; // Date
    endTime!: string;  // Date
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
