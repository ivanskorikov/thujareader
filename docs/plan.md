# Implementation Plan

This implementation plan describes how to realize the `thujareader` requirements. Each plan item references the related requirements by ID (R1, R2, …) from `docs/requirements.md` and is assigned a priority.

## 1. Project Foundation and Architecture

**P1. Initialize Go Project Structure and Dependencies**  
Priority: High  
Related Requirements: R1, R16

- Define the main module layout (`cmd/thujareader`, `internal/ui`, `internal/reader`, `internal/config`, etc.).
- Add and configure Bubble Tea and supporting TUI libraries (lipgloss, bubbles, etc.).
- Set up a minimal main loop that starts and stops cleanly in common terminals.

**P2. Core Domain Models for Books and Positions**  
Priority: High  
Related Requirements: R7, R8, R9, R10, R12, R18

- Design data structures representing a logical book (metadata, chapters/sections, linearized text stream).
- Define a navigation/position model (chapter index, offset, percentage) independent of UI.
- Define bookmark and TOC entry models.

## 2. TUI Shell Emulating DOS `edit.exe`

**P3. Base TUI Layout (Menu Bar, Text Area, Status Bar)**  
Priority: High  
Related Requirements: R1, R4, R6, R16

- Implement a Bubble Tea model that renders a full-screen layout matching `edit.exe` (top menu bar, main area, bottom status bar).
- Implement box-drawing and pseudo-graphics borders consistent with DOS.
- Handle terminal resize events to recompute layout and text wrapping.

**P4. Color Scheme and Theme Configuration**  
Priority: High  
Related Requirements: R4, R6, R14, R16

- Define a color palette approximating the classic `edit.exe` scheme.
- Implement theme abstraction with fallbacks for limited color support.
- Expose minimal theme overrides via configuration.

**P5. Keyboard Input and Menu Interaction Matching `edit.exe`**  
Priority: High  
Related Requirements: R1, R3, R5, R11, R12, R13, R17

- Map keybindings (Alt+menu, arrows, PgUp/PgDn, Home/End, F-keys) to internal commands.
- Implement menu navigation and activation using pseudo-graphics menus.
- Wire keybindings for Open, Exit, Find, Help, bookmarks, recent files, and TOC.

## 3. Format Support and Parsing

**P6. EPUB Parsing and Normalization Layer**  
Priority: High  
Related Requirements: R2, R7, R9, R10, R11, R18

- Integrate or implement an EPUB parsing component to read metadata, spine, and content documents.
- Normalize HTML/XHTML content into a sequence of plain text paragraphs suitable for terminal display.
- Provide TOC extraction and mapping from TOC entries to internal positions.

**P7. FB2 Parsing and Normalization Layer**  
Priority: High  
Related Requirements: R2, R8, R9, R10, R11, R18

- Implement FB2 XML parsing for metadata and text sections.
- Convert FB2 structural tags (sections, titles, paragraphs) into normalized paragraphs.
- Provide TOC-like navigation based on FB2 sections.

**P8. Unified Reader Abstraction**  
Priority: High  
Related Requirements: R2, R7, R8, R9, R10, R11, R18

- Define an interface for book loaders (EPUB, FB2) exposing a unified view of text, TOC, and metadata.
- Implement adapters for EPUB and FB2 that satisfy this interface.
- Implement error reporting channels for parsing and I/O problems.

## 4. Core Reading Experience

**P9. Text Rendering and Scrolling**  
Priority: High  
Related Requirements: R4, R5, R6, R7, R8, R9, R18

- Implement line wrapping and pagination respecting terminal width and height.
- Implement smooth scrolling behavior matching `edit.exe` semantics.
- Integrate the navigation model so that scroll operations update and consume positions.

**P10. Navigation, TOC, and Section Jumping**  
Priority: Medium  
Related Requirements: R9, R10, R11, R12, R17

- Implement a TOC dialog accessible via menu/keybinding.
- Support jumping to TOC entries and syncing displayed position.
- Maintain and display current location (percent, chapter) in the status bar.

**P11. Search Within Document**  
Priority: Medium  
Related Requirements: R11, R18

- Implement a Find dialog matching `edit.exe` style.
- Implement incremental or buffered text search over the normalized text model.
- Support find-next behavior and wrap/no-wrap messaging.

**P12. Bookmarks and Recent Files**  
Priority: Medium  
Related Requirements: R9, R12, R13, R14

- Implement in-memory bookmark management and persistence.
- Implement a recent files list with a dedicated dialog/menu.
- Integrate both features into the TUI menus and keybindings.

## 5. Persistence, Configuration, and State Management

**P13. Persistence of Reading State and Metadata**  
Priority: Medium  
Related Requirements: R9, R12, R13, R14

- Design a small on-disk store (e.g., JSON or Bolt-like DB) for bookmarks, last positions, and recent files.
- Implement loading and saving of this state at startup/shutdown.
- Handle corruption and versioning safely with fallbacks.

**P14. Configuration System**  
Priority: Low  
Related Requirements: R6, R14, R16

- Define a configuration file format and default location per OS.
- Implement loading with default fallbacks and validation.
- Expose relevant settings such as theme overrides, paths, and recent list size.

## 6. Error Handling, Robustness, and Cross-Platform Concerns

**P15. Error Dialogs and Non-Fatal Failures**  
Priority: High  
Related Requirements: R2, R7, R8, R11, R15

- Implement a reusable error dialog component styled like `edit.exe` message boxes.
- Ensure all parsing and I/O errors route through this component instead of panicking.

**P16. Terminal Capability Detection and Adaptation**  
Priority: Medium  
Related Requirements: R4, R5, R6, R16

- Detect color and keycode capabilities where possible.
- Provide fallbacks for limited terminals while keeping core navigation functional.

## 7. Help, Documentation, and UX Polish

**P17. In-App Help and Key Reference**  
Priority: Medium  
Related Requirements: R5, R17

- Implement a help screen/dialog listing keybindings and basic workflows.
- Integrate Help into the menu bar in the same position as `edit.exe` where feasible.

**P18. User Documentation and Terminal Setup Guide**  
Priority: Low  
Related Requirements: R4, R6, R14, R16

- Document how to configure terminal fonts, colors, and sizes to best emulate DOS `edit.exe`.
- Provide examples and notes on differences between platforms.
