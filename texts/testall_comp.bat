@echo off
SetLocal EnableDelayedExpansion
for %%f IN (*.txt) DO call :eventuallyProcessFile %%f
goto :eof

:eventuallyProcessFile
set x=%1
if /i "%x%" EQU "!x:_homophone=!" call :eventuallyProcessFile2 %1
exit /b

:eventuallyProcessFile2
set y=%1
if /i "%y%" EQU "!y:_decrypted=!" call :processFile %1
exit /b

:processFile
set n=%~n1
set x=%~x1
set decr=%n%_decrypted%x%
if not exist %decr% exit /b
echo on
fc /c %1 %decr%
@echo off
exit /b
