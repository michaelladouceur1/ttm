; ttm_installer.nsi
Name "TTM"
OutFile "ttm_setup.exe"
InstallDir "$PROGRAMFILES\TTM"
Section
  SetOutPath $INSTDIR
  File "ttm.exe"
SectionEnd