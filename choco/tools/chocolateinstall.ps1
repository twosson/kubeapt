
$ErrorActionPreference = 'Stop';
$packageName = 'apt'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64 = 'https://github.com/twosson/kubeapt/releases/download/v0.1.1/apt_0.1.1_Windows-64bit.zip'
$checksum64 = '9150f16fe49834ebd3ad065ac5a77acd506cd6ec5232fbe3b2fa191588670a47'
$checksumType64= 'sha256'

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  url64bit      = $url64

  softwareName  = 'apt*'

  checksum64    = $checksum64
  checksumType64= 'sha256'
}

Install-ChocolateyZipPackage @packageArgs

