
# Radar CLI

A Go command-line tool that I use as my personal swiss army knife 

## Building

1. **Prerequisites:** Ensure you have Go installed. ([Download Go](https://golang.org/))

2. **Get the code:** 
   * If you're using a version control system (like Git):
     ```bash
     git clone https://github.com/cristipercu/radar
     cd radar 
     ```
   * Otherwise, download the source code and navigate to the project's main directory

3. **Build:**
   ```bash
   go build


## Usage

```bash
./radar <command> [subcommand] [flags]
```

* `sync`: 
    * `create-config [-dirname <dir>]`: Create a config file (default: `.radar`)
    * `push`: Push local changes

* `mm`: Start the WIP... tool that does nothing, but keeps your laptop awake

* `--help`: Show help



## Examples

* `radar sync create-config`
* `radar sync create-config -dirname my_config`
* `radar sync push`
* `radar mm`


## More info
The sync command uses the rsync tool with a ssh connection, and if your connection uses a private key, you need to create a config file in the .ssh dir.

```bash
cd ~/.ssh
vi config
```

You can set a default setting for your ssh connection like this:
```
Host *
    IdentityFile /path/to/key
```

## Contributing

Fork and submit a pull request.

## License

