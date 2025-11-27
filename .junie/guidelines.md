# Junie Guidelines for `thujareader`

## Working with `docs/tasks.md`

- Use the existing phase structure (Phase 1, Phase 2, …) as the primary organization of work. Do not rename or remove phases; add new tasks within the most appropriate phase, or introduce a new phase only if absolutely necessary.
- Each task line must follow the pattern: `N. [ ] Description (Plan: P*, …; Reqs: R*, …)`.
- When a task is completed, change its marker from `[ ]` to `[x]` without altering its number, description, or links.
- If you need to add a task:
  - Place it in the correct phase, preserving numeric order (you may append at the end of the phase or renumber tasks in that phase consistently).
  - Ensure it references at least one plan item ID (P*) from `docs/plan.md` and one requirement ID (R*) from `docs/requirements.md`.
  - Use the same wording style and formatting as existing tasks.
- When changing or adding tasks that imply new functionality or behavior not covered by any current requirement or plan item:
  - First update `docs/requirements.md` with a new requirement (next sequential number and new R* ID).
  - Then update `docs/plan.md` with a new or revised plan item referencing that requirement.
  - Finally, add or update the task in `docs/tasks.md` so its `(Plan: …; Reqs: …)` section is accurate.
- Keep cross-references in sync: if a plan item ID or requirement ID changes, update all affected tasks.
- Maintain valid Markdown structure (headings, lists, checkboxes) and avoid introducing formatting styles that are inconsistent with the existing documents.
