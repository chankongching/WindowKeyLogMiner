Dim WinScriptHost
Set WinScriptHost = CreateObject("WScript.shell")
WinScriptHost.Run Chr(34) & "x:\WindowKeyLogMiner\WindowKeyLogMiner.exe" & Chr(34),0
Set WinScriptHost = Nothing