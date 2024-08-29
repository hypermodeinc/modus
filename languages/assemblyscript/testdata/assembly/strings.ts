// "Hello World" in Japanese
const testString = "こんにちは、世界"

export function testStringInput(s: string): void {
    assert(s == testString);
}

export function testStringOutput(): string {
    return testString;
}
