## remindme change

Change reminder message and/or notification time

### Synopsis

Change reminder message and/or notification time.

The command expects a reminder ID to be provided via the "--id" flag - otherwise, the error will be produced.
Additionally, "--about", "--time" and/or "--postpone" flags should be provided - otherwise, the error will be produced.
Either "--time" or "--postpone" should be provided, not both - otherwise, the error will be produced.
It is valid to provide both "--about" and "--time" OR "--about" and "--postpone" flags together.

If the "--postpone" flag is provided, "--sec", "--min" and/or "--hr" flags should be provided alongside - otherwise, the error will be produced.
Negative integer values are not accepted - the error will be produced in such case.

List the upcoming reminders with the "list" command.

```
remindme change [flags]
```

### Options

```
  -a, --about string      Reminder message to change to
  -h, --help              help for change
      --hr --postpone     Hours to shift the existing notification time with - should be passed alongside the --postpone flag
      --id int            Reminder ID to change
      --min --postpone    Minutes to shift the existing notification time with - should be passed alongside the --postpone flag
      --postpone --time   If provided, specifies that no new time will be provided by --time flag, but rather a desired shift in time using `--sec`, `--min` and/or `--hr` flags (e.g. if the notification time is 15:30 and `--postpone --min 20` is provided, the new time will be 15:50)
      --sec --postpone    Seconds to shift the existing notification time with - should be passed alongside the --postpone flag
  -t, --time string       Time to change the notification to in 24-hours HH:MM format: e.g. 16:30, 07:45, 00:00
```

### SEE ALSO

* [remindme](remindme.md)	 - A tool to set reminders from the terminal

###### Auto generated by spf13/cobra on 6-Aug-2023