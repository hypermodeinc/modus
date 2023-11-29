// @ts-ignore
@json
export class DQLResponse<T> {
    data!: T;
    // TODO: errors, extensions
}

// @ts-ignore
@json
export class DQLMutationResponse {
    code!: string;
    message!: string;
    queries: string | null = null;
    // uids!: Map<string, string>; // Doesn't work :(
    uids!: Hack;
}

// TODO: Remove the hack and use the map instead.
// See https://github.com/JairusSW/as-json/issues/56

// @ts-ignore
@json
class Hack {
    x: string | null = null;

    get(key: string): string | null {
        if (key == "x") {
            return this.x;
        }
        
        return null;
    }
}
