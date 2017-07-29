# `hotslogs`
Command line replay uploader for HOTS Logs

## Installation
If you have a Go environment set up: `go get -u -v github.com/yi-jiayu/hotslogs` will add the `hotslogs` binary to `$GOPATH/bin`.

Otherwise, check in the releases section for binary downloads and add to your `$PATH`.

## Setup
Run `hotslogs config init` and set your replay directory, which may or may not be automatically detected.

According to the HOTS Logs [upload page](https://www.hotslogs.com/Account/Upload), the default folder for replays is:
- Windows: `$HOME\Documents\Heroes of the Storm\Accounts\########\#-Hero-#-######\Replays\`
- Mac: `~/Library/Application Support/Blizzard/Heroes of the Storm/Accounts/########/#-Hero-#-######/Replays/`

## Usage
Run `hotslogs update` to upload all new replays since the last time you ran the command.

## Roadmap
- [ ] Calculate replay ID and check if a replay has already been uploaded before uploading it.
