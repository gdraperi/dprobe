for selfsigned in delgkey1 delgkey2 delgkey3 delgkey4; do
        subj='/C=US/ST=CA/L=SanFrancisco/O=Docker/CN=delegation'

        openssl genrsa -out "$***REMOVED***selfsigned***REMOVED***.key" 2048
        openssl req -new -key "$***REMOVED***selfsigned***REMOVED***.key" -out "$***REMOVED***selfsigned***REMOVED***.csr" -sha256 -subj "$***REMOVED***subj***REMOVED***"
        cat > "$***REMOVED***selfsigned***REMOVED***.cnf" <<EOL
[selfsigned]
basicConstraints = critical,CA:FALSE
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage=codeSigning
subjectKeyIdentifier=hash
EOL

        openssl x509 -req -days 3560 -in "$***REMOVED***selfsigned***REMOVED***.csr" -signkey "$***REMOVED***selfsigned***REMOVED***.key" -sha256 \
                -out "$***REMOVED***selfsigned***REMOVED***.crt" -extfile "$***REMOVED***selfsigned***REMOVED***.cnf" -extensions selfsigned

        rm "$***REMOVED***selfsigned***REMOVED***.cnf" "$***REMOVED***selfsigned***REMOVED***.csr"
done
