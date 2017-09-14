# hotslogs
Command line replay uploader for HOTS Logs

## Background
After finding uploading HOTS replays to HOTS Logs through the browser to be tedious, yet not wanting to use a graphical uploader, I decided to write a command line tool for doing so. Upload functionality based on [eivindveg/HotSUploader](https://github.com/eivindveg/HotSUploader).

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
Run `hotslogs upload` to upload all new replays since the last time you ran the command.

```
PS C:\Users\jiayu> hotslogs up
Using config file: C:\Users\jiayu\.hotslogs.yaml
Looking for new replays since: 2017-08-01 20:58:15 +0800 SGT
Found 4 new replay(s) since last upload.
Uploading new replays...
  Blackheart's Bay (16).StormReplay: DONE (Duplicate)
  Haunted Mines (20).StormReplay: DONE (Success)
  Infernal Shrines (26).StormReplay: DONE (Success)
  Towers of Doom (21).StormReplay: DONE (Success)
Updating config file... Done.
PS C:\Users\jiayu> 
```

## Hots Api support
`hotslogs upload` accepts the `--destinations` flag for choosing where to upload your replays to, which defaults to HOTS Logs. You can pass "hotsapi" to upload to Hots Api instead: `hotslogs up --destinations hotsapi`.

If you want to upload all your previous replays to Hots Api, you can use the `--config` flag with a dummy config file to ignore the last upload time: `hotslogs --config no-such-file up --destinations hotsapi`.

## Roadmap
- [ ] Calculate replay ID and check if a replay has already been uploaded before uploading it.
