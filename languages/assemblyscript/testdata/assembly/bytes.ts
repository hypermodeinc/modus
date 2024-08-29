export function testArrayBufferInput(b: ArrayBuffer): void {
    const view = Uint8Array.wrap(b);
    assert(view.length == 4);
    assert(view[0] == 1);
    assert(view[1] == 2);
    assert(view[2] == 3);
    assert(view[3] == 4);
}

export function testArrayBufferOutput(): ArrayBuffer {
    const b = new ArrayBuffer(4);
    const view = Uint8Array.wrap(b);
    view[0] = 1;
    view[1] = 2;
    view[2] = 3;
    view[3] = 4;
    return b;
}
