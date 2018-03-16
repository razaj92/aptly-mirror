#!/bin/sh

# VARIABLES
export GNUPGHOME=${GNUPGHOME:-"/var/lib/aptly/gpg"}

echo "-> Configuring aptly-mirror"

mkdir -p $GNUPGHOME && chmod 700 $GNUPGHOME

if [[ -f /gpg_pub.gpg ]] && [[ -f /gpg_key.gpg ]] ; then

  echo "-> Importing GPG Keys"

  gpg --batch --import /gpg_pub.gpg
  gpg --batch --allow-secret-key-import --import /gpg_key.gpg

elif [[ ! -f "${GNUPGHOME}/pubring.gpg" ]]; then

  echo "-> No Keys found.. creating generic keys"

  cat <<EOT >> /gpg.txt
Key-Type: 1
Key-Length: 2048
Subkey-Type: 1
Subkey-Length: 2048
Name-Real: foo bar
Name-Email: root@foo.bar
Expire-Date: 0
EOT

  gpg --batch --gen-key /gpg.txt

fi

aptly-mirror $@
