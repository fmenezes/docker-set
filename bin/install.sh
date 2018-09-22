file="$GOPATH/src/github.com/fmenezes/docker-set/bin/docker-set.sh"
[[ -s ~/.bash_profile ]] && echo "[[ -s \"$file\" ]] && source \"$file\" # Load docker-set" >> ~/.bash_profile
[[ -s ~/.bashrc ]] && echo "[[ -s \"$file\" ]] && source \"$file\" # Load docker-set" >> ~/.bashrc
[[ -s ~/.zshrc ]] && echo "[[ -s \"$file\" ]] && source \"$file\" # Load docker-set" >> ~/.zshrc

[[ -s "$file" ]] && source "$file" # Load docker-set