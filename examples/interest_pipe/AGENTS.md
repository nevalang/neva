# AGENTS.md

This example models account state updates in three sequential steps:

1. initial deposit
2. interest with a rate change
3. interest with a balance change

Expected printed values:

- `5000` for `5 * 1000 * (1-0)` after initial deposit
- `5000` for `5 * 1000 * (2-1)` before rate update applies
- `6000` for `6 * 1000 * (3-2)` after rate update
