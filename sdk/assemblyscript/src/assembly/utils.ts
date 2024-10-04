export function resultIsInvalid<T>(result: T): bool {
  return changetype<usize>(result) == 0;
}
