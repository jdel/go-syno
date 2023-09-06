#!/bin/bash -e
package=$(go list)
domain=${package%/*}
prefix=${package#*/}
packagesJson=$(jq -nrR '[inputs]' <<< "$(go list -f '{{ .ImportPath }}' ./... | tail -n+2 | sed -e "s/${domain}\/${prefix}\///g")")

{
jq -rn \
   --arg domain "$domain" \
   --arg prefix "$prefix" \
   --argjson subs "$packagesJson" \
   '{
        "domain": $domain,
        "repositories": [
            {
                "prefix": $prefix,
                "subs": $subs,
                "url": [ "https://github.com/jdel/", $prefix ] | join("")
            }
        ]
    }'
} > vangen.json

git clone "https://${VANGEN_TOKEN}@github.com/jdel/jdel.github.io.git" html
vangen -config vangen.json -out html
cd html
git config user.name "github-actions"
git config user.email "github-actions@users.noreply.github.com"
git add "$prefix"
git commit -m "Vangen update for $domain/$prefix" || true
git push
