# hotslogs
Command line replay uploader for HOTS Logs

## Installation
If you have a Go environment set up: run `go get -u -v github.com/yi-jiayu/hotslogs` to add the `hotslogs` binary to `$GOPATH/bin`.

Otherwise, check [Releases](https://github.com/yi-jiayu/hotslogs/releases) for binaries.

## Setup
`hotslogs` will work out of the box if your replays are located in the default locations. 

According to the HOTS Logs [upload page](https://www.hotslogs.com/Account/Upload), the default folder for replays is:
- Windows: `$HOME\Documents\Heroes of the Storm\Accounts\########\#-Hero-#-######\Replays\`
- Mac: `~/Library/Application Support/Blizzard/Heroes of the Storm/Accounts/########/#-Hero-#-######/Replays/`

Otherwise, you can run `hotslogs config init` to set your replay directory manually.

## Usage
Run `hotslogs update` to upload all new replays since the last time you ran the command.

## Roadmap
- [ ] Calculate replay ID and check if a replay has already been uploaded before uploading it.
