# HardLink Deduplicator - Detect and Manage Duplicate Files with Hard Links

## Overview

HardLink Deduplicator is a tool for detecting and managing duplicate files on your system by utilizing hard links. It helps in reducing storage space usage by creating hard links between identical files.

## Features

- Detect and report duplicate files.
- Deduplicate files by creating hard links.
- Option to keep only one extra copy of the file or remove all duplicates.
- Minimum file size setting to consider for deduplication.

## Installation

### Prerequisites

- Supported operating systems: Windows
- Go version 1.14 or higher

### Installation Steps

1. Clone the repository:

```bash
git clone https://github.com/NIR3X/hardlinkdeduplicator
```

2. Change to the project directory:

```bash
cd hardlinkdeduplicator
```

3. Build the project:

```bash
go build -o hardlinkdeduplicator.exe .\cmd\hardlinkdeduplicator
```

4. Run the executable:

```bash
.\hardlinkdeduplicator -h
```

## Usage

```bash
.\hardlinkdeduplicator [options] path
```

## Options

* `-a`: Remove all duplicates (default is to keep one extra copy of the file).
* `-d`: Deduplicate files (not just report duplicates).
* `-s`: Minimum file size to consider for deduplication (in bytes).
* `-v`: Verbose output.

## Example

```bash
.\hardlinkdeduplicator -a -d -s 1024 -v C:\Path\To\Directory
```

## License
[![GNU AGPLv3 Image](https://www.gnu.org/graphics/agplv3-155x51.png)](https://www.gnu.org/licenses/agpl-3.0.html)  

This program is Free Software: You can use, study share and improve it at your
will. Specifically you can redistribute and/or modify it under the terms of the
[GNU Affero General Public License](https://www.gnu.org/licenses/agpl-3.0.html) as
published by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
