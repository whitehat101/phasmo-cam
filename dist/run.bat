@ECHO OFF
SETLOCAL
SET CHANNELID=<Number from "Copy ID" of target channel>
SET BOTTOKEN=<token from https://discord.com/developers/applications/xxx/bot>
SET SAVEDIR=C:\Users\%USERNAME%\Pictures\Phasmophobia
IF NOT EXIST %SAVEDIR% MKDIR %SAVEDIR%
.\phasmo-cam.exe --discord-token=%BOTTOKEN% --discord-channelID=%CHANNELID% --save-dir=%SAVEDIR%
PAUSE