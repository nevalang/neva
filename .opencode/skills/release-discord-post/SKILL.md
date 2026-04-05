# Release Discord Post

Create a Discord-ready release announcement in English for the target Neva Discord news-style channel provided by the caller.

## Tools
- `Read` for release input context from the prompt.
- `Write` to save the final message file.

## Workflow
1. Read the release context from the prompt: repository, target channel name, tag, release name, release URL, and release notes.
2. Read `.opencode/shared/release-marketing.md` for positioning, CTA, and tone guidance.
3. Produce a concise Discord post in English using this exact high-level structure:
   - Line 1: release heading like `Neva vX.Y.Z` or `Neva vX.Y.Z — Release Title`
   - Optional short punchy framing line when the release has a clear theme
   - Short intro paragraph, 1-2 sentences max
   - `Highlights` or `What shipped` label
   - 4-8 short bullet lines with the most relevant release changes
   - `Full notes:` label or short CTA line
   - Final line: release URL by itself
4. Keep the intro and bullets grounded in the actual release notes. If the notes are sparse, stay conservative and do not invent features or guarantees.
5. Optimize for Discord readability:
   - plain text or light Markdown only
   - short paragraphs
   - scannable bullets
   - no Telegram-style HTML
6. Keep the tone clear, sharp, mildly energetic, and community-facing. Avoid corporate phrasing and avoid overselling.
7. Favor the underlying communication style from prior Neva Discord release posts:
   - strong, direct opener
   - practical framing around what changed and why it matters
   - compact highlight list
   - final pointer to full release notes
   Do not imitate exact wording or reproduce specific old content unless it is also present in the current release notes.
8. Use `@everyone` only when the caller explicitly asks for it in the prompt. Otherwise do not mention everyone or here-style pings.
9. Use `.opencode/shared/release-marketing.md` to frame Neva well, but only when it helps explain real release value.
10. Keep the message under 1800 characters.
11. Save output to the exact file path requested by the caller.

## Output Rules
- Return only the final Discord message body in the file.
- The file content is consumed as raw text in a webhook JSON `content` field, so it must be directly sendable without edits.
- Do not add wrappers like "Here is your post", "Result:", "Output:", or follow-up suggestions.
- Do not add YAML frontmatter or surrounding quotes.
- Do not include code fences.
- Do not include markdown headings.
- Do not turn the post into a long prose block.
- Do not bury the release URL inside another sentence; keep it on its own final line.
- Prefer short, high-signal bullets over exhaustive coverage of every patch note.
- Do not force a generic CTA if a simple `Full notes:` handoff reads better for this release.
