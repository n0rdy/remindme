## remindme at

Create a reminder to be notified about it at some point in time

### Synopsis

Create a reminder to be notified about it at some point in time.

The command accepts the exact time the notification should be sent at in either 24 hours "hh:mm" format (e.g. 13:05 or 09:45) via --time flag, or 12 hours "hh:mm" A.M./P.M. format via --am/--pm flag.
The provided time should be in future - otherwise, the error will be produced.
The command expects a reminder message to be provided via the "--about" flag - otherwise, the error will be produced.

List the upcoming reminders with the "list" command.

```
remindme at [flags]
```

### Options

```
  -a, --about string   Reminder message
      --am at          A.M. time to remind at for at command in 12-hours HH:MM format: e.g. 07:45
  -h, --help           help for at
      --pm at          P.M. time to remind at for at command in 12-hours HH:MM format: e.g. 07:45
  -t, --time at        Time to remind at for at command in 24-hours HH:MM format: e.g. 16:30, 07:45, 00:00
```

### SEE ALSO

* [remindme](remindme.md)	 - A tool to set reminders from the terminal

###### Auto generated by spf13/cobra on 13-Aug-2023
