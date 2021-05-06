# BlackOut
BlackOut is a cookie extractor for Chrome on macOS which uses the Keychain to decrypt cookie values. You can export to CSV or JSON, or print to STDOUT.

## Installing
Either download the correct binary for your platform, or build locally using `make build`.

## Usage
- `-o` will take the output filename. Passing `.json` as a file extension outputs cookies in JSON format. By using `.csv` or omitting a file extension, CSV will be used. Defaults to `stdout`.

- `-p` defines the Chrome Profile to use. Defaults to `Default`.

- `-d` performs domain filtering on the cookie hostname. Defaults to all.
