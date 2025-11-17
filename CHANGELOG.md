| Release | Date        | Comments                                                                                                        |
|---------|-------------|-----------------------------------------------------------------------------------------------------------------|
| 4.0.0   | 2025.11.13  | **Breaking changes**<br>Removed the customError package from this package, to avoid circular dependencies hell  |
| 4.0.0   | _continued_ | Glyph function calls are cleaner<br>All log init options are now wrapped in a single struct                     |                                                                                      
| 3.1.0   | 2025.11.10  | Swapped parameters in most glyph functions in terminalfx so that the text will stay in the usual white color    |
| 3.0.8   | 2025.11.10  | Added another option to prepend the messsage text with the username calling the tool                            |
| 3.0.6   | 2025.11.09  | Added an optional prefix option for each logfile entry                                                          |
| 3.0.5   | 2025.10.08  | Added a new "in progress" glyph                                                                                 |
| 3.0.4   | 2025.10.07  | Simplified glyph handling and printing                                                                          |
| 3.0.3   | 2025.10.06  | Removed mentions of customError/v2<br>GO version bump: 1.25.2                                                   |
| 3.0.2   | 2025.09.26  | **Breaking change**<br>terminal functions are now in their own subpackage                                       |
| 2.4.1   | 2025.09.16  | Added coloured (terminal) symbols in terminal.go                                                                |
| 2.3.1   | 2025.08.25  | Added proc information to Init()                                                                                |
| 2.2.4   | 2025.08.14  | GO version bump to 1.25.0<br>Updated CustomError dependency (version bump)                                      |
| 2.2.2   | 2025.08.13  | Added a ParseLevel() helper to convert a string loglevel to a LogLevelIota variable                             |
| 2.2.1   | 2025.08.13  | Simplified the logging subpackage, which now handles loglevels itself instead of relying on the calling tool    |
| 2.1.2   | 2025.08.11  | Reverted that warning loglevel; some brain freeze where I confounded loglevel and errorleveel                   |
| 2.1.1   | 2025.08.11  | Added a "Warning" loglevel; how could this have been forgotten ??                                               |
| 2.1.0   | 2025.08.09  | Hybrid approach to logging (WIP)                                                                                |
| 2.0.4   | 2025.08.09  | Updated builddeps<br>Added a String method to LogLevel                                                          |
| 2.0.3   | 2025.08.09  | Minor refactoring in the logging subpackage<br>GO version bump                                                  |
| 2.0.2   | 2025.08.05  | Version bump to ensure that we got rid of customError v1.x.y                                                    |
| 2.0.1   | 2025.07.27  | Added basic logging facilities                                                                                  |
| 1.12.2  | 2025.05.27  | Modified the paginating terminal size logic                                                                     |
| 1.12.1  | 2025.05.27  | Corrected error-handling in the Pager() functions                                                               |
| 1.12.0  | 2025.05.26  | Reworked a bit around the GetPassword() in debug mode                                                           |
| 1.11.1  | 2025.05.15  | Added a pager function, like the `more` command                                                                 |
| 1.10.1  | 2025.04.21  | Added a repository information extraction function                                                              |
| 1.9.0   | 2025.04.17  | Added a wrapper so that if we call GetPassword in DEBUG mode, we use                                            |
| 1.8.0   | 2025.02.18  | Added Center() and Right() for more text formatting features, GO version bump (1.23.6)                          |
| 1.7.0   | 2024.11.19  | Blue() actually printed in yellow; how comes I never noticed that ?!<br/>Updated to GO 1.23.3                   |
| 1.6.x   | 2024.xx.yy  |                                                                                                                 |
| 1.5.2   | 2024.06.22  | Added function to check mountpoint type (only supports linux)                                                   |
| 1.4.0   | 2024.04.29  | Added a generic prompt function                                                                                 |
| 1.3.2   | 2024.04.14  | Reverted v1.3.1                                                                                                 |
| 1.3.1   | 2024.04.14  | Reverted Encode/Decode strings to a default cipherkey                                                           |
| 1.3.0   | 2024.04.12  | SI() now supports more than uint64, is and also sign-aware                                                      |
| 1.2.1   | 2024.04.10  | Added Changelog()                                                                                               |
| 1.1.0   | 2024.04.10  | Added doc, CHANGELOG, extra functions in {encodeDecodePassword,terminal}.go                                     |
| 1.0.0   | 2024.04.10  | Initial version                                                                                                 |




