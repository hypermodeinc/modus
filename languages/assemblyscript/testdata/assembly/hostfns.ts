
// @ts-expect-error: decorator
@external("test", "add")
declare function hostAdd(a: i32, b: i32): i32;

// @ts-expect-error: decorator
@external("test", "echo")
declare function hostEcho(message: string): string;

// @ts-expect-error: decorator
@external("test", "echoObject")
declare function hostEchoObject(obj: TestHostObject): TestHostObject;

class TestHostObject {
    a!: i32;
    b!: bool;
    c!: string;
}

export function add(a: i32, b: i32): i32 {
    return hostAdd(a, b);
}

export function echo(message: string): string {
    return hostEcho(message);
}

export function echoObject(obj: TestHostObject): TestHostObject {
    return hostEchoObject(obj);
}
