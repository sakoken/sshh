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
Make SSH connection after host registered.<br>
If the host is already registered, only SSH connection will be made.

## After exec sshh 
#### `sshh>> #[positionNo]`
Do ssh connection.<br>
You can select with the up and down arrows.

#### `sshh>> add`
Add a new host<br>
You can register
- Host name
- User name
- Port number
- Password(Encrypted with a pass phrase)
- Explanation

#### `sshh>> mod [positionNo]`
Modify a host

#### `sshh>> del [positionNo]`
Delete a host





