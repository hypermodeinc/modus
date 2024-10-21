#!/usr/bin/env node --no-warnings=ExperimentalWarning
import { execute } from "@oclif/core";

await execute({ dir: import.meta.url });
