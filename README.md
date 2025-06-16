# pflags

## Build
  ```sh
  rm -f -r /tmp/pflags
  git clone https://github.com/chahal-p/pflags.git /tmp/pflags
  cd /tmp/pflags
  make build
  
  ```
## Install
  ### User installation
  ```sh
  make install INSTALLATION_PATH=$HOME/.local/bin
  ```
  Make sure `$HOME/.local/bin` is path of PATH env variable.
  If not export this.
  ```sh
  export PATH="$PATH:$HOME/.local/bin"
  ```
  For bash, simply below snippet can be added to `.bashrc`
  ```sh
  [[ ":$PATH:" == *":$HOME/.local/bin:"* ]] || export PATH="$PATH:$HOME/.local/bin"
  ```
  ### System installation
  ```sh
  sudo make install INSTALLATION_PATH=/usr/local/bin
  ```
## Uninstall
  From user:
  ```sh
  make install INSTALLATION_PATH=$HOME/.local/bin
  ```
  From system:
  ```sh
  sudo make install INSTALLATION_PATH=/usr/local/bin
  ```

## Usage
  ```sh
  pflags --help
  ```
