# hotslogs
Command line replay uploader for HOTS Logs and HotS API.

## Background
After finding uploading HOTS replays to HOTS Logs through the browser to be tedious, yet not wanting to use a graphical uploader, I decided to write a command line tool for doing so. Upload functionality based on [eivindveg/HotSUploader](https://github.com/eivindveg/HotSUploader).

## Installation
If you have a Go environment set up: run `go get -u -v github.com/yi-jiayu/hotslogs` to add the `hotslogs` binary to `$GOPATH/bin`.

Otherwise, check [Releases](https://github.com/yi-jiayu/hotslogs/releases) for binaries.

## Usage
Run `hotslogs upload` to upload all new replays since the last time you ran the command.

```
PS C:\Users\jiayu> hotslogs up
Using config file: C:\Users\jiayu\.hotslogs\settings.yaml
[17:34] Starting upload to HOTS Logs
Found 2 new replays.
Uploaded Sky Temple.StormReplay (Success)...
Uploaded Tomb of the Spider Queen.StormReplay (Success)...
[17:34] Finished upload to HOTS Logs
[17:34] Starting upload to HotS API
Found 2 new replays.
Uploaded Sky Temple.StormReplay (Success)...
Uploaded Tomb of the Spider Queen.StormReplay (Success)...
[17:35] Finished upload to HotS API
```

## Hots Api support
`hotslogs` reads a configuration file at `$HOME/.hotslogs/settings.yaml`. Currently the only supported setting is where to upload your replays. A sample `settings.yaml` is:
```
destinations
  - hotslogs
  - hotsapi
```
Currently only hotslogs and hotsapi are supported as destinations.

## Roadmap
- [ ] Calculate replay ID and check if a replay has already been uploaded before uploading it.
