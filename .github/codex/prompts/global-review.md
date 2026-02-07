# Global Review Guidance

Scope:
- Review the PR diff between base and head.
- Identify correctness, security, reliability, or breaking-change risks.
- Call out missing tests or docs updates.
- Suggest concrete fixes with file paths and line numbers when possible.

Focus checks:
- Configuration, secrets, and CI workflow changes.
- Backward compatibility and migration steps.
- Error handling and logging consistency.

Output format (strict):
- Summary (1-3 bullets)
- Findings (bulleted list, include severity: High, Medium, Low)
- Tests (list suggested commands; use "Not run" if none)
- Inline Comments (JSON) in a ```json code block

Inline Comments JSON schema:
- An array of objects: { "path": "relative/file", "line": 123, "side": "RIGHT", "body": "comment", "severity": "High|Medium|Low" }
- "line" must be the line number in the PR head (RIGHT side).
- If there are no inline comments, output an empty array: [].

If no issues are found, say "No blocking issues found" and still include Tests and the JSON block.
