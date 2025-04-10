# c2k

A tool for converting a curl command to the corresponding Kotlin code with Ktor.

## Usage

Just place the `c2k` command before the `curl` command. For example:
```shell
c2k curl -v --data-raw 'request-body' https://example.com
```

As a result, the generated code will be printed to stdout.

## Install

1. Go to the [GitHub releases page](https://github.com/Stexxe/c2k/releases) and download the appropriate binary
2. Make it executable `chmod +x <binary-name>  # Replace with the correct filename`
3. Move it to a PATH directory: `sudo mv <binary-name> /usr/local/bin/c2k`
4. Verify it works: `c2k --version`