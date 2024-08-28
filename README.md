
# Radar CLI

A Go command-line tool that I use it as my personal swiss army knife 

## Building

1. **Prerequisites:** Ensure you have Go installed. ([Download Go](https://golang.org/))

2. **Get the code:** 
   * If you're using a version control system (like Git):
     ```bash
     git clone <repository-url>
     cd <project-directory> 
     ```
   * Otherwise, download the source code and navigate to the project's main directory

3. **Build:**
   ```bash
   go build


## Usage

```bash
./radar <command> [subcommand] [flags]

* `sync`: 
    * `create-config [-dirname <dir>]`: Create a config file (default: `.radar`)
    * `push`: Push local changes
    * `mm`: Keep laptop awake (WIP)

* `--help`: Show help

## Examples

* `radar sync create-config`
* `radar sync create-config -dirname my_config`
* `radar sync push`
* `radar sync mm`

## Contributing

Fork and submit a pull request.

## License

MIT License
