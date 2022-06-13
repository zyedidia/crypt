# Crypt: create password-protected files

Crypt is a simple tool to encrypt/decrypt files and directories and lock them
behind a password. There are three commands:

## Lock

Creates a secure `.crypt` archive from a set of files.

```
$ crypt lock file.txt
```

## Unlock

Extracts a `.crypt` archive.

```
$ crypt unlock file.txt.crypt
```

## Op

Operates in-place on a `.crypt` archive. The `op` command lets you run an
arbitrary command on the decrypted archive, available in the temporary file
system. Once complete, the files will be re-locked, including any modifications
made, and the temporarily decrypted files will be removed. For example, you can
edit an encrypted file:

```
$ crypt op file.txt.crypt micro file.txt
```

Or simple open an interactive shell:

```
$ crypt op file.txt.crypt bash
```

You can then view/edit the files from in the shell.

# Install

```
go install github.com/zyedidia/crypt@latest
```
