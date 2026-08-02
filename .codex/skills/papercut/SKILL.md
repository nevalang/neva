---
name: "papercut"
description: "Capture or triage small Neva repository and engineering-harness frictions without derailing feature work."
---

# Papercut

Use this skill when a person or agent encounters minor friction while doing
other Neva work: a flaky command, stale cache, misleading error, obsolete
script, or documentation that contradicts observed behavior.

The journal is [`.codex/papercuts.md`](../../papercuts.md).
It is not an issue tracker. Product, language, compiler, and runtime defects
belong in their owning GitHub issue or change; already-tracked work stays there.

## Capture

1. Keep working on the original task unless the friction blocks it.
2. Add one entry directly below the journal marker, newest first:

   ```text
   ## YYYY-MM-DD HH:MM TZ — author or model

   What I was doing -> what got in the way. Include the relevant command or
   error, workaround, and likely cause or fix when known.
   ```

3. Record only evidence already observed. Never put credentials, access tokens,
   or private data in an entry.

## Triage

1. Read the full journal and group duplicate symptoms.
2. Reproduce or otherwise verify each candidate before assigning a cause.
3. Convert a meaningful item into the smallest durable outcome: a test,
   automation, documentation update, or GitHub issue.
4. Link that outcome in the pull request or issue, then remove the triaged
   entry from the journal. Leave unverified entries intact.

Do not silently broaden the current task to implement a separate fix. Ask for
direction when the remediation is material and not already in scope.
