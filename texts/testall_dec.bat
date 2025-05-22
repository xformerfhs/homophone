SetLocal EnableDelayedExpansion
for %%f IN (*_homophone.txt) DO call :processFile %%f
goto :eof

:processFile
..\homophone.exe d -in %1
exit /b
