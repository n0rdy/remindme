## remindme admin logs delete

Delete logs files

### Synopsis

Delete logs files.

Accepts either --server or --client flag to specify which logs to delete.
If no flag is provided, deletes both client and server logs by default.

If both flags are provided, the error message is printed.

If the remindme app didn't manage to find the logs file, nothing is deleted.

```
remindme admin logs delete [flags]
```

### Options

```
  -h, --help   help for delete
```

### Options inherited from parent commands

```
  -c, --client   Request client logs to be printed (default true)
  -s, --server   Request server logs to be printed
```

### SEE ALSO

* [remindme admin logs](remindme_admin_logs.md)	 - Admin logs commands: print and delete logs

###### Auto generated by spf13/cobra on 13-Aug-2023
