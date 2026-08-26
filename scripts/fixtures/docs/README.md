Negative fixtures for `scripts/verify-documentation.mjs`. Each case in `cases.json` must fail the named check. They are not published documentation.

`generated-missing-pages` is an llms.txt that includes `/start` and `/guides/api-conventions` but omits `/products/beans`. It must fail under `--require-generated`.

