# Terminal reminder app

## Description
Hello there! =)

This is a simple terminal reminder app. It allows you to add, delete, and view reminders. Reminders are stored in-memory and are not persistent.

The idea behind this app is to use it for a short-term reminders, e.g. to remind you to do something in 2 hours.

Reminders are cross-platform compatible and are supposed to work on Windows, Linux, and MacOS (however, as of now it has only been tested on MacOS - I'll update the docs once more OS platforms have been tested).

Have fun!

## Installation
### MacOS
- via Homebrew:
```shell
brew tap n0rdy/n0rdy
brew install remindme
```

## Usage
### Starting app
To start the app, run the following command in the terminal:
```shell
remindme start
```

### Adding a reminder
There are several ways to add a reminder:

- to be reminded in a certain amount of time, e.g. in 2 hours 30 minutes 10 seconds:
```shell
remindme in --hr 2 --min 30 --sec 10 --about "Do something cool"
```

- to be reminded at a certain time, e.g. at 22:30:
```shell
remindme at --time 22:30 --about "Do something cool"
```
It is also to specify the time in 12-hour A.M./P.M format:
```shell
remindme at --pm 10:30 --about "Do something cool"
``` 
or for 10:30 A.M.:
```shell
remindme at --am 10:30 --about "Do something cool"
```

### List the existing reminders
- to see the list of all reminders, run:
```shell
remindme list
```

`list` command also supports sorting the list of reminders by ID, message or time in ascending or descending order. By default the list is provided in random order. If the sorting is requested, by the order of sorting is not specified, the ascending order is used.
- to sort by ID, run:
```shell
remindme list --sort --id
```
- to sort by message, run:
```shell
remindme list --sort --message
```
- to sort by time in descending order, run:
```shell
remindme list --sort --time --desc
```

### Canceling a reminder
- to cancel a reminder, run:
```shell
remindme cancel --id 1
```
where `1` is the ID of the reminder to be canceled. The ID can be obtained by running `remindme list` command.

- to cancel all reminders, run:
```shell
remindme cancel --all
```

### Changing the reminder message or/and time
- to change the message of a reminder, run:
```shell
remindme change --id 1 --message "Do something else, but also cool"
```
where `1` is the ID of the reminder to be changed. The ID can be obtained by running `remindme list` command.

- to change the time of a reminder, run:
```shell
remindme change --id 1 --time 22:30
```
where `1` is the ID of the reminder to be changed. The ID can be obtained by running `remindme list` command.

- it is possible to change the time by postponing it by a certain amount of time, e.g. by 2 hours 30 minutes 10 seconds:
```shell
remindme change --id 1 --postpone --hr 2 --min 30 --sec 10
```
where `1` is the ID of the reminder to be changed. The ID can be obtained by running `remindme list` command.

### Stopping the app
- to stop the app, run the following command in the terminal:
```shell
remindme stop
```

### Help
- to see the list of all available commands, run:
```shell
remindme help
```
Please, note that some flags have shortcuts, e.g. `--about` can be replaced with `-a`, `--time` with `-t`, etc. To see the full list of shortcuts, run the `help` command.
