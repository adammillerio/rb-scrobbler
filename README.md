# Minimal Rockbox Scrobbler (rb-scrobbler)

_rb-scrobbler_ is a minimal CLI tool which allows a user to scrobble(submit) .scrobbler.log files generated by Rockbox CFW for DAPs to [last.fm](https://last.fm).
![badge-as](badge.gif)

![usage](in-action.png)

Compared to other tools rb-scrobbler has advantages such as:

        1. Login via session key - Your password is not saved anywhere on disk
        2. Does not rely on outdated dependencies
        3. Cross platform with any OS/Arch that can compile Rust
        4. Single binary - no extra dependencies required

## Usage

Pre compiled binaries for Windows and Linux x64 are available in releases.

The following instructions assume that the binary is within your PATH or your prompt is inside the binary's directory.

(**Windows**: Shift + Right Click -> Open Command Window Here
**Linux**: You already know :) )

On first start you'll need to authenticate yourself with last.fm

`rb-scrobbler -f=[FILE PATH TO .scrobbler.log] -o=[OFFSET FROM UTC (IN HOURS)] auth -u=[USERNAME] -p[PASSWORD]`

NOTE: Offset is in hours. If you say happen to live in a timezone -10:30 from UTC enter it as -10.5

Further uses of the program can be executed with just the file path and timezone offset. Further authentication uses a session key only.
Your password is never kept on disk and is only used for the initial authentication.

### Compilation

If a binary does not exist for your system compilation is trivial.
All that is required is an installation of the Rust toolchain and a last.fm [API key](https://www.last.fm/api/account/create).

        1. Navigate to `src/api_keys.rs` in your editor of choice and fill in the empty fields with your API Keys
        2. In the main project directory run `cargo build --release`
        3. Your compiled binary will be within 'target/release'

## .scrobbler.log

A .scrobbler.log file looks like this:

(placeholder)

It's a tab separated file of fields [ARTIST, ALBUM, TRACK, TRACK_POS, DURATION(sec), "RATING", TIMESTAMP(unix)].
rb-scrobbler iterates over each line, looks if the "RATING" field is marked "listened" (\tL\t),
copies **Artist**, **Album**, **Track** & **Timestamp** (and converts it back/forward to UTC if necessary) and then submits it to your
last.fm.

Official documentation on this format is available [here](https://web.archive.org/web/20170107015006/http://www.audioscrobbler.net/wiki/Portable_Player_Logging).
I haven't implemented the specification perfectly and have just implemented the parts implemented by Rockbox (_i.e_ RB doesn't keep timezone data
so rb-scrobbler doesn't parse the file looking for it)

As with everything built on the last.fm API; scrobbles older than two weeks will not propagate on their database despite what the response will tell you.

## TODO

        * Send batches of scrobbles instead of one by one to reduce request amounts
        * If you try to reauthenticate, the program panics when trying to create a new directory, will fix soon
        * Use the desktop authentication stream

## Built With

        * [rustfm-scrobble](https://github.com/bobbo/rustfm-scrobble) - This library is seriously simple and allowed me to bang this out in a week
        * [dirs-rs](https://github.com/soc/dirs-rs) - Intuitive, small, cross-platform library for getting config directories

## Author

        * Eddie Jeselnik 2020

## License

This project is licensed under the GNU GPL v3 - see [LICENSE](LICENSE) for details.
