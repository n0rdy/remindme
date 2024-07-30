# Terminal reminder app

## Description
Hello there!

The idea behind this app is to use it for short-term reminders, e.g. to remind you to do something in 2 hours or in 10 minutes.
The reminders are set in terminal, and the app will notify you about the reminder via OS-specific reminders.

Reminders are stored in SQLite database, so they will persist even if the app is stopped and started again.
However, if there is an issue with the SQLite setup, the app will still work, but the reminders will be stored in memory and will be lost after the app is stopped.
SQLite is bundled with the app, so there is no need to install it separately.

The flow of the app is as follows:
![](https://github.com/n0rdy/remindme/blob/master/docs/flow.gif)

Reminders are cross-platform compatible and are supposed to work on Windows, Linux, and MacOS.
So far, the app has been tested on:
- MacOS Ventura 13.5 with Intel Core i9 x64 CPU and zsh 5.9 (x86_64-apple-darwin22.1.0)
- MacOS Ventura 13.5 with Apple M2 Max x64 CPU and zsh 5.9 (x86_64-apple-darwin22.0)
- Ubuntu 22.04.2 LTS (via VirtualBox) with Intel Core i9 x64 CPU and bash 5.1.16(1)-release (x86_64-pc-linux-gnu)
- Fedora Workstation 38 (via VirtualBox) with Intel Core i9 x64 CPU and bash 5.2.15(1)-release (x86_64-redhat-linux-gnu)
- Windows 10 Home (version 21H2) with Intel Core i5 x64 CPU and Windows PowerShell 5.1.19041.1682

Have fun! =)

## Installation
### Prerequisites
- [Go](https://golang.org/doc/install) (version 1.22 or higher) if you want to build the app from the source code.

### Manual
Download the latest release for your OS from [GitHub](https://github.com/n0rdy/remindme/releases).

### MacOS
- via Homebrew:
```shell
brew tap n0rdy/n0rdy
brew install remindme
```

### Linux
#### via APT
##### Prerequisites
To enable, add the following file /etc/apt/sources.list.d/fury.list:
```text
deb [trusted=yes] https://apt.fury.io/n0rdy/ /
```
You can do this either manually or by running the following command:
```shell
echo "deb [trusted=yes] https://apt.fury.io/n0rdy/ /" > /etc/apt/sources.list.d/fury.list
```
If you experienced this error:
```text
bash: /etc/apt/sources.list.d/fury.list: Permission denied
```
try to do the following:
```shell
sudo -i
# enter your password
echo "deb [trusted=yes] https://apt.fury.io/n0rdy/ /" > /etc/apt/sources.list.d/fury.list
# click Ctrl+D to exit
```

##### Installation
```shell
sudo apt update && sudo apt install remindme
```

#### via YUM
##### Prerequisites
To enable, add the following file /etc/yum.repos.d/fury.repo:
```text
[fury]
name=Gemfury n0rdy Private Repo
baseurl=https://yum.fury.io/n0rdy/
enabled=1
gpgcheck=0
```

##### Installation
```shell
sudo yum install remindme
```

## Usage
### Starting app
To start the app, run the following command in the terminal:
```shell
remindme start
```
Under the hood, the HTTP server is started on port 15555.

It is possible to specify the port to start the server on by either:
- using the `--port` flag, e.g.:
```shell
remindme start --port 15556
```
- setting the `REMINDME_PORT` environment variable to the desired port number.

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
It is also possible to specify the time in 12-hour A.M./P.M format:
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

The `list` command also supports sorting the list of reminders by ID, message or time in an ascending or descending order. 
By default, the list is provided in a random order. If the sorting is requested, but the sorting order is not specified, the ascending order is used.
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
where `1` is the ID of the reminder to be cancelled. The ID can be obtained by running `remindme list` command.

- to cancel all reminders, run:
```shell
remindme cancel --all
```

### Changing the reminder message or/and time
- to change the message of a reminder, run:
```shell
remindme change --id 1 --about "Do something else, but also cool"
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

### Logs
#### Printing logs
- to print the logs, run the following command in the terminal:
```shell
remindme admin logs print
```
It is possible to specify which logs to print: client or server logs by using the `--client` or `--server` flags respectively.
By default, client logs are printed.

#### Deleting log files
- to delete the log files, run the following command in the terminal:
```shell
remindme admin logs delete
```
It is possible to specify which log file to delete: client or server logs by using the `--client` or `--server` flags respectively.
By default, both client and server log files are deleted.

### Help
- to see the list of all available commands, run:
```shell
remindme help
```

### Completion
remindme app provides a possibility to use completion for the commands and flags. To generate the completion, run the following command in the terminal:
```shell
remindme completion
```
This will generate the completion for the current shell. 
To enable the completion for the current shell run, use this command:
```shell
source <(remindme completion)
```
To enable the completion for the current shell permanently, add this line to your `.bashrc`, `.zshrc` or `config.fish` file:
```shell
source <(remindme completion)
```

If you are on Windows with PowerShell, it is possible to generate the completion by running the following command:
```shell
remindme completion
```
Please, check the [PowerShell documentation](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/register-argumentcompleter?view=powershell-7.3) for more information about loading completions for this shell type.

### Documentation generation
To generate the documentation for the app in the MD-format, run the following command in the terminal:
```shell
remindme docs --dir "path_to_dir"
```
If the `--dir` flag is not specified, the documentation will be generated in the temp directory on your machine - the path to the directory will be printed to the terminal as a result of `docs` command execution.
You can find the same docs (as generated) within the [docs](https://github.com/n0rdy/remindme/blob/master/docs/remindme.md) folder on GitHub.

### Additional information
Please, note that some flags have shortcuts, e.g. `--about` can be replaced with `-a`, `--time` with `-t`, etc. To see the full list of shortcuts, run the `remindme help`.
