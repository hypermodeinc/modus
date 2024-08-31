class TestClass1 {
    a!: bool;
}

class TestClass2 {
    a!: bool;
    b!: isize;
}

class TestClass3 {
    a!: bool;
    b!: isize;
    c!: string;
}

class TestClass4 {
    a!: bool;
    b!: isize;
    c!: string | null;
}

class TestClass5 {
    a!: bool;
    b!: TestClass3;
}

class TestRecursiveClass {
    a!: bool;
    b!: TestRecursiveClass | null;
}

const testClass1 = <TestClass1>{ a: true };
const testClass2 = <TestClass2>{ a: true, b: 123 };
const testClass3 = <TestClass3>{ a: true, b: 123, c: "abc" };
const testClass4 = <TestClass4>{ a: true, b: 123, c: "abc" };
const testClass4_withNull = <TestClass4>{ a: true, b: 123, c: null };
const testClass5 = <TestClass5>{ a: true, b: testClass3 };

const testRecursiveClass = <TestRecursiveClass>{ a: true };
testRecursiveClass.b = testRecursiveClass;


export function testClassInput1(o: TestClass1): void {
    assert(o.a == testClass1.a);
}

export function testClassOutput1(): TestClass1 {
    return testClass1;
}

export function testClassInput2(o: TestClass2): void {
    assert(o.a == testClass2.a);
    assert(o.b == testClass2.b);
}

export function testClassOutput2(): TestClass2 {
    return testClass2;
}

export function testClassInput3(o: TestClass3): void {
    assert(o.a == testClass3.a);
    assert(o.b == testClass3.b);
    assert(o.c == testClass3.c);
}

export function testClassOutput3(): TestClass3 {
    return testClass3;
}

export function testClassInput4_withNull(o: TestClass4): void {
    assert(o.a == testClass4_withNull.a);
    assert(o.b == testClass4_withNull.b);
    assert(o.c == testClass4_withNull.c);
}

export function testClassOutput4(): TestClass4 {
    return testClass4;
}

export function testClassInput4(o: TestClass4): void {
    assert(o.a == testClass4.a);
    assert(o.b == testClass4.b);
    assert(o.c == testClass4.c);
}

export function testClassOutput4_withNull(): TestClass4 {
    return testClass4_withNull;
}

export function testClassInput5(o: TestClass5): void {
    assert(o.a == testClass5.a);
    assert(o.b.a == testClass5.b.a);
    assert(o.b.b == testClass5.b.b);
    assert(o.b.c == testClass5.b.c);
}

export function testClassOutput5(): TestClass5 {
    return testClass5;
}

export function testRecursiveClassInput(o: TestRecursiveClass): void {
    assert(o.a == testRecursiveClass.a);
    assert(o.b == testRecursiveClass.b);
}

export function testRecursiveClassOutput(): TestRecursiveClass {
    return testRecursiveClass;
}

export function testNullableClassInput1(o: TestClass1 | null): void {
    assert(o != null);
    assert(o!.a == testClass1.a);
}

export function testNullableClassOutput1(): TestClass1 | null {
    return testClass1;
}

export function testNullableClassInput1_null(o: TestClass1 | null): void {
    assert(o == null);
}

export function testNullableClassOutput1_null(): TestClass1 | null {
    return null;
}
