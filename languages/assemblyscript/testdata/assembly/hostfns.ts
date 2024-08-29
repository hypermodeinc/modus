
// @ts-expect-error: decorator
@external("test", "add")
declare function hostAdd(a: i32, b: i32): i32;

// @ts-expect-error: decorator
@external("test", "echo")
declare function hostEcho(message: string): string;

export function add(a: i32, b: i32): i32 {
    return hostAdd(a, b);
}

export function echo(message: string): string {
    return hostEcho(message);
}
