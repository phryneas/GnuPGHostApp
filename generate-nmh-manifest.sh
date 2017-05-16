[ -z "$1" ] && { echo please pass your extension id as first parameter; exit 1;  }

cat <<EOF > nmh-manifest.json
{
  "name": "de.phryneas.gpg.hostapp",
  "description": "GPG Host App",
  "path": "${PWD}/HostApp",
  "type": "stdio",
  "allowed_origins": [
    "chrome-extension://$1/"
  ]
}
EOF

cat <<EOF
nmh-manifest.json has been written

Windows/Chrome:
    REG ADD "HKCU\Software\Google\Chrome\NativeMessagingHosts\de.phryneas.gpg.hostapp" /ve /t REG_SZ /d "${PWD}/nmh-manifest.json" /f
MacOS/Chrome:
    mkdir -p "~/Library/Application Support/Google/Chrome/NativeMessagingHosts"
    cp nmh-manifest.json "~/Library/Application Support/Google/Chrome/NativeMessagingHosts/de.phryneas.gpg.hostapp.json"
MacOS/Chromium:
    mkdir -p "~/Library/Application Support/Chromium/NativeMessagingHosts"
    cp nmh-manifest.json "~/Library/Application Support/Chromium/NativeMessagingHosts/de.phryneas.gpg.hostapp.json"
Linux/Chrome:
    mkdir -p ~/.config/google-chrome/NativeMessagingHosts
    cp nmh-manifest.json ~/.config/google-chrome/NativeMessagingHosts/de.phryneas.gpg.hostapp.json
Linux/Chromium:
    mkdir -p ~/.config/chromium/NativeMessagingHosts
    cp nmh-manifest.json ~/.config/chromium/NativeMessagingHosts/de.phryneas.gpg.hostapp.json
EOF
