# sshh : ssh hosts manager

[![GitHub release](https://img.shields.io/github/release/sakoken/sshh.svg)](https://github.com/sakoken/sshh/releases/latest)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://github.com/sakoken/sshh/blob/master/LICENSE)
<br>
sshh is a management application of ssh's host<br>
You can ssh with your own secret phrase (server passwordless)<br>
![sshh](https://res.cloudinary.com/dwarv2f81/image/upload/v1554915394/sshh/sshh_movie.gif)

## Install
Download from here https://github.com/sakoken/sshh/releases

## Usage
#### `sshh`
- Show the list of host
- Search a host
- Can do ssh connection
- Add,Modify,Delete a host

#### `sshh [user]@[hostname] -p [port]`
Register and make SSH connection.<br>
If the host is already registered, only SSH connection will be made.<br>
Default port is 22

## After exec sshh 
#### `sshh>> #[positionNo]`
Do ssh connection.<br>
You can select with the up and down arrows.

#### `sshh>> add`
Add a new host<br>
You can register
- Host name
- User name
- Port number(Default 22)
- Password(Encrypted with a pass phrase)
- Explanation

#### `sshh>> mod [positionNo]`
Modify a host

#### `sshh>> del [positionNo]`
Delete a host

----

# License
MIT

# Author
Kenji Sakoda


