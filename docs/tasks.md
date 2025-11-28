# Development Tasks

This checklist breaks the implementation plan into concrete technical tasks. Each task references its plan item (P*) from `docs/plan.md` and related requirements (R*) from `docs/requirements.md`.

Mark tasks as `[x]` when completed.

## Phase 1 – Project Setup and Architecture

1. [x] Initialize Go module layout with `cmd/thujareader` and `internal` packages (Plan: P1; Reqs: R1, R16)
2. [x] Add Bubble Tea and supporting TUI dependencies to `go.mod` (Plan: P1; Reqs: R1, R4, R16)
3. [x] Implement minimal `main` that starts a Bubble Tea program and exits cleanly (Plan: P1; Reqs: R1, R16)
4. [x] Define core domain structs for books, chapters/sections, and positions (Plan: P2; Reqs: R7, R8, R9, R10, R12, R18)
5. [x] Define bookmark and TOC entry types and interfaces (Plan: P2; Reqs: R9, R10, R12)

## Phase 2 – TUI Shell and DOS `edit.exe` Emulation

6. [x] Implement Bubble Tea model with top menu bar, main text area, and bottom status bar (Plan: P3; Reqs: R1, R4, R6, R16)
7. [x] Implement pseudo-graphics borders and box-drawing for `edit.exe`-style layout (Plan: P3; Reqs: R4, R6)
8. [x] Handle terminal resize messages and recompute layout and wrapping (Plan: P3; Reqs: R6, R16)
9. [x] Implement classic `edit.exe` color palette in a theme module (Plan: P4; Reqs: R4, R6)
10. [x] Implement color fallback paths for limited terminals (Plan: P4; Reqs: R4, R16)
11. [x] Expose minimal theme overrides via configuration hooks (Plan: P4; Reqs: R14, R16)
12. [x] Map `edit.exe`-style keybindings (Alt+menu, arrows, PgUp/PgDn, Home/End, F-keys) to internal commands (Plan: P5; Reqs: R1, R3, R5, R11, R12, R13, R17)
13. [x] Implement menu bar navigation and activation (Plan: P5; Reqs: R3, R5, R11, R12, R13, R17)
14. [x] Wire menu items and keybindings for Open, Exit, Help, Find, TOC, bookmarks, and recent files (Plan: P5; Reqs: R1, R3, R5, R11, R12, R13, R17)

## Phase 3 – Format Support and Unified Reader

15. [x] Implement EPUB loader that parses metadata, spine, and content documents (Plan: P6; Reqs: R2, R7, R9, R10, R11, R18)
16. [x] Normalize EPUB content to plain text paragraphs suitable for terminal display (Plan: P6; Reqs: R7, R9, R11, R18)
17. [x] Implement EPUB TOC extraction and mapping to internal positions (Plan: P6; Reqs: R7, R9, R10)
18. [x] Implement FB2 XML parser for sections, titles, and paragraphs (Plan: P7; Reqs: R2, R8, R9, R10, R11, R18)
19. [x] Normalize FB2 structure into paragraphs and section metadata (Plan: P7; Reqs: R8, R9, R10)
20. [x] Implement unified reader interface and adapters for EPUB and FB2 (Plan: P8; Reqs: R2, R7, R8, R9, R10, R11, R18)
21. [x] Implement command-line argument handling to open a file on startup via the unified reader (Plan: P8; Reqs: R2, R9)
22. [x] Implement file open flow from the File → Open menu item using the unified reader (Plan: P5, P8; Reqs: R2, R3, R9)

## Phase 4 – Core Reading Experience

23. [x] Implement text wrapping based on current terminal width and height (Plan: P9; Reqs: R4, R6, R7, R8, R9, R18, R19)
24. [x] Implement scrolling behavior (line, page up/down) matching `edit.exe` semantics (Plan: P9; Reqs: R5, R9, R18)
25. [x] Integrate navigation model so scroll commands update and consume positions (Plan: P9; Reqs: R9)
26. [x] Implement TOC dialog and allow jumping to selected entries (Plan: P10; Reqs: R9, R10)
27. [x] Display current location (chapter, percentage) in the status bar (Plan: P10; Reqs: R9, R10)
28. [x] Implement Find dialog styled like `edit.exe` (Plan: P11; Reqs: R11)
29. [x] Implement text search over normalized content with find-next behavior (Plan: P11; Reqs: R11, R18)
30. [x] Show user feedback when no more matches are found (Plan: P11; Reqs: R11)

## Phase 5 – Persistence and Advanced Features

31. [ ] Implement in-memory bookmark management commands (add, list, delete, jump) (Plan: P12; Reqs: R9, R12)
32. [ ] Implement bookmarks persistence on disk (Plan: P12; Reqs: R9, R12, R13)
33. [ ] Implement recent files tracking in memory (Plan: P12; Reqs: R13)
34. [ ] Persist recent files between sessions (Plan: P12; Reqs: R13, R14)
35. [ ] Design on-disk store (e.g., JSON files) for positions, bookmarks, and recent files (Plan: P13; Reqs: R9, R12, R13, R14)
36. [ ] Implement loading of persisted state at startup and saving on shutdown (Plan: P13; Reqs: R9, R12, R13, R14)
37. [ ] Handle corrupted state files gracefully with fallbacks (Plan: P13; Reqs: R15)
38. [ ] Implement configuration file format and default search paths per OS (Plan: P14; Reqs: R14, R16)
39. [ ] Load configuration with defaults and validation on startup (Plan: P14; Reqs: R14, R16)
40. [ ] Apply configuration options (theme overrides, recent list size, default library path) (Plan: P14; Reqs: R14)

## Phase 6 – Robustness, Cross-Platform Behavior, and UX Polish

41. [ ] Implement reusable error dialog component styled like `edit.exe` (Plan: P15; Reqs: R2, R7, R8, R11, R15)
42. [ ] Route parsing and I/O errors through error dialogs instead of panics (Plan: P15; Reqs: R2, R7, R8, R11, R15)
43. [ ] Ensure terminal state is restored on all error exit paths (Plan: P15; Reqs: R1, R15, R16)
44. [ ] Detect terminal color and keycode capabilities where feasible (Plan: P16; Reqs: R4, R5, R6, R16)
45. [ ] Implement behavior fallbacks for limited terminals (Plan: P16; Reqs: R4, R5, R6, R16)
46. [ ] Test behavior on Windows, Linux, and macOS terminals and adjust mappings (Plan: P16; Reqs: R16)
47. [ ] Implement in-app Help screen listing keybindings and core concepts (Plan: P17; Reqs: R5, R17)
48. [ ] Wire Help into menu bar and keybindings mirroring `edit.exe` (Plan: P17; Reqs: R5, R17)
49. [ ] Write user-facing documentation describing terminal setup to emulate DOS `edit.exe` (Plan: P18; Reqs: R4, R6, R14, R16)

## Phase 7 – Testing and Quality Assurance

50. [ ] Add unit tests for EPUB and FB2 parsing and normalization (Plan: P6, P7; Reqs: R7, R8, R18)
51. [ ] Add tests for navigation, search, and bookmarks logic at the domain level (Plan: P2, P9, P10, P11, P12; Reqs: R9, R10, R11, R12, R18)
52. [ ] Add integration tests or scripted runs to verify startup, open, navigation, and shutdown flows (Plan: P1, P3, P5, P9; Reqs: R1, R2, R3, R4, R5, R16)
53. [ ] Manually verify appearance and behavior against a reference DOS `edit.exe` setup (Plan: P3, P4, P5, P18; Reqs: R4, R5, R6)

## Phase 8 – Advanced Format Support

54. [ ] Extend EPUB parsing to support complex structures and edge cases (Plan: P19; Reqs: R2, R7, R9, R10, R11, R18)
55. [ ] Extend FB2 parsing to support complex structures and edge cases (Plan: P20; Reqs: R2, R8, R9, R10, R11, R18)
56. [ ] Implement optional word hyphenation for long words in wrapped text (Plan: P9; Reqs: R20)
57. [ ] Implement optional text justification mode with safe fallbacks (Plan: P9, P18; Reqs: R6, R19, R21)
