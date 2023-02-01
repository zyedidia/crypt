# Crypt: create password-protected files

Crypt is a simple tool to encrypt/decrypt files and directories and lock them
behind a password. There are three commands:

## Lock

Creates a secure `.crypt` archive from a set of files.

```
$ crypt lock file.txt
```

```
lock [OPTS] FILES:
  Archive and encrypt files or directories.
  -o string
    	output crypt file name
```

## Unlock

Extracts a `.crypt` archive.

```
$ crypt unlock file.txt.crypt
```

```
unlock [OPTS] FILE:
  Decrypt and extract crypt archives.
  -o string
    	directory to place extracted files
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

Or open an interactive shell and view/edit the files from there:

```
$ crypt op file.txt.crypt bash
```

The `-pipe` option will tell crypt to send the decrypted file contents to the
op command's stdin via a pipe rather than via the temporary file system. The
command's stdout will be written back to the file. This only works for
single-file crypt archives. The benefit of this approach is that the
decrypted file is only ever in process memory, rather than available in the
file system, which makes it more secure. For example to edit a file with micro and `-pipe`:

```
$ crypt op -pipe file.txt.crypt micro
```

```
op FILE COMMAND:
  Operate on a crypt file.
  -pipe
        pipe encrypted file instead of using temp files
```

# Install

Prebuilt binary (see releases):

```
eget zyedidia/crypt
```

From source:

```
go install github.com/zyedidia/crypt@latest
```
