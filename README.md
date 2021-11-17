# Speeddial

Remember complex shell commands with ease, and never forget another one!


## Usage

```sh
# Search over saved commands and prefill the next prompt
$ spd

# Add the previous command to Speeddial
$ spd add

# Add the specified command to Speeddial
$ spd add <command string>

# Remove a command
$ spd rm
```

## How to Install

### Download/Build
#### From GitHub Releases
1. Download the relevant binary from the latest release.

2. Add the downloaded `speeddial` binary to your PATH (move it into a relevant directory etc.)

3. Follow the steps in the `Initialize` section.

#### From Source
1. Clone the repo.

    ```sh
    git clone github.com/rithvikp/speeddial
    ```

2. Build the package.

    ```bash
    make build
    ```

3. Add the created `speeddial` binary to your PATH (move it into a relevant directory etc.)

4. Follow the steps in the `Initialize` section.

### Initialize

2. Update the relevant shell configuration file to initialize Speeddial for every shell session.

    Zsh:
    ```sh
    eval "$(speeddial init zsh)"
    ```

    Fish:
    ```sh
    speeddial init fish | source
    ```

3. After sourcing your shell configuration (or creating a new terminal session), start using
   Speeddial through the `spd` command name!
