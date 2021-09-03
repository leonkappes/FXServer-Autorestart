# FXServer-Restarter
This is a script, written in Go(-lang) for restarting an FXServer.exe process, which is the Server executable of the [FiveM](https://fivem.net) Platform.
Automaticlly restarts at the given hours in `restart.csv`, which are full hours seperated by commas. Running this scripts enables unattendes restarts.
The Server executable is killed by the default windows executable `taskkill.exe`, while it checks for a running server with `tasklist.exe` and starts the server in a newly
summoned `cmd.exe` shell-window.

Pressing `r` in the terminal window allows you to do a manual restart.
