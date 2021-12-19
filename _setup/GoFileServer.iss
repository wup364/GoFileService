[Setup]
AppName=File Service
AppVerName=GoFileService
VersionInfoVersion=1.0.0  
VersionInfoTextVersion=2020-08-16
LicenseFile=package-info\license.txt
DefaultDirName={code:getInstallDir}
DefaultGroupName=GoFileService
DisableProgramGroupPage=yes
;OutputDir=.\package-result
OutputDir=..\release\windows
OutputBaseFilename=GoFileService
Compression=lzma/max
SolidCompression=yes
PrivilegesRequired=admin
;SetupIconFile=.\package-info\SetupIcon.ico
ShowUndisplayableLanguages=yes

[Languages]
Name: "english"; MessagesFile: "package-language\English.isl";
Name: "chs"; MessagesFile: "package-language\Chinese.isl";

[Files]
Source: package-info\license.txt; DestDir: {app};
;Source: "package-source\*"; Excludes:""; DestDir: "{app}"; Flags: ignoreversion recursesubdirs createallsubdirs;
Source: "package-source\bin\*"; Excludes:""; DestDir: "{app}\bin"; Flags: ignoreversion recursesubdirs createallsubdirs;
Source: "..\bin\*"; Excludes:""; DestDir: "{app}"; Flags: ignoreversion recursesubdirs createallsubdirs;

[Icons]
;��ʼ�˵�
Name: "{group}\{cm:UninstallProgram,File Service}"; Filename: "{uninstallexe}"

[Registry]


[Run]
;����ǽ
Filename: "{sys}\netsh.exe"; Parameters: "firewall delete allowedprogram ""{app}\gofs.exe"" ";
Filename: "{sys}\netsh.exe"; Parameters: "firewall add allowedprogram ""{app}\gofs.exe"" ""Gofs Service"" ENABLE ALL";


[UninstallRun]

[UninstallDelete]
;�����װĿ¼���ļ�
Name: {app}; Type: filesandordirs


[Code]
//�ж��ַ����Ƿ��� ���֡�Ӣ����ĸ��ð�š���б�����
function IsChar(Str: string): Boolean;
var
  i: Integer;
  flag: Boolean;
begin
  Result   :=   True;
  for   i  :=   1   to   Length(Str)   do begin
    flag   :=   False;
    if ((Ord(Str[i])>=48) and (Ord(Str[i])<=57)) then begin
      flag := True;
    end;
    if ((Ord(Str[i])>=65) and (Ord(Str[i])<=90)) then begin
      flag := True;
    end;
    if ((Ord(Str[i])>=97) and (Ord(Str[i])<=122)) then begin
      flag := True;
    end;
    if ((Ord(Str[i])=58) or (Ord(Str[i])=92)) then begin
      flag := True;
    end;

    if not flag then begin
      Result   :=   False;
      Break;
    end;
  end;
end;

// ��ȡĬ�ϰ�װλ��
function getInstallDir(Param: String): String;
var 
    Tx_Disk: String;
begin
    Tx_Disk := ExpandConstant('{sd}')+'\';
    if DirExists('D:\') then
    begin
       Tx_Disk := 'D:\';
    end else begin 
        if DirExists('E:\') then
        begin
           Tx_Disk := 'E:\';
        end else begin
            if DirExists('F:\') then
            begin
               Tx_Disk := 'F:\';
            end else begin
                if DirExists('G:\') then
                begin
                   Tx_Disk := 'G:\';
                end;
            end;
        end;
    end;

    result := Tx_Disk+'GoFileServer';
end;

// �رս���
function closeRuningProgram():boolean;
var errorCode:Integer;
begin
  Exec(ExpandConstant('{cmd}'),'/C taskkill /f /im gofs.exe', '', SW_HIDE,ewWaitUntilTerminated, errorCode);
  Result:=True;
end;

// ��ʼ������
function installService():boolean;
var errorCode:Integer;
begin
  Exec(ExpandConstant('{app}')+'\bin\nssm.exe','install GoFileService '+ExpandConstant('{app}')+'\gofs.exe', '', SW_SHOWNORMAL,ewWaitUntilTerminated, errorCode);
  Exec(ExpandConstant('{app}')+'\bin\nssm.exe','set GoFileService AppDirectory '+ExpandConstant('{app}'), '', SW_SHOWNORMAL,ewWaitUntilTerminated, errorCode);
  Exec(ExpandConstant('{app}')+'\bin\nssm.exe','set GoFileService ObjectName \"Local Service\"', '', SW_SHOWNORMAL,ewWaitUntilTerminated, errorCode);
 
  Exec(ExpandConstant('{app}')+'\bin\nssm.exe','start GoFileService', '', SW_SHOWNORMAL,ewWaitUntilTerminated, errorCode);
  // Exec(ExpandConstant('{cmd}'), '/C '+ExpandConstant('{app}')+'\bin\openPage.url', '', SW_HIDE,ewWaitUntilTerminated, errorCode);
  Result:=True;
end;

// ж��ע��ķ���
function uninstallService():boolean;
var errorCode:Integer;
begin
  Exec(ExpandConstant('{app}')+'\bin\nssm.exe','stop GoFileService', '', SW_SHOWNORMAL,ewWaitUntilTerminated, errorCode);
  Exec(ExpandConstant('{app}')+'\bin\nssm.exe','remove GoFileService confirm', '', SW_SHOWNORMAL,ewWaitUntilTerminated, errorCode);
  Result:=True;
end;


function NextButtonClick(CurPageID: Integer): Boolean;
begin
  Result := True;
  if CurPageID = wpSelectDir then
  begin
    Result := IsChar(WizardDirValue);
    if not Result then MsgBox('Please choose other folder', mbError, MB_OK);
  end;
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep=usAppMutexCheck then
  begin
     uninstallService();
  end;
end;

procedure CurPageChanged(CurPageID: Integer);
begin

  if CurPageID = wpFinished then 
  begin
    installService();

  end else if CurPageID = wpReady then  
  begin
    

  end else if CurPageID = wpInstalling  then
  begin

  end;
end;

