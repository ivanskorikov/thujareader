# Development Tasks

This checklist breaks the implementation plan into concrete technical tasks. Each task references its plan item (P*) from `docs/plan.md` and related requirements (R*) from `docs/requirements.md`.

Mark tasks as `[x]` when completed.

## Phase 1 – Project Setup and Architecture

1. [ ] Initialize Go module layout with `cmd/thujareader` and `internal` packages (Plan: P1; Reqs: R1, R16)
2. [ ] Add Bubble Tea and supporting TUI dependencies to `go.mod` (Plan: P1; Reqs: R1, R4, R16)
3. [ ] Implement minimal `main` that starts a Bubble Tea program and exits cleanly (Plan: P1; Reqs: R1, R16)
4. [ ] Define core domain structs for books, chapters/sections, and positions (Plan: P2; Reqs: R7, R8, R9, R10, R12, R18)
5. [ ] Define bookmark and TOC entry types and interfaces (Plan: P2; Reqs: R9, R10, R12)

## Phase 2 – TUI Shell and DOS `edit.exe` Emulation

6. [ ] Implement Bubble Tea model with top menu bar, main text area, and bottom status bar (Plan: P3; Reqs: R1, R4, R6, R16)
7. [ ] Implement pseudo-graphics borders and box-drawing for `edit.exe`-style layout (Plan: P3; Reqs: R4, R6)
8. [ ] Handle terminal resize messages and recompute layout and wrapping (Plan: P3; Reqs: R6, R16)
9. [ ] Implement classic `edit.exe` color palette in a theme module (Plan: P4; Reqs: R4, R6)
10. [ ] Implement color fallback paths for limited terminals (Plan: P4; Reqs: R4, R16)
11. [ ] Expose minimal theme overrides via configuration hooks (Plan: P4; Reqs: R14, R16)
12. [ ] Map `edit.exe`-style keybindings (Alt+menu, arrows, PgUp/PgDn, Home/End, F-keys) to internal commands (Plan: P5; Reqs: R1, R3, R5, R11, R12, R13, R17)
13. [ ] Implement menu bar navigation and activation (Plan: P5; Reqs: R3, R5, R11, R12, R13, R17)
14. [ ] Wire menu items and keybindings for Open, Exit, Help, Find, TOC, bookmarks, and recent files (Plan: P5; Reqs: R1, R3, R5, R11, R12, R13, R17)

## Phase 3 – Format Support and Unified Reader

15. [ ] Implement EPUB loader that parses metadata, spine, and content documents (Plan: P6; Reqs: R2, R7, R9, R10, R11, R18)
16. [ ] Normalize EPUB content to plain text paragraphs suitable for terminal display (Plan: P6; Reqs: R7, R9, R11, R18)
17. [ ] Implement EPUB TOC extraction and mapping to internal positions (Plan: P6; Reqs: R7, R9, R10)
18. [ ] Implement FB2 XML parser for sections, titles, and paragraphs (Plan: P7; Reqs: R2, R8, R9, R10, R11, R18)
19. [ ] Normalize FB2 structure into paragraphs and section metadata (Plan: P7; Reqs: R8, R9, R10)
20. [ ] Implement unified reader interface and adapters for EPUB and FB2 (Plan: P8; Reqs: R2, R7, R8, R9, R10, R11, R18)
21. [ ] Implement command-line argument handling to open a file on startup via the unified reader (Plan: P8; Reqs: R2, R9)

## Phase 4 – Core Reading Experience

22. [ ] Implement text wrapping based on current terminal width and height (Plan: P9; Reqs: R4, R6, R7, R8, R9, R18)
23. [ ] Implement scrolling behavior (line, page up/down) matching `edit.exe` semantics (Plan: P9; Reqs: R5, R9, R18)
24. [ ] Integrate navigation model so scroll commands update and consume positions (Plan: P9; Reqs: R9)
25. [ ] Implement TOC dialog and allow jumping to selected entries (Plan: P10; Reqs: R9, R10)
26. [ ] Display current location (chapter, percentage) in the status bar (Plan: P10; Reqs: R9, R10)
27. [ ] Implement Find dialog styled like `edit.exe` (Plan: P11; Reqs: R11)
28. [ ] Implement text search over normalized content with find-next behavior (Plan: P11; Reqs: R11, R18)
29. [ ] Show user feedback when no more matches are found (Plan: P11; Reqs: R11)

## Phase 5 – Persistence and Advanced Features

30. [ ] Implement in-memory bookmark management commands (add, list, delete, jump) (Plan: P12; Reqs: R9, R12)
31. [ ] Implement bookmarks persistence on disk (Plan: P12; Reqs: R9, R12, R13)
32. [ ] Implement recent files tracking in memory (Plan: P12; Reqs: R13)
33. [ ] Persist recent files between sessions (Plan: P12; Reqs: R13, R14)
34. [ ] Design on-disk store (e.g., JSON files) for positions, bookmarks, and recent files (Plan: P13; Reqs: R9, R12, R13, R14)
35. [ ] Implement loading of persisted state at startup and saving on shutdown (Plan: P13; Reqs: R9, R12, R13, R14)
36. [ ] Handle corrupted state files gracefully with fallbacks (Plan: P13; Reqs: R15)
37. [ ] Implement configuration file format and default search paths per OS (Plan: P14; Reqs: R14, R16)
38. [ ] Load configuration with defaults and validation on startup (Plan: P14; Reqs: R14, R16)
39. [ ] Apply configuration options (theme overrides, recent list size, default library path) (Plan: P14; Reqs: R14)

## Phase 6 – Robustness, Cross-Platform Behavior, and UX Polish

40. [ ] Implement reusable error dialog component styled like `edit.exe` (Plan: P15; Reqs: R2, R7, R8, R11, R15)
41. [ ] Route parsing and I/O errors through error dialogs instead of panics (Plan: P15; Reqs: R2, R7, R8, R11, R15)
42. [ ] Ensure terminal state is restored on all error exit paths (Plan: P15; Reqs: R1, R15, R16)
43. [ ] Detect terminal color and keycode capabilities where feasible (Plan: P16; Reqs: R4, R5, R6, R16)
44. [ ] Implement behavior fallbacks for limited terminals (Plan: P16; Reqs: R4, R5, R6, R16)
45. [ ] Test behavior on Windows, Linux, and macOS terminals and adjust mappings (Plan: P16; Reqs: R16)
46. [ ] Implement in-app Help screen listing keybindings and core concepts (Plan: P17; Reqs: R5, R17)
47. [ ] Wire Help into menu bar and keybindings mirroring `edit.exe` (Plan: P17; Reqs: R5, R17)
48. [ ] Write user-facing documentation describing terminal setup to emulate DOS `edit.exe` (Plan: P18; Reqs: R4, R6, R14, R16)

## Phase 7 – Testing and Quality Assurance

49. [ ] Add unit tests for EPUB and FB2 parsing and normalization (Plan: P6, P7; Reqs: R7, R8, R18)
50. [ ] Add tests for navigation, search, and bookmarks logic at the domain level (Plan: P2, P9, P10, P11, P12; Reqs: R9, R10, R11, R12, R18)
51. [ ] Add integration tests or scripted runs to verify startup, open, navigation, and shutdown flows (Plan: P1, P3, P5, P9; Reqs: R1, R2, R3, R4, R5, R16)
52. [ ] Manually verify appearance and behavior against a reference DOS `edit.exe` setup (Plan: P3, P4, P5, P18; Reqs: R4, R5, R6)
