// @ts-expect-error: decorator
@external("hypermode", "log")
declare function log(level: string, message: string): void;

// @ts-expect-error: decorator
@unsafe
@external("wasi_snapshot_preview1", "proc_exit")
declare function exit(rval: u32): void;

export default function modus_abort(
  message: string | null = null,
  fileName: string | null = null,
  lineNumber: u32 = 0,
  columnNumber: u32 = 0,
): void {
  let msg = message ? message : "abort()";
  if (fileName) {
    msg += ` at ${fileName}:${lineNumber}:${columnNumber}`;
  }

  log("fatal", msg);
  exit(255);
}
