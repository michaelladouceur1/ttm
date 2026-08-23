Unicode true
RequestExecutionLevel admin

!ifndef VERSION
  !define VERSION "dev"
!endif

Name "TTM"
OutFile "ttm_setup.exe"
InstallDir "$PROGRAMFILES\TTM"
ShowInstDetails show
ShowUnInstDetails show

Page directory
Page instfiles
UninstPage uninstConfirm
UninstPage instfiles

Section "TTM" SecTTM
  SetOutPath "$INSTDIR"
  File "ttm.exe"
  WriteUninstaller "$INSTDIR\Uninstall.exe"

  CreateDirectory "$SMPROGRAMS\TTM"
  CreateShortcut "$SMPROGRAMS\TTM\TTM.lnk" "$INSTDIR\ttm.exe"
  CreateShortcut "$SMPROGRAMS\TTM\Uninstall TTM.lnk" "$INSTDIR\Uninstall.exe"

  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TTM" "DisplayName" "TTM"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TTM" "DisplayVersion" "${VERSION}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TTM" "DisplayIcon" "$INSTDIR\ttm.exe"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TTM" "UninstallString" "$\"$INSTDIR\Uninstall.exe$\""
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TTM" "NoModify" 1
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TTM" "NoRepair" 1
SectionEnd

Section "Uninstall"
  Delete "$SMPROGRAMS\TTM\TTM.lnk"
  Delete "$SMPROGRAMS\TTM\Uninstall TTM.lnk"
  RMDir "$SMPROGRAMS\TTM"
  Delete "$INSTDIR\ttm.exe"
  Delete "$INSTDIR\Uninstall.exe"
  RMDir "$INSTDIR"
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TTM"
SectionEnd
