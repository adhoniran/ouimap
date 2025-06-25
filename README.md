# OUImap 

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=adhoniran_ouimap&metric=alert_status)](https://sonarcloud.io/dashboard?id=adhoniran_ouimap)
[![Donate with PayPal](https://img.shields.io/badge/Donate%20with-PayPal-blue?logo=paypal&logoColor=white)](https://www.paypal.com/donate/?business=N6AH25Q2D4BL8&no_recurring=0&item_name=Contribute+to+the+future+of+our+projects.+Your+donation+via+PayPal+empowers+us+to+keep+creating+and+growing...+Thank+you%21&currency_code=USD)
[![Download](https://img.shields.io/badge/Download-Latest-red?logo=github)](https://github.com/adhoniran/ouimap/releases/latest)

OUImap is a command-line (CLI) tool for querying Organizationally Unique Identifier (OUI) information and MAC address prefixes from an always up-to-date Wireshark database. It allows you to quickly look up which vendor is associated with a given MAC address (or part of it) and also supports text-based searches for vendor names.

## Table of Contents

- [Overview](#overview)
- [Main Features](#main-features)
- [Build](#build)
- [Usage](#usage)
- [Execution Example](#execution-example)
- [Automatic Updates](#automatic-updates)
- [Contributing](#contributing)
- [License](#license)
- [Donations](#donations)

## Overview

OUImap is designed for developers, network analysts and enthusiasts who need to quickly identify the vendor of a complete or partial MAC address. It downloads the vendor database from Wireshark, keeps a local copy for reference, and provides an interactive search mode.

## Main Features

- Cross-platform support: Linux, macOS and Windows
- Multi-architecture: AMD64 and ARM64 support
- Quick mode: Pass search terms directly as command-line arguments for quick lookups
- Interactive mode: Multi-line input interface for complex searches
- Weekly automatic download of the Wireshark vendor database
- Integrity check of the downloaded database file
- Local storage and automatic database update
- MAC prefix (OUI) lookups in standardized formats (XX:XX:XX)
- Support for various input formats (e.g., 00-50-56, 00:50:56, 001a.b623.3499, etc.)
- Text-based search for vendor names (e.g., “Intel,” “Dell”)
- Detailed results, including the search time and the total number of records found
- Simple and interactive command-line tool
- Automatic version checking: Notifies when new versions are available


## Build

### Requirements

If you want to build OUImap from source, you need:

- Git 2.47.1 (or higher).
- Go compiler (version 1.23.4 or later), available at https://go.dev.
- Internet connection to download the vendor database if no local copy exists or if your local copy is out of date.
- Compatible operating system (Linux, macOS, Windows or any Go supported).

### Step 1: Clone the Repository

If you want to compile from source:

1. Clone the repository:
   ```bash
   git clone https://github.com/adhoniran/ouimap.git ouimap
   ```

2. Go to the project folder:  
   ```bash
   cd ouimap
   ```

### Step 2: Compile

In the terminal, run:

For Linux or macOS:
```bash
sh build.run
```

For Windows 11 (Powershell):
```powershell
Invoke-Expression (Get-Content -Raw .\build.run)
```

This will produce the “ouimap” executable (or “ouimap.exe” on Windows) in the same project folder.

### Step 3: Run

On Linux or macOS:
```
./ouimap
```

On Windows:
```
.\ouimap.exe
```

## Usage

When running OUImap, you can provide your search parameters interactively or by command line parameters.

1. Open a terminal and execute OUImap passing the MAC address and/or vendor name list.
2. Check the results to see which vendors and OUI ranges matched your query.

## Execution Example

```
$ ./ouimap
```
```
OUImap v1.2.0+build.g6ab869e (windows/amd64) 
Copyright © 2025 Adhoniran Gomes
This program comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it under certain conditions;
Visit https://www.gnu.org/licenses/gpl-3.0.html for details.


Enter a multi-line list of OUIs, MAC addresses and/or descriptions. Separate OUI/MAC address parts with colons, hyphens or periods.
Press ENTER on a blank line to start the search, or CTRL+C to exit.


suse linux
0050.56
00-19-1D-0F-DA-08

Searching...

0C:FD:37:00:00:00/24      SUSE Linux GmbH                                                                                                                                                                                           
00:50:56:00:00:00/24      VMware, Inc.                                                                                                                                                                                              
00:19:1D:00:00:00/24      Nintendo Co., Ltd.                                                                                                                                                                                        
                                                                                                                                                                                                                                    
>> 3 record(s) found from your search parameters...
>> Search completed in 2.1153ms.

```
or

```
$ ./ouimap 00-50-56 'atari, inc' 00:90:75:00:1A:BC
```
```
OUImap v1.2.0+build.g6ab869e (windows/amd64) 
Copyright © 2025 Adhoniran Gomes
This program comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it under certain conditions;
Visit https://www.gnu.org/licenses/gpl-3.0.html for details.

Searching...

00:50:56:00:00:00/24      VMware, Inc.
44:03:77:10:00:00/28      Atari, Inc.
00:90:75:00:00:00/24      Nec Do Brasil S.A.

>> 3 record(s) found from your search parameters...
>> Search completed in 2.1345ms.

```

## Automatic Updates

- OUImap periodically checks whether your local database is older than seven days. If so, it downloads the latest manufacturer list from the Wireshark server and automatically replaces your local copy.
- Additionally, whenever OUImap starts, it notifies if a new version is available.

## Contributing

Contributions are welcome! To contribute improvements, fixes, or new features:

1. Fork this repository.
2. Create a branch for your contribution:  
   git switch -b my-feature
3. Commit your changes:  
   git commit -m "Implemented a new feature"
4. Push to your repository:  
   git push origin my-feature
5. Open a Pull Request describing your changes so they can be reviewed.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE.md](LICENSE.md) file for details.

Copyright (C) 2025 Adhonian Gomes

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License only.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.


## Donations

Donations are greatly appreciated! You can make your contribution via the PayPal button below.

[![Donate with PayPal](https://img.shields.io/badge/Donate%20with-PayPal-blue?logo=paypal&logoColor=white)](https://www.paypal.com/donate/?business=N6AH25Q2D4BL8&no_recurring=0&item_name=Contribute+to+the+future+of+our+projects.+Your+donation+via+PayPal+empowers+us+to+keep+creating+and+growing...+Thank+you%21&currency_code=USD)
