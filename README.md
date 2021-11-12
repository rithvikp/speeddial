# Speeddial

Remember complex shell commands with ease, and never forget another one!


## How to Install (To be improved)

### From Source
1. Clone the repo

    ```bash
    git clone github.com/rithvikp/speeddial
    ```

2. Build the package

    ```bash
    go build github.com/rithvikp/speeddial
    ```

3. Add the created `speeddial` binary to your path (move it into a relevant directory etc.)

4. Update the relevant shell configuration file to initialize Speeddial for every shell session.
   Currently only Zsh is supported.

    Zsh:
    ```bash
    eval "$(speeddial init zsh)"
    ```

5. After sourcing your shell configuration (or creating a new terminal session), start using
   Speeddial through the `spd` command name!
