# Requirements Document

## Introduction

`thujareader` is a terminal-based EPUB and FB2 reader implemented in Go using the Bubble Tea TUI framework. The application emulates the classic DOS `edit.exe` editor used in Norton Commander, replicating its pseudo-graphics user interface, color scheme, keyboard controls, and overall interaction model while providing modern e-book reading capabilities. The focus is on a fast, keyboard-centric reading experience that works consistently across terminals and platforms.

## Functional Requirements

1. **R1 – Application Startup and Shutdown**
   - **User Story**: As a user, I want to start and exit the reader cleanly so that I can reliably use it in my terminal sessions.
   - **Acceptance Criteria**:
     - WHEN the user runs the `thujareader` command with no arguments THEN the system SHALL start in the main screen with an empty editor area and a DOS `edit.exe`-style menu/status layout.
     - WHEN the user presses the standard `edit.exe` quit key combination (e.g., `Alt+F` then `X`) or invokes Exit from the menu THEN the system SHALL close the TUI, restore the terminal state, and exit with code 0.

2. **R2 – Open EPUB and FB2 Files from CLI Arguments**
   - **User Story**: As a user, I want to open EPUB and FB2 files directly from the command line so that I can start reading immediately.
   - **Acceptance Criteria**:
     - WHEN the user runs `thujareader <path-to-file>` and the file is a valid EPUB or FB2 THEN the system SHALL open it directly into the reader and position the caret at the beginning of the book.
     - WHEN the user provides a path to a non-existent file THEN the system SHALL show an `edit.exe`-style error dialog and exit with a non-zero code after user confirmation.

3. **R3 – File Open Dialog (Internal File Picker)**
   - **User Story**: As a user, I want to open files through an in-app file dialog so that I can browse for books without remembering paths.
   - **Acceptance Criteria**:
     - WHEN the user triggers the `Open` action using the same key sequence as DOS `edit.exe` (e.g., `Alt+F` then `O`) THEN the system SHALL show a pseudo-graphics file open dialog listing files and directories.
     - WHEN the user selects an EPUB or FB2 file in the dialog and confirms THEN the system SHALL load and display the book.
     - WHEN the user attempts to open a file with an unsupported extension THEN the system SHALL display an error dialog and stay in the file picker.

4. **R4 – DOS `edit.exe`-Style Layout and Color Scheme**
   - **User Story**: As a user, I want the UI to look like the classic DOS `edit.exe` so that the interface feels familiar and nostalgic.
   - **Acceptance Criteria**:
     - WHEN the application starts THEN the system SHALL render a top menu bar, main text area, and bottom status bar using box-drawing pseudo-graphics identical in layout to DOS `edit.exe`.
     - WHEN the terminal supports ANSI colors THEN the system SHALL apply a color palette matching the default `edit.exe` scheme (blue background, cyan menu bar, etc.).
     - WHEN the terminal does not support colors or runs in a limited palette THEN the system SHALL fall back to a safe, high-contrast approximation while preserving the layout.

5. **R5 – Keyboard Controls Matching DOS `edit.exe`**
   - **User Story**: As a user familiar with `edit.exe`, I want the same keyboard shortcuts so that I can navigate and control the reader without learning new bindings.
   - **Acceptance Criteria**:
     - WHEN the user presses arrow keys, Page Up/Down, Home, End THEN the system SHALL scroll the text view exactly as `edit.exe` would scroll a text file.
     - WHEN the user invokes menu actions via `Alt+<letter>` and navigates menus with arrow keys and Enter THEN the system SHALL open and execute actions consistent with `edit.exe` (Open, Exit, Help, etc.).
     - WHEN the user presses keys not bound to any command THEN the system SHALL ignore them and keep the UI stable.

6. **R6 – Text Rendering and Fonts**
   - **User Story**: As a user, I want the reader to visually resemble `edit.exe`, including font feel, so that it evokes the DOS reading/editing experience.
   - **Acceptance Criteria**:
     - WHEN the application starts THEN the system SHALL assume a fixed-width terminal font and render text using a layout compatible with DOS codepage-style fonts.
     - WHEN documenting usage THEN the system SHALL describe how to configure the terminal (font and size) to best match the DOS `edit.exe` appearance.
     - WHEN the terminal window is resized THEN the system SHALL reflow the visible text region while preserving monospaced alignment and pseudo-graphics.

7. **R7 – EPUB Parsing and Normalization**
   - **User Story**: As a user, I want EPUB books to render cleanly as plain text so that I can read them comfortably in a terminal.
   - **Acceptance Criteria**:
     - WHEN the user opens a valid EPUB file THEN the system SHALL parse its content, extract the reading order, strip unsupported formatting, and present readable text wrapped to the terminal width.
     - WHEN the EPUB contains images, complex layouts, or unsupported media THEN the system SHALL gracefully omit them, optionally replacing them with textual placeholders (e.g., `[Image omitted]`).

8. **R8 – FB2 Parsing and Normalization**
   - **User Story**: As a user, I want FB2 books to display correctly so that I can read my existing library without conversion.
   - **Acceptance Criteria**:
     - WHEN the user opens a valid FB2 file THEN the system SHALL parse XML, extract text sections, titles, and paragraphs, and render them as wrapped text.
     - WHEN the FB2 is malformed THEN the system SHALL show an error dialog indicating parsing failure and refuse to open the file.

9. **R9 – Document Navigation and Position Persistence**
   - **User Story**: As a user, I want to navigate within a book and resume where I left off so that I can read in multiple sessions.
   - **Acceptance Criteria**:
     - WHEN the user scrolls through a book using the keyboard THEN the system SHALL update the current position without noticeable lag.
     - WHEN the user exits the application after reading part of a book THEN the system SHALL persist the current position.
     - WHEN the user reopens the same book later THEN the system SHALL restore the reading position to the last saved location.

10. **R10 – Table of Contents and Section Jumping**
    - **User Story**: As a user, I want to jump quickly to chapters or sections so that I can navigate long books efficiently.
    - **Acceptance Criteria**:
      - WHEN a book has a table of contents THEN the system SHALL expose it via a menu or keybinding mirroring `edit.exe`-style dialogs.
      - WHEN the user selects a chapter/section from the TOC THEN the system SHALL jump the view to the corresponding location in the book.

11. **R11 – Search Within Document**
    - **User Story**: As a user, I want to search for text within the current book so that I can locate specific passages.
    - **Acceptance Criteria**:
      - WHEN the user invokes the Find functionality using the standard `edit.exe` keybinding THEN the system SHALL open a DOS-style search dialog.
      - WHEN the user enters a search term and confirms THEN the system SHALL highlight and scroll to the next occurrence, if any.
      - WHEN there are no more matches THEN the system SHALL display a small status message indicating no further matches.

12. **R12 – Bookmarks and Quick Jump**
    - **User Story**: As a user, I want to set bookmarks so that I can quickly return to important spots in a book.
    - **Acceptance Criteria**:
      - WHEN the user triggers the bookmark creation command THEN the system SHALL store a named bookmark at the current position.
      - WHEN the user opens the bookmark list THEN the system SHALL show saved bookmarks in a pseudo-graphics dialog and allow jumping to them.
      - WHEN the user deletes a bookmark THEN the system SHALL remove it from persistent storage.

13. **R13 – Recent Files List**
    - **User Story**: As a user, I want a list of recently opened books so that I can reopen them quickly.
    - **Acceptance Criteria**:
      - WHEN the user opens books over time THEN the system SHALL maintain a recent files list stored between sessions.
      - WHEN the user opens the recent files menu item THEN the system SHALL show the list in an `edit.exe`-style menu or dialog and allow selecting a book to open.

14. **R14 – Configuration and Preferences**
    - **User Story**: As a user, I want to configure basic behavior (paths, recent list size, visual tweaks) so that the reader fits my workflow.
    - **Acceptance Criteria**:
      - WHEN the user edits the configuration file or uses configuration options (if exposed via UI) THEN the system SHALL apply settings such as default library path, recent list size, and optional color overrides.
      - WHEN the configuration file is missing or invalid THEN the system SHALL fall back to sensible defaults and, if invalid, report the problem in a non-blocking status message.

15. **R15 – Error Handling and Robustness**
    - **User Story**: As a user, I want clear error messages and stable behavior so that problems do not crash my terminal or lose my place.
    - **Acceptance Criteria**:
      - WHEN the system encounters I/O errors, parsing errors, or unsupported files THEN it SHALL show a clear, concise error dialog styled like `edit.exe` and avoid panics.
      - WHEN a fatal error prevents continuing THEN the system SHALL restore the terminal state before exiting.

16. **R16 – Cross-Platform Terminal Support**
    - **User Story**: As a user, I want the reader to work on common platforms so that I can use it on my preferred OS.
    - **Acceptance Criteria**:
      - WHEN run on modern Windows, Linux, or macOS terminals THEN the system SHALL render the TUI correctly using Bubble Tea without corrupting the console.
      - WHEN terminal capabilities differ (e.g., function keys, color support) THEN the system SHALL adapt bindings or fall back gracefully while keeping core navigation and reading functional.

17. **R17 – Help and Key Reference**
    - **User Story**: As a new user, I want an in-app help screen so that I can learn the DOS-style controls.
    - **Acceptance Criteria**:
      - WHEN the user invokes Help (using the same key sequence/menu as `edit.exe`) THEN the system SHALL display a help screen listing key commands, navigation, and basic usage.
      - WHEN the user closes the help screen THEN the system SHALL return to the previous view without losing reading position.

18. **R18 – Performance on Large Books**
    - **User Story**: As a user, I want smooth scrolling even in large books so that reading is comfortable.
    - **Acceptance Criteria**:
      - WHEN the user opens a large EPUB or FB2 (e.g., several MBs, thousands of pages) THEN the system SHALL remain responsive and scroll without noticeable freezes.
      - WHEN memory usage would become excessive for extremely large books THEN the system SHALL use streaming or chunked loading strategies to remain stable.
