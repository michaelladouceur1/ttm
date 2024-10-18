; Script generated by the Inno Setup Script Wizard.
; SEE THE DOCUMENTATION FOR DETAILS ON CREATING INNO SETUP SCRIPT FILES!

[Setup]
; NOTE: The value of AppId uniquely identifies this application.
; Do not use the same AppId value in installers for other applications.
; (To generate a new GUID, click Tools | Generate GUID inside the IDE.)
AppId={{c38c515e-c49a-45b2-a1c6-ebb8a36350be}}
AppName=TTM
AppVersion=0.1.0
AppPublisher=Michael LaDouceur  
AppPublisherURL=https://yourwebsite.com
AppSupportURL=https://yourwebsite.com/support
AppUpdatesURL=https://yourwebsite.com/updates
DefaultDirName={pf}\TTM
DefaultGroupName=TTM
OutputBaseFilename=TTMInstaller
Compression=lzma
SolidCompression=yes

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
; Specify the files to be included in the installer
Source: "path\to\your\binary\ttm.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "path\to\your\config\config.json"; DestDir: "{app}"; Flags: ignoreversion
Source: "path\to\your\database\ttm.db"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\TTM"; Filename: "{app}\ttm.exe"
Name: "{group}\{cm:UninstallProgram,TTM}"; Filename: "{uninstallexe}"

[Run]
; Add any additional actions to be performed after installation
Filename: "{app}\ttm.exe"; Description: "Launch TTM"; Flags: nowait postinstall skipifsilent