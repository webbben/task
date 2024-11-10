# task manager cli

This is a personal project for a task manager CLI tool. It's mainly designed to help me be more productive at work.

## Generating completions

If you need to generate new completions, do the following:

```bash
./task completion zsh > zsh_completion.sh
```

The above command generates the bash script for completions for all commands in this CLI tool, specifically for zsh. It saves it to `zsh_completion.sh`

you can then source it so that the completions work in the terminal:

```bash
source zsh_completion.sh
```
