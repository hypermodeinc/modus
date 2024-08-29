export function testArrayInput_i32(arr: i32[]): void {
    assert(arr.length == 3);
    assert(arr[0] == 1);
    assert(arr[1] == 2);
    assert(arr[2] == 3);
}

export function testArrayInput_string(arr: string[]): void {
    assert(arr.length == 3);
    assert(arr[0] == "abc");
    assert(arr[1] == "def");
    assert(arr[2] == "ghi");
}

export function testArrayOutput_i32(): i32[] {
    return [1, 2, 3];
}

export function testArrayOutput_string(): string[] {
    return ["abc", "def", "ghi"];
}

class TestObject1 {
    constructor(public a: i32, public b: i32) { }
}

export function testArrayIteration(arr: TestObject1[]): void {
    for (let i = 0; i < arr.length; i++) {
        let obj = arr[i];
        console.log(`[${i}]: a=${obj.a}, b=${obj.b}`);
    }
}
