export function add(a: u32, b: u32): u32 {
  return a + b;
}

export function now(): i64 {
  return Date.now();
}

export function spin(duration: i64): i64 {
  const start = performance.now();

  let d = Date.now();
  while (Date.now() - d <= duration) {
    // do nothing
  }

  const end = performance.now();

  return i64(end - start);
}
