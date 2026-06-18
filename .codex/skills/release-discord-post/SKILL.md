---
name: "release-discord-post"
description: "Create a Discord-ready Neva release announcement from a GitHub release payload. Use for official Neva Discord release posts."
---

# Release Discord Post

Create a Discord-ready release announcement in English for the target Neva Discord news-style channel provided by the caller.

## Workflow

1. Read the release context from the prompt: repository, target channel name,
   tag, release name, release URL, and release notes.
2. Read `.codex/shared/release-marketing.md` for positioning, CTA, and tone
   guidance.
3. Produce a concise Discord post in English using this high-level structure:
   - Line 1: release heading like `Neva vX.Y.Z` or
     `Neva vX.Y.Z - Release Title`
   - Optional short punchy framing line when the release has a clear theme
   - Short intro paragraph, 1-2 sentences max
   - `Highlights` or `What shipped` label
   - 4-8 short bullet lines with the most relevant release changes
   - `Full notes:` label or short CTA line
   - Final line: release URL by itself
4. Keep the intro and bullets grounded in the actual release notes. If the notes
   are sparse, stay conservative and do not invent features or guarantees.
5. Optimize for Discord readability: plain text or light Markdown, short
   paragraphs, scannable bullets, no Telegram-style HTML.
6. Keep the tone clear, sharp, mildly energetic, and community-facing. Avoid
   corporate phrasing and overselling.
7. Use `@everyone` only when the caller explicitly asks for it.
8. Keep the message under 1800 characters.
9. Save output to the exact file path requested by the caller.

## Output Rules

- Return only the final Discord message body in the file.
- Do not add wrappers like "Here is your post", "Result:", or "Output:".
- Do not add YAML frontmatter, surrounding quotes, or code fences.
- Do not bury the release URL inside another sentence; keep it on its own final
  line.
- Prefer short, high-signal bullets over exhaustive coverage of every patch
  note.
