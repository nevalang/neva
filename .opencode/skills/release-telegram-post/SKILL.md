# Release Telegram Post

Create a Telegram-ready release announcement in English for the Neva community group.

## Tools
- `Read` for release input context from the prompt.
- `Write` to save the final message file.

## Workflow
1. Read the release context from the prompt: repository, tag, release name, release URL, and release notes.
2. Read `.opencode/shared/release-marketing.md` for positioning, CTA, and tone guidance.
3. Read `README.md` and `.agent/notes/0.neva_vision_raw.md` for project language, ambition, and framing.
4. Produce a concise Telegram post in English with 2-3 short paragraphs:
   - Paragraph 1: what changed in this release.
   - Paragraph 2: why it matters for users/developers, with a strong but grounded product framing.
   - Paragraph 3: release link and a compact CTA block.
5. Make the post feel ambitious and clear, not dry and not hypey. Sell the direction, not just the technical properties.
6. Include CTA at the end (short, not noisy):
   - Open Collective support link.
   - GitHub star ask.
   - Share ask with open-source motivation.
7. Use Telegram HTML formatting only (`<b>`, `<i>`, `<code>`, `<a href=\"...\">`).
8. Keep the message easy to scan and under 1500 characters.
9. Save output to the exact file path requested by the caller.

## Output Rules
- Return only the final Telegram message body in the file.
- The file content is consumed as raw text in a JSON request (`jq --arg text ...`), so it must be directly sendable without edits.
- Do not add wrappers like "Here is your post", "Result:", "Output:", or follow-up suggestions.
- Do not add YAML frontmatter or surrounding quotes.
- Do not include code fences.
- Do not include unsupported HTML tags.
- Do not reduce the project description to a dry list of technical facts such as "statically typed" or "compiled" unless they support a stronger product point.
