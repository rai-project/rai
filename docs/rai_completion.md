---
date: 2017-03-12T13:24:31Z
title: "rai completion"
slug: rai_completion
url: /commands/rai_completion/
---

## rai completion

Generates bash completition files.

### Synopsis

Output shell completion code for the specified shell (bash or zsh).
The shell code must be evalutated to provide interactive
completion of kubectl commands.  This can be done by sourcing it from
the .bash_profile.
Note: this requires the bash-completion framework, which is not installed
by default on Mac.  This can be installed by using homebrew:

    $ brew install bash-completion

Once installed, bash_completion must be evaluated.  This can be done by adding the
following line to the .bash_profile

    $ source $(brew --prefix)/etc/bash_completion

Note for zsh users: [1] zsh completions are only supported in versions of zsh >= 5.2

#### Install bash completion on a Mac using homebrew

brew install bash-completion
printf "\\n# Bash completion support\\nsource $(brew --prefix)/etc/bash_completion\\n" >> $HOME/.bash_profile
source $HOME/.bash_profile

#### Load the rai completion code for bash into the current shell

  source &lt;(rai completion)

#### Write bash completion code to a file and source if from .bash_profile

  rai completion bash > ~/.rai/completion.bash.inc
  printf "\\n# Rai shell completion\\nsource '$HOME/.rai/completion.bash.inc'\\n" >> $HOME/.bash_profile
  source $HOME/.bash_profile

#### Load the rai completion code for zsh[1] into the current shell

    source &lt;(rai completion)

    rai completion

### Options inherited from parent commands

      -c, --color         Toggle color output.
      -d, --debug         Toggle debug mode.
      -v, --verbose       Toggle verbose mode.

### SEE ALSO

-   [rai](rai.md)	 - The client is used to submit jobs to the server.
