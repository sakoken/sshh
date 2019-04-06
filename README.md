# sshh
sshh is a management application for ssh's host

## Install
Download from here https://github.com/sakoken/sshh/releases

## Usage
#### `sshh`
- Show the list of host
- Search a host
- Can do ssh connection
- Add,Modify,Delete a host

#### `sshh [hostname]`
After register the host, do ssh connection
If the host already is registered, do only ssh connection


## After exec sshh 
#### `sshh>> #[positionNo]`
Do ssh connection
It you can select by Up,Down Arrow

#### `sshh>> add`
Add a new host interactive<br>
You can register
- Host name
- User name
- Port number
- Password(Encrypted with a pass phrase)
- Explanation

#### `sshh>> mod [positionNo]`
Modify a host interactive

#### `sshh>> del [positionNo]`
Delete a host





