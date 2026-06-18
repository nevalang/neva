---
name: "release-telegram-post"
description: "Create a Telegram-ready Neva release announcement from a GitHub release payload. Use for official Neva Telegram release posts."
---

# Release Telegram Post

Create a Telegram-ready release announcement in English for the Neva community group.

## Workflow

1. Read the release context from the prompt: repository, tag, release name,
   release URL, and release notes.
2. Read `.codex/shared/release-marketing.md` for positioning, CTA, and tone
   guidance.
3. Produce a concise Telegram post in English using this high-level structure:
   - Line 1: `Neva vX.Y - Release Title`
   - Line 2: standalone release URL
   - Short intro paragraph, 1-2 sentences max; state whether there are
     language-level breaking changes when the notes make that clear
   - `Highlights:` label
   - 6-8 short numbered highlight lines using emoji numerals
   - Final CTA paragraph with repo star/share/Open Collective links
4. Keep each highlight line short, concrete, and scannable.
5. Make the post ambitious and clear, not dry and not hypey.
6. Use Telegram HTML formatting only (`<b>`, `<i>`, `<code>`,
   `<a href="...">`).
7. Keep the message easy to scan and under 1500 characters.
8. Save output to the exact file path requested by the caller.

## Output Rules

- Return only the final Telegram message body in the file.
- Do not add wrappers like "Here is your post", "Result:", or "Output:".
- Do not add YAML frontmatter, surrounding quotes, code fences, or unsupported
  HTML tags.
- Do not reduce the project description to a dry list of technical facts unless
  they support a stronger product point.
- Do not bury the release URL inside the CTA paragraph; it should appear near
  the top as a standalone line.
