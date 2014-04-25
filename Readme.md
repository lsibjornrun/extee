Usage: `extee -e <regular expression> -x <command to execute> [-l <logfile from command>] [-q] [-v]`

Example:
`cat sample.txt | extee -e "echo (?P<def>[a-z]+)" -x "echo <def>" -l   log.txt`

Named variables in regular expression:
`?P<name>  ` 

This can be used in command as 
`<name>`

The log file will be created if it does not exists, and appended if it exists.

Use -q to quiet the output. If -l is not used, only output from command execution will be shown.

Use -v to increase verbosity and help debug regular expression.
