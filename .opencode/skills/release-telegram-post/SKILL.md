# Release Telegram Post

Create a Telegram-ready release announcement in English for the Neva community group.

## Tools
- `Read` for release input context from the prompt.
- `Write` to save the final message file.

## Workflow
1. Read the release context from the prompt: repository, tag, release name, release URL, and release notes.
2. Read `docs/marketing/shared-messaging.md` and reuse the project blurb and CTA lines.
3. Produce a concise Telegram post in English with 2-3 short paragraphs:
   - Paragraph 1: what changed in this release.
   - Paragraph 2: why it matters for users/developers.
   - Paragraph 3: release link and a compact CTA block.
4. Include CTA at the end (short, not noisy):
   - Open Collective support link.
   - GitHub star ask.
   - Share ask with open-source motivation.
5. Use Telegram HTML formatting only (`<b>`, `<i>`, `<code>`, `<a href=\"...\">`).
6. Keep the message easy to scan and under 1500 characters.
7. Save output to the exact file path requested by the caller.

## Output Rules
- Return only the final Telegram message body in the file.
- The file content is consumed as raw text in a JSON request (`jq --arg text ...`), so it must be directly sendable without edits.
- Do not add wrappers like "Here is your post", "Result:", "Output:", or follow-up suggestions.
- Do not add YAML frontmatter or surrounding quotes.
- Do not include code fences.
- Do not include unsupported HTML tags.
